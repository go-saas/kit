package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "github.com/go-saas/kit/order/api/order/v1"
	v15 "github.com/go-saas/kit/payment/api/subscription/v1"
	"github.com/go-saas/kit/payment/private/biz"
	"github.com/go-saas/kit/pkg/authn"
	stripe2 "github.com/go-saas/kit/pkg/stripe"
	"github.com/go-saas/kit/pkg/utils"
	v13 "github.com/go-saas/kit/product/api/price/v1"
	v14 "github.com/go-saas/kit/product/api/product/v1"
	biz2 "github.com/go-saas/kit/product/private/biz"
	v12 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	"github.com/go-saas/saas/gorm"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v76"
	stripeclient "github.com/stripe/stripe-go/v76/client"

	pb "github.com/go-saas/kit/payment/api/checkout/v1"
)

type CheckoutService struct {
	orderInternalSrv v1.OrderInternalServiceServer
	userInternalSrv  v12.UserInternalServiceServer
	prodInternalSrv  v14.ProductInternalServiceServer
	subsRepo         biz.SubscriptionRepo
	stripeClient     *stripeclient.API
}

var _ pb.CheckoutServiceServer = (*CheckoutService)(nil)

func NewCheckoutService(
	orderInternalSrv v1.OrderInternalServiceServer,
	userInternalSrv v12.UserInternalServiceServer,
	prodInternalSrv v14.ProductInternalServiceServer,
	subsRepo biz.SubscriptionRepo,
	stripeClient *stripeclient.API,
) *CheckoutService {
	return &CheckoutService{
		orderInternalSrv: orderInternalSrv,
		userInternalSrv:  userInternalSrv,
		prodInternalSrv:  prodInternalSrv,
		subsRepo:         subsRepo,
		stripeClient:     stripeClient,
	}
}

func (s *CheckoutService) CheckoutNow(ctx context.Context, req *pb.CheckoutNowRequest) (*pb.CheckoutNowReply, error) {
	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	uid := ui.GetId()
	_, ok := lo.Find(req.Items, func(itemParams *pb.CheckoutItemParams) bool {
		return len(itemParams.PriceId) == 0
	})
	if ok {
		return nil, errors.BadRequest("", "")
	}
	// disable tenant filter to find products and prices from different tenant
	ctx = data.NewMultiTenancyDataFilter(ctx, false)
	ctx = data.NewAutoSetTenantId(ctx, false)
	var prices []*v13.Price
	tenantId := ""
	var priceType = ""
	var currencyCode = req.CurrencyCode
	for i, item := range req.Items {
		price, err := s.prodInternalSrv.GetInternalPrice(ctx, &v14.GetInternalPriceRequest{Id: item.PriceId})
		if err != nil {
			return nil, err
		}
		if len(currencyCode) == 0 {
			currencyCode = price.CurrencyCode
		}
		if i == 0 {
			tenantId = price.TenantId
			priceType = price.Type
		} else {
			if tenantId != price.TenantId {
				return nil, pb.ErrorProductAcrossTenantLocalized(ctx, nil, nil)
			}
			if priceType != price.Type {
				return nil, pb.ErrorPriceTypeUnsupportedLocalized(ctx, nil, nil)
			}
		}
		//check currency
		_, ok := lo.Find(price.CurrencyOptions, func(option *v13.PriceCurrencyOption) bool {
			return option.CurrencyCode == currencyCode
		})
		if price.CurrencyCode != currencyCode && !ok {
			return nil, pb.ErrorCurrencyUnsupportedLocalized(ctx, nil, nil)
		}
		prices = append(prices, price)
	}
	ret := &pb.CheckoutNowReply{}
	if priceType == string(biz2.PriceTypeRecurring) {
		subs, err := s.checkoutSubscription(ctx, req, prices, uid, tenantId)
		if err != nil {
			return nil, err
		}
		ret.Subscription = subs
	} else if priceType == string(biz2.PriceTypeOneTime) {
		order, err := s.checkoutOneTime(ctx, req, prices, uid, tenantId, currencyCode)
		if err != nil {
			return nil, err
		}
		ret.Order = order
	} else {
		return nil, pb.ErrorPriceTypeUnsupportedLocalized(ctx, nil, nil)
	}
	return ret, nil
}

