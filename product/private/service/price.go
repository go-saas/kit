package service

import (
	"context"
	"github.com/go-saas/kit/pkg/price"
	v12 "github.com/go-saas/kit/product/api/price/v1"
	"github.com/go-saas/kit/product/private/biz"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapBizPrice2Pb(ctx context.Context, a *biz.Price, b *v12.Price) {
	b.Id = a.ID.String()

	b.CreatedAt = timestamppb.New(a.CreatedAt)
	b.UpdatedAt = timestamppb.New(a.UpdatedAt)
	b.TenantId = a.TenantId.String

	b.Default = price.MustNewFromInt64(a.DefaultAmount, a.CurrencyCode).ToPricePb(ctx)
	b.Discounted = price.MustNewFromInt64(a.DiscountedAmount, a.CurrencyCode).ToPricePb(ctx)
	b.DiscountText = a.DiscountText
	b.DenyMoreDiscounts = a.DenyMoreDiscounts

	b.BillingScheme = string(a.BillingScheme)
	b.CurrencyOptions = lo.Map(a.CurrencyOptions, func(t biz.PriceCurrencyOption, i int) *v12.PriceCurrencyOption {
		r := &v12.PriceCurrencyOption{}
		mapBizCurrencyOption2Pb(ctx, &t, r)
		return r
	})

	if a.Recurring != nil {
		b.Recurring = &v12.PriceRecurring{}
		mapBizPriceRecurring2Pb(a.Recurring, b.Recurring)
	}
	b.Tiers = lo.Map(a.Tiers, func(t biz.PriceTier, i int) *v12.PriceTier {
		r := &v12.PriceTier{}
		mapBizPriceTier2Pb(&t, r)
		return r
	})
	b.TiersMode = string(a.TiersMode)
	b.TransformQuantity = &v12.PriceTransformQuantity{}
	mapBizPriceTransformQuantity2Pb(&a.TransformQuantity, b.TransformQuantity)
	b.Type = string(a.Type)
}

func mapBizCurrencyOption2Pb(ctx context.Context, a *biz.PriceCurrencyOption, b *v12.PriceCurrencyOption) {
	b.Default = price.MustNewFromInt64(a.DefaultAmount, a.CurrencyCode).ToPricePb(ctx)
	b.Discounted = price.MustNewFromInt64(a.DiscountedAmount, a.CurrencyCode).ToPricePb(ctx)
	b.DiscountText = a.DiscountText
	b.DenyMoreDiscounts = a.DenyMoreDiscounts
	b.CurrencyCode = a.CurrencyCode
	b.Tiers = lo.Map(a.Tiers, func(t biz.PriceCurrencyOptionTier, i int) *v12.PriceCurrencyOptionTier {
		r := &v12.PriceCurrencyOptionTier{}
		mapBizPriceCurrencyOptionTier2Pb(&t, r)
		return r
	})
}

func mapBizPriceCurrencyOptionTier2Pb(a *biz.PriceCurrencyOptionTier, b *v12.PriceCurrencyOptionTier) {
	b.FlatAmount = a.FlatAmount
	b.UnitAmount = a.UnitAmount
	b.UpTo = a.UpTo
}

func mapBizPriceRecurring2Pb(a *biz.PriceRecurring, b *v12.PriceRecurring) {
	b.Interval = string(a.Interval)
	b.IntervalCount = a.IntervalCount
	b.TrialPeriodDays = a.TrialPeriodDays
	b.AggregateUsage = string(a.AggregateUsage)
	b.UsageType = string(a.UsageType)
}

func mapBizPriceTier2Pb(a *biz.PriceTier, b *v12.PriceTier) {
	b.FlatAmount = a.FlatAmount
	b.UnitAmount = a.UnitAmount
	b.UpTo = a.UpTo
}

func mapBizPriceTransformQuantity2Pb(a *biz.PriceTransformQuantity, b *v12.PriceTransformQuantity) {
	b.DivideBy = a.DivideBy
	b.Round = string(a.Round)
}

func mapPbCreatePrice2Biz(a *v12.CreatePriceRequest, b *biz.Price) {

	b.DefaultAmount = a.DefaultAmount
	b.DiscountedAmount = a.DiscountedAmount

	b.DiscountText = a.DiscountText
	b.DenyMoreDiscounts = a.DenyMoreDiscounts

	b.BillingScheme = biz.PriceBillingScheme(a.BillingScheme)
	b.CurrencyOptions = lo.Map(a.CurrencyOptions, func(t *v12.PriceCurrencyOptionPrams, i int) biz.PriceCurrencyOption {
		r := &biz.PriceCurrencyOption{}
		mapPbCurrencyOption2Biz(t, r)
		return *r
	})

	if a.Recurring != nil {
		b.Recurring = &biz.PriceRecurring{}
		mapPbPriceRecurring2Biz(a.Recurring, b.Recurring)
	}
	b.Tiers = lo.Map(a.Tiers, func(t *v12.PriceTier, i int) biz.PriceTier {
		r := &biz.PriceTier{}
		mapPbPriceTier2Biz(t, r)
		return *r
	})
	b.TiersMode = biz.PriceTiersMode(a.TiersMode)
	b.TransformQuantity = biz.PriceTransformQuantity{}
	mapPbPriceTransformQuantity2Biz(a.TransformQuantity, &b.TransformQuantity)
	b.Type = biz.PriceType(a.Type)
}

func mapPbUpdatePrice2Biz(a *v12.UpdatePrice, b *biz.Price) {

	b.DefaultAmount = a.DefaultAmount
	b.DiscountedAmount = a.DiscountedAmount

	b.DiscountText = a.DiscountText
	b.DenyMoreDiscounts = a.DenyMoreDiscounts

	b.BillingScheme = biz.PriceBillingScheme(a.BillingScheme)
	b.CurrencyOptions = lo.Map(a.CurrencyOptions, func(t *v12.PriceCurrencyOptionPrams, i int) biz.PriceCurrencyOption {
		r := &biz.PriceCurrencyOption{}
		mapPbCurrencyOption2Biz(t, r)
		return *r
	})

	if a.Recurring != nil {
		b.Recurring = &biz.PriceRecurring{}
		mapPbPriceRecurring2Biz(a.Recurring, b.Recurring)
	}
	b.Tiers = lo.Map(a.Tiers, func(t *v12.PriceTier, i int) biz.PriceTier {
		r := &biz.PriceTier{}
		mapPbPriceTier2Biz(t, r)
		return *r
	})
	b.TiersMode = biz.PriceTiersMode(a.TiersMode)
	b.TransformQuantity = biz.PriceTransformQuantity{}
	mapPbPriceTransformQuantity2Biz(a.TransformQuantity, &b.TransformQuantity)
	b.Type = biz.PriceType(a.Type)
}

func mapPbCurrencyOption2Biz(a *v12.PriceCurrencyOptionPrams, b *biz.PriceCurrencyOption) {
	b.DefaultAmount = a.DefaultAmount
	b.DiscountedAmount = a.DiscountedAmount
	b.DiscountText = a.DiscountText
	b.DenyMoreDiscounts = a.DenyMoreDiscounts
	b.CurrencyCode = a.CurrencyCode
	b.Tiers = lo.Map(a.Tiers, func(t *v12.PriceCurrencyOptionTier, i int) biz.PriceCurrencyOptionTier {
		r := &biz.PriceCurrencyOptionTier{}
		mapPbPriceCurrencyOptionTier2Biz(t, r)
		return *r
	})
}

func mapPbPriceRecurring2Biz(a *v12.PriceRecurring, b *biz.PriceRecurring) {
	b.Interval = biz.PriceRecurringInterval(a.Interval)
	b.IntervalCount = a.IntervalCount
	b.TrialPeriodDays = a.TrialPeriodDays
	b.AggregateUsage = biz.PriceRecurringAggregateUsage(a.AggregateUsage)
	b.UsageType = biz.PriceRecurringUsageType(a.UsageType)
}

func mapPbPriceCurrencyOptionTier2Biz(a *v12.PriceCurrencyOptionTier, b *biz.PriceCurrencyOptionTier) {
	b.FlatAmount = a.FlatAmount
	b.UnitAmount = a.UnitAmount
	b.UpTo = a.UpTo
}
func mapPbPriceTier2Biz(a *v12.PriceTier, b *biz.PriceTier) {
	b.FlatAmount = a.FlatAmount
	b.UnitAmount = a.UnitAmount
	b.UpTo = a.UpTo
}
func mapPbPriceTransformQuantity2Biz(a *v12.PriceTransformQuantity, b *biz.PriceTransformQuantity) {
	b.DivideBy = a.DivideBy
	b.Round = biz.PriceTransformQuantityRound(a.Round)
}
