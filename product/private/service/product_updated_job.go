package service

import (
	"context"
	"fmt"
	"github.com/go-saas/kit/pkg/job"
	"github.com/go-saas/kit/pkg/query"
	stripe2 "github.com/go-saas/kit/pkg/stripe"
	"github.com/go-saas/kit/product/private/biz"
	"github.com/go-saas/saas"
	"github.com/go-saas/uow"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v76"
	stripeclient "github.com/stripe/stripe-go/v76/client"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"strings"
	"time"
)

const (
	JobTypeProductUpdated = "product" + ":" + "product_updated"
)

func NewProductUpdatedTask(prams *biz.ProductUpdatedJobParam) (*asynq.Task, error) {
	payload, err := protojson.Marshal(prams)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(JobTypeProductUpdated, payload, asynq.Queue(string(biz.ConnName)), asynq.Retention(time.Hour*24*30)), nil
}

func NewProductUpdatedTaskHandler(productRepo biz.ProductRepo, priceRepo biz.PriceRepo, uowMgr uow.Manager, client *stripeclient.API) *job.Handler {
	return job.NewHandlerFunc(JobTypeProductUpdated, func(ctx context.Context, t *asynq.Task) error {
		return uowMgr.WithNew(ctx, func(ctx context.Context) error {
			msg := &biz.ProductUpdatedJobParam{}
			if err := protojson.Unmarshal(t.Payload(), msg); err != nil {
				return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
			}
			//change to product tenant
			ctx = saas.NewCurrentTenant(ctx, msg.TenantId, "")
			var product *biz.Product
			var err error
			if !msg.IsDelete {
				product, err = productRepo.Get(ctx, msg.ProductId)
				if err != nil {
					return err
				}
				if product == nil {
					return fmt.Errorf("can not find %s: %w", msg.ProductId, asynq.SkipRetry)
				}
				if product.Version.String != msg.ProductVersion {
					return fmt.Errorf("product:%s version mismatch, should be:%s, got:%s: %w", msg.ProductId, msg.ProductVersion, product.Version.String, asynq.SkipRetry)
				}
			}

			group, ctx := errgroup.WithContext(ctx)
			//sync with stripe
			group.Go(func() error {
				return syncWithStripe(ctx, client, productRepo, priceRepo, product, msg)
			})
			return group.Wait()

		})
	})
}