func (s *CheckoutService) CheckoutOrder(ctx context.Context, req *pb.CheckOutOrderRequest) (*pb.CheckoutOrderReply, error) {
	ui, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.orderInternalSrv.GetInternalOrder(ctx, &v1.GetInternalOrderRequest{Id: &req.OrderId})
	if err != nil {
		return nil, err
	}
	if order.CustomerId != ui.GetId() {
		return nil, errors.NotFound("", "")
	}
	info, err := s.checkoutOrderInternal(ctx, order, req.Provider)
	if err != nil {
		return nil, err
	}

	return &pb.CheckoutOrderReply{PaymentProviderInfo: info}, nil
}

func (s *CheckoutService) checkoutOrderInternal(ctx context.Context, order *v1.Order, provider string) (*v1.OrderPaymentProviderInfo, error) {
	//TODO other providers
	customer, err := s.userInternalSrv.FindOrCreateStripeCustomer(ctx, &v12.FindOrCreateStripeCustomerRequest{
		UserId: &order.CustomerId,
	})
	if err != nil {
		return nil, err
	}

	var intent *stripe.PaymentIntent
	stripeInfo, ok := lo.Find(order.PaymentProviders, func(provider *v1.OrderPaymentProvider) bool {
		return provider.Provider == stripe2.ProviderName
	})
	if ok {
		//find intent
		intent, err = s.stripeClient.PaymentIntents.Get(stripeInfo.ProviderKey, nil)
		if err != nil {
			return nil, handleStripeError(err)
		}
	} else {
		//create intent
		//TODO find or create payment intent
		paymentIntentParams := &stripe.PaymentIntentParams{
			Amount:   &order.TotalPrice.Amount,
			Currency: &order.TotalPrice.CurrencyCode,
			Customer: &customer.StripeCustomerId,
		}
		paymentIntentParams.Metadata = map[string]string{
			"user_id":   order.CustomerId,
			"order_id":  order.Id,
			"tenant_id": order.TenantId,
		}
		intent, err = s.stripeClient.PaymentIntents.New(paymentIntentParams)
		if err != nil {
			return nil, handleStripeError(err)
		}
		//update payment provider
		_, err = s.orderInternalSrv.UpdateInternalOrderPaymentProvider(ctx, &v1.UpdateInternalOrderPaymentProviderRequest{
			OrderId:     order.Id,
			Provider:    stripe2.ProviderName,
			ProviderKey: intent.ID,
		})
		if err != nil {
			return nil, err
		}
	}
	ephemeralKey, err := s.stripeClient.EphemeralKeys.New(&stripe.EphemeralKeyParams{
		Customer:      &customer.StripeCustomerId,
		StripeVersion: stripe.String(stripe.APIVersion),
	})
	stripePaymentIntent := &stripe2.PaymentIntent{
		Id:           intent.ID,
		ClientSecret: intent.ClientSecret,
		Status:       string(intent.Status),
	}
	key := &stripe2.EphemeralKey{}
	stripe2.MapStripeEphemeralKey(ephemeralKey, key)
	return &v1.OrderPaymentProviderInfo{
		Stripe: &v1.OrderPaymentStripeInfo{
			PaymentIntent: stripePaymentIntent,
			EphemeralKey:  key,
		}}, nil
}

