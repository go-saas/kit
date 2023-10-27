package service

import (
	"context"
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/pkg/job"
	"github.com/go-saas/kit/product/private/biz"
	"github.com/go-saas/saas"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v76"
	stripeclient "github.com/stripe/stripe-go/v76/client"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

const (
	JobTypeProductUpdated = "product" + ":" + "product_updated"
)

func NewProductUpdatedTask(prams *biz.ProductUpdatedJobPram) (*asynq.Task, error) {
	payload, err := protojson.Marshal(prams)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(JobTypeProductUpdated, payload, asynq.Queue(string(biz.ConnName))), nil
}

func NewProductUpdatedTaskHandler(repo biz.ProductRepo, client *stripeclient.API) *job.Handler {
	return job.NewHandlerFunc(JobTypeProductUpdated, func(ctx context.Context, t *asynq.Task) error {
		msg := &biz.ProductUpdatedJobPram{}
		if err := protojson.Unmarshal(t.Payload(), msg); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}
		//change to product tenant
		ctx = saas.NewCurrentTenant(ctx, msg.TenantId, "")
		product, err := repo.Get(ctx, msg.ProductId)
		if err != nil {
			return err
		}
		if product == nil {
			return fmt.Errorf("can not find %s: %w", msg.ProductId, asynq.SkipRetry)
		}
		if product.Version.String != msg.ProductVersion {
			return fmt.Errorf("product:%s version mismatch, should be:%s, got:%s: %w", msg.ProductId, msg.ProductVersion, product.Version.String, asynq.SkipRetry)
		}
		group, ctx := errgroup.WithContext(ctx)
		//sync with stripe
		group.Go(func() error {
			return syncWithStripe(ctx, client, repo, product, msg.ProductVersion)
		})
		return group.Wait()

	})
}

func syncWithStripe(ctx context.Context, client *stripeclient.API, repo biz.ProductRepo, product *biz.Product, version string) error {
	if client == nil {
		return nil
	}
	links, err := repo.GetSyncLinks(ctx, product)
	if err != nil {
		return err
	}
	stripeInfo, find := lo.Find(links, func(link biz.ProductSyncLink) bool {
		return string(biz.ProductManageProviderStripe) == link.ProviderName
	})
	var stripeId string
	if !find {
		// create stripe object
		params := mapBizProduct2Stripe(product)
		stripeProduct, err := client.Products.New(params)
		if err != nil {
			return err
		}
		stripeId = stripeProduct.ID
		t := time.Now()
		err = repo.UpdateSyncLink(ctx, product, &biz.ProductSyncLink{
			ProviderName: string(biz.ProductManageProviderStripe),
			ProviderId:   stripeId,
			LastSyncTime: &t,
		})
		if err != nil {
			return err
		}
	} else {
		stripeId = stripeInfo.ProviderId
		//update product if needed
		stripeProduct, err := client.Products.Get(stripeId, &stripe.ProductParams{})
		if stripeProduct.Metadata["version"] != version {
			klog.Infof("product_id:%s version:%s same with stipe, skip updates", product.ID.String(), version)
			return nil
		}
		params := mapBizProduct2Stripe(product)
		_, err = client.Products.Update(stripeId, params)
		if err != nil {
			return err
		}
	}

	//TODO update price

	return nil
}

func mapBizProduct2Stripe(product *biz.Product) *stripe.ProductParams {
	return &stripe.ProductParams{
		Active:      stripe.Bool(product.Active),
		Name:        stripe.String(product.Title),
		Description: stripe.String(product.Desc),
		Shippable:   stripe.Bool(product.NeedShipping),
		Metadata:    map[string]string{"version": product.Version.String},
		//TODO type

	}
}