func syncWithStripe(ctx context.Context, client *stripeclient.API, productRepo biz.ProductRepo, priceRepo biz.PriceRepo, product *biz.Product, jobParams *biz.ProductUpdatedJobParam) error {
	if client == nil {
		return nil
	}
	if jobParams.IsDelete {
		stripeInfo, ok := lo.Find(jobParams.SyncLinks, func(t *biz.ProductUpdatedJobParam_SyncLink) bool {
			return t.ProviderName == string(biz.ProductManageProviderStripe)
		})
		if ok {
			// stripe product which has prices can't be deleted
			//_, err := client.Products.Del(stripeInfo.ProviderId, &stripe.ProductParams{})
			//if err != nil {
			//	return err
			//}
			_, err := client.Products.Update(stripeInfo.ProviderId, &stripe.ProductParams{Active: stripe.Bool(false)})
			if err != nil {
				return err
			}
		}
		return nil
	}

	links, err := productRepo.GetSyncLinks(ctx, product)
	if err != nil {
		return err
	}
	stripeInfo, find := lo.Find(links, func(link biz.ProductSyncLink) bool {
		return string(biz.ProductManageProviderStripe) == link.ProviderName
	})
	var stripeProductId string
	if !find {
		// create stripe object
		params := mapBizProduct2Stripe(product)
		stripeProduct, err := client.Products.New(params)
		if err != nil {
			return err
		}
		stripeProductId = stripeProduct.ID
		t := time.Now()
		err = productRepo.UpdateSyncLink(ctx, product, &biz.ProductSyncLink{
			ProviderName: string(biz.ProductManageProviderStripe),
			ProviderId:   stripeProductId,
			LastSyncTime: &t,
		})
		if err != nil {
			return err
		}
	} else {
		stripeProductId = stripeInfo.ProviderId
		//update product if needed
		stripeProduct, err := client.Products.Get(stripeProductId, &stripe.ProductParams{})
		if err != nil {
			return err
		}
		params := mapBizProduct2Stripe(product)
		_, err = client.Products.Update(stripeProduct.ID, params)
		if err != nil {
			return err
		}
	}

	// update price
	allKeys := lo.Map(product.Prices, func(t biz.Price, _ int) string {
		return t.ID.String()
	})
	priceIter := client.Prices.List(&stripe.PriceListParams{Product: &stripeProductId})
	if priceIter.Err() != nil {
		return priceIter.Err()
	}
	var stripePrices []*stripe.Price
	for priceIter.Next() {
		if priceIter.Err() != nil {
			return priceIter.Err()
		}
		stripePrices = append(stripePrices, priceIter.Price())
	}
	//delete  prices
	toDeactivate := lo.Filter(stripePrices, func(price *stripe.Price, _ int) bool {
		return !lo.Contains(allKeys, price.LookupKey)
	})
	for _, price := range toDeactivate {
		active := false
		_, err = client.Prices.Update(price.ID, &stripe.PriceParams{Active: &active})
		if err != nil {
			return err
		}
	}
	for _, price := range product.Prices {
		stripePrice, ok := lo.Find(stripePrices, func(p *stripe.Price) bool {
			return p.LookupKey == price.ID.String()
		})

		if !ok {
			params := mapBizPrice2CreateStripe(stripeProductId, &price)
			//create
			stripePrice, err = client.Prices.New(params)
			if err != nil {
				return err
			}
		} else {
			params := mapBizPrice2UpdateStripe(stripeProductId, &price)
			_, err = client.Prices.Update(stripePrice.ID, params)
			if err != nil {
				return err
			}
		}
		if price.StripePriceId == nil {
			price.StripePriceId = &stripePrice.ID
			err = priceRepo.Update(ctx, price.ID.String(), &price, query.NewField(&fieldmaskpb.FieldMask{Paths: []string{"stripe_price_id"}}))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func mapBizProduct2Stripe(product *biz.Product) *stripe.ProductParams {
	return &stripe.ProductParams{
		Active:      stripe.Bool(product.Active),
		Name:        stripe2.String(product.Title),
		Description: stripe2.String(product.Desc),
		Shippable:   stripe.Bool(product.NeedShipping),
		Metadata:    map[string]string{"version": product.Version.String},
		//TODO type
	}
}

func mapBizPrice2CreateStripe(stripeProductId string, price *biz.Price) *stripe.PriceParams {
	r := &stripe.PriceParams{
		BillingScheme: stripe2.String(string(price.BillingScheme)),
		Currency:      stripe2.String(strings.ToLower(price.CurrencyCode)),
		LookupKey:     stripe2.String(price.ID.String()),
		Product:       stripe2.String(stripeProductId),
		TiersMode:     stripe2.String(string(price.TiersMode)),
	}
	if price.TransformQuantity.DivideBy >= 1 {
		r.TransformQuantity = &stripe.PriceTransformQuantityParams{
			DivideBy: stripe.Int64(price.TransformQuantity.DivideBy),
			Round:    stripe2.String(string(price.TransformQuantity.Round)),
		}
	}
	r.UnitAmount = stripe.Int64(price.GetNeedPayAmount())

	if len(price.CurrencyOptions) > 0 {
		r.CurrencyOptions = lo.SliceToMap(price.CurrencyOptions, func(t biz.PriceCurrencyOption) (string, *stripe.PriceCurrencyOptionsParams) {
			cop := &stripe.PriceCurrencyOptionsParams{
				UnitAmount: stripe.Int64(t.DefaultAmount),
			}
			if len(t.Tiers) > 0 {
				cop.Tiers = lo.Map(t.Tiers, func(t biz.PriceCurrencyOptionTier, _ int) *stripe.PriceCurrencyOptionsTierParams {
					return &stripe.PriceCurrencyOptionsTierParams{
						FlatAmount: stripe.Int64(t.FlatAmount),
						UnitAmount: stripe.Int64(t.UnitAmount),
						UpTo:       stripe.Int64(t.UpTo),
					}
				})
			}
			return strings.ToLower(t.CurrencyCode), cop
		})
	}

	if price.Recurring != nil {
		r.Recurring = &stripe.PriceRecurringParams{
			AggregateUsage:  stripe2.String(string(price.Recurring.AggregateUsage)),
			Interval:        stripe2.String(string(price.Recurring.Interval)),
			IntervalCount:   stripe.Int64(price.Recurring.IntervalCount),
			TrialPeriodDays: stripe.Int64(price.Recurring.TrialPeriodDays),
			UsageType:       stripe2.String(string(price.Recurring.UsageType)),
		}
	}
	if len(price.Tiers) > 0 {
		r.Tiers = lo.Map(price.Tiers, func(t biz.PriceTier, _ int) *stripe.PriceTierParams {
			return &stripe.PriceTierParams{
				FlatAmount: stripe.Int64(t.FlatAmount),
				UnitAmount: stripe.Int64(t.UnitAmount),
				UpTo:       stripe.Int64(t.UpTo),
			}
		})
	}

	return r
}

func mapBizPrice2UpdateStripe(stripeProductId string, price *biz.Price) *stripe.PriceParams {
	//https://github.com/stripe/stripe-node/issues/916
	//https://stripe.com/docs/api/prices/update
	r := &stripe.PriceParams{
		LookupKey: stripe2.String(price.ID.String()),
	}
	r.CurrencyOptions = map[string]*stripe.PriceCurrencyOptionsParams{}
	if len(price.CurrencyOptions) > 0 {
		r.CurrencyOptions = lo.SliceToMap(price.CurrencyOptions, func(t biz.PriceCurrencyOption) (string, *stripe.PriceCurrencyOptionsParams) {
			cop := &stripe.PriceCurrencyOptionsParams{
				UnitAmount: stripe.Int64(t.GetNeedPayAmount()),
			}
			if len(t.Tiers) > 0 {
				cop.Tiers = lo.Map(t.Tiers, func(t biz.PriceCurrencyOptionTier, _ int) *stripe.PriceCurrencyOptionsTierParams {
					return &stripe.PriceCurrencyOptionsTierParams{
						FlatAmount: stripe.Int64(t.FlatAmount),
						UnitAmount: stripe.Int64(t.UnitAmount),
						UpTo:       stripe.Int64(t.UpTo),
					}
				})
			}
			return strings.ToLower(t.CurrencyCode), cop
		})
	}
	return r
}