func (s *CheckoutService) checkoutSubscription(ctx context.Context,
	req *pb.CheckoutNowRequest, prices []*v13.Price, userId, tenantId string) (*v15.Subscription, error) {
	//TODO other providers
	customer, err := s.userInternalSrv.FindOrCreateStripeCustomer(
		ctx, &v12.FindOrCreateStripeCustomerRequest{UserId: &userId})
	if err != nil {
		return nil, err
	}
	subsNewParams := &stripe.SubscriptionParams{
		Customer: stripe2.String(customer.StripeCustomerId),
		Items: lo.Map(prices, func(t *v13.Price, i int) *stripe.SubscriptionItemsParams {
			return &stripe.SubscriptionItemsParams{Price: t.StripePriceId, Quantity: stripe2.Int64(int64(req.Items[i].Quantity))}
		}),
		PaymentBehavior: stripe.String("default_incomplete"),
		//TrialPeriodDays:            nil,
	}
	subsNewParams.AddExpand("latest_invoice.payment_intent")
	subs, err := s.stripeClient.Subscriptions.New(subsNewParams)
	if err != nil {
		return nil, err
	}
	localSubs := &biz.Subscription{Provider: stripe2.ProviderName, ProviderKey: subs.ID, UserId: userId}
	localSubs.Items = lo.Map(req.Items, func(t *pb.CheckoutItemParams, i int) biz.SubscriptionItem {
		return biz.SubscriptionItem{
			PriceID:        t.PriceId,
			PriceOwnerID:   prices[i].OwnerId,
			PriceOwnerType: prices[i].OwnerType,
			ProductID:      prices[i].ProductId,
			Quantity:       t.Quantity,
			BizPayload:     utils.Structpb2Map(t.BizPayload),
		}
	})
	localSubs.TenantId = gorm.NewTenantId(tenantId)
	MapStripeSubscription2Biz(subs, localSubs)
	err = s.subsRepo.Create(ctx, localSubs)
	if err != nil {
		return nil, err
	}
	ret := &v15.Subscription{}
	mapBizSubscription2Pb(localSubs, ret)
	infoSubs := &stripe2.Subscription{}
	stripe2.MapStripeSubscription(subs, infoSubs)
	ret.ProviderInfo = &v15.SubscriptionProviderInfo{Stripe: &v15.SubscriptionStripeInfo{
		Subscription: infoSubs,
	}}
	return ret, nil
}

func (s *CheckoutService) checkoutOneTime(ctx context.Context,
	req *pb.CheckoutNowRequest, prices []*v13.Price, userId, tenantId string, currencyCode string) (*v1.Order, error) {
	//create order and request payment
	orderRequest := &v1.CreateInternalOrderRequest{
		CustomerId:   userId,
		CurrencyCode: currencyCode,
	}
	// disable tenant filter to find products and prices from different tenant
	ctx = data.NewMultiTenancyDataFilter(ctx, false)

	var orderItems []*v1.CreateInternalOrderItem
	//cal prices
	for i, item := range req.Items {
		price := prices[i]
		var priceAmount, originalPriceAmount int64
		if price.CurrencyCode == currencyCode {
			if price.Discounted != nil {
				priceAmount = price.Discounted.Amount
			} else {
				priceAmount = price.Default.Amount
			}
			originalPriceAmount = price.Default.Amount
		} else {
			//find from currency options
			currencyOption, _ := lo.Find(price.CurrencyOptions, func(option *v13.PriceCurrencyOption) bool {
				return option.CurrencyCode == currencyCode
			})
			if currencyOption.Discounted != nil {
				priceAmount = currencyOption.Discounted.Amount
			} else {
				priceAmount = currencyOption.Default.Amount
			}
			originalPriceAmount = currencyOption.Default.Amount
		}
		//find product by id
		product, err := s.prodInternalSrv.GetInternalProduct(ctx, &v14.GetInternalProductRequest{Id: price.ProductId})
		if err != nil {
			return nil, err
		}
		mainPic := ""
		if product.MainPic != nil {
			mainPic = product.MainPic.Url
		}
		oi := &v1.CreateInternalOrderItem{
			Qty:                 item.Quantity,
			PriceAmount:         priceAmount,
			OriginalPriceAmount: originalPriceAmount,
			IsGiveaway:          false,
			BizPayload:          item.BizPayload,
			Product: &v1.OrderProduct{
				Name:    product.Title,
				MainPic: mainPic,
				Id:      &price.ProductId,
				Version: product.Version,
				PriceId: &price.Id,
			},
		}
		orderItems = append(orderItems, oi)
	}
	orderRequest.Items = orderItems
	thisTenantCtx := saas.NewCurrentTenant(ctx, tenantId, "")
	order, err := s.orderInternalSrv.CreateInternalOrder(thisTenantCtx, orderRequest)
	if err != nil {
		return nil, err
	}
	info, err := s.checkoutOrderInternal(ctx, order, req.Provider)
	//TODO requery order?
	if err != nil {
		return nil, err
	}
	order.PaymentProviderInfo = info
	return order, nil
}
