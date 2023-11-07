package service

import (
	"context"
	"github.com/go-saas/kit/pkg/price"
	v12 "github.com/go-saas/kit/product/api/price/v1"
	"github.com/go-saas/kit/product/private/biz"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapBizPrice2Pb(ctx context.Context, a *biz.Price, b *v12.Price) {
	b.Id = a.ID.String()

	b.CreatedAt = timestamppb.New(a.CreatedAt)
	b.UpdatedAt = timestamppb.New(a.UpdatedAt)
	b.TenantId = a.TenantId.String
	b.CurrencyCode = a.CurrencyCode
	b.ProductId = a.ProductID
	b.OwnerId = a.OwnerID
	b.OwnerType = a.OwnerType

	b.Default = price.MustNewFromInt64(a.DefaultAmount, a.CurrencyCode).ToPricePb(ctx)
	if a.DiscountedAmount != nil {
		b.Discounted = price.MustNewFromInt64(*a.DiscountedAmount, a.CurrencyCode).ToPricePb(ctx)
	}
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
		mapBizPriceTier2Pb(ctx, a.CurrencyCode, &t, r)
		return r
	})
	b.TiersMode = string(a.TiersMode)
	b.TransformQuantity = &v12.PriceTransformQuantity{}
	mapBizPriceTransformQuantity2Pb(&a.TransformQuantity, b.TransformQuantity)
	b.Type = string(a.Type)
	b.StripePriceId = a.StripePriceId
}

func mapBizCurrencyOption2Pb(ctx context.Context, a *biz.PriceCurrencyOption, b *v12.PriceCurrencyOption) {
	b.Default = price.MustNewFromInt64(a.DefaultAmount, a.CurrencyCode).ToPricePb(ctx)
	if a.DiscountedAmount != nil {
		b.Discounted = price.MustNewFromInt64(*a.DiscountedAmount, a.CurrencyCode).ToPricePb(ctx)
	}
	b.DiscountText = a.DiscountText
	b.DenyMoreDiscounts = a.DenyMoreDiscounts
	b.CurrencyCode = a.CurrencyCode
	b.Tiers = lo.Map(a.Tiers, func(t biz.PriceCurrencyOptionTier, i int) *v12.PriceCurrencyOptionTier {
		r := &v12.PriceCurrencyOptionTier{}
		mapBizPriceCurrencyOptionTier2Pb(ctx, a.CurrencyCode, &t, r)
		return r
	})
}

func mapBizPriceCurrencyOptionTier2Pb(ctx context.Context, currency string, a *biz.PriceCurrencyOptionTier, b *v12.PriceCurrencyOptionTier) {
	b.Flat = price.MustNewFromInt64(a.FlatAmount, currency).ToPricePb(ctx)
	b.Unit = price.MustNewFromInt64(a.UnitAmount, currency).ToPricePb(ctx)
	b.UpTo = a.UpTo
}

func mapBizPriceRecurring2Pb(a *biz.PriceRecurring, b *v12.PriceRecurring) {
	b.Interval = string(a.Interval)
	b.IntervalCount = a.IntervalCount
	b.TrialPeriodDays = a.TrialPeriodDays
	b.AggregateUsage = string(a.AggregateUsage)
	b.UsageType = string(a.UsageType)
}

func mapBizPriceTier2Pb(ctx context.Context, currency string, a *biz.PriceTier, b *v12.PriceTier) {
	b.Flat = price.MustNewFromInt64(a.FlatAmount, currency).ToPricePb(ctx)
	b.Unit = price.MustNewFromInt64(a.UnitAmount, currency).ToPricePb(ctx)
	b.UpTo = a.UpTo
}

func mapBizPriceTransformQuantity2Pb(a *biz.PriceTransformQuantity, b *v12.PriceTransformQuantity) {
	b.DivideBy = a.DivideBy
	b.Round = string(a.Round)
}

func mapPbPrice2Biz(a *v12.PriceParams, b *biz.Price) {
	if len(a.Id) > 0 {
		b.UIDBase.ID = uuid.MustParse(a.Id)
	}
	b.CurrencyCode = a.CurrencyCode
	b.DefaultAmount = price.MustNew(a.DefaultAmountDecimal, a.CurrencyCode).Amount
	if a.DiscountedAmountDecimal != nil {
		dis := price.MustNew(*a.DiscountedAmountDecimal, a.CurrencyCode).Amount
		b.DiscountedAmount = &dis
	}

	b.DiscountText = a.DiscountText
	b.DenyMoreDiscounts = a.DenyMoreDiscounts

	b.BillingScheme = biz.PriceBillingScheme(a.BillingScheme)
	b.CurrencyOptions = lo.Map(a.CurrencyOptions, func(t *v12.PriceCurrencyOptionParams, i int) biz.PriceCurrencyOption {
		r := &biz.PriceCurrencyOption{}
		mapPbCurrencyOption2Biz(t, r)
		return *r
	})

	if a.Recurring != nil {
		b.Recurring = &biz.PriceRecurring{}
		mapPbPriceRecurring2Biz(a.Recurring, b.Recurring)
	}
	b.Tiers = lo.Map(a.Tiers, func(t *v12.PriceTierParams, i int) biz.PriceTier {
		r := &biz.PriceTier{}
		mapPbPriceTier2Biz(a.CurrencyCode, t, r)
		return *r
	})
	b.TiersMode = biz.PriceTiersMode(a.TiersMode)
	if a.TransformQuantity != nil {
		b.TransformQuantity = biz.PriceTransformQuantity{}
		mapPbPriceTransformQuantity2Biz(a.TransformQuantity, &b.TransformQuantity)
	}
	b.Type = biz.PriceType(a.Type)
}

func mapPbUpdatePrice2Biz(a *v12.PriceParams, b *biz.Price) {
	if len(a.Id) > 0 {
		b.UIDBase.ID = uuid.MustParse(a.Id)
	}

	b.DiscountText = a.DiscountText
	b.DenyMoreDiscounts = a.DenyMoreDiscounts

	b.CurrencyOptions = lo.Map(a.CurrencyOptions, func(t *v12.PriceCurrencyOptionParams, i int) biz.PriceCurrencyOption {
		r := &biz.PriceCurrencyOption{}
		mapPbCurrencyOption2Biz(t, r)
		return *r
	})

}

func mapPbCurrencyOption2Biz(a *v12.PriceCurrencyOptionParams, b *biz.PriceCurrencyOption) {
	b.DefaultAmount = price.MustNew(a.DefaultAmountDecimal, a.CurrencyCode).Amount
	if a.DiscountedAmountDecimal != nil {
		dis := price.MustNew(*a.DiscountedAmountDecimal, a.CurrencyCode).Amount
		b.DiscountedAmount = &dis
	}
	b.DiscountText = a.DiscountText
	b.DenyMoreDiscounts = a.DenyMoreDiscounts
	b.CurrencyCode = a.CurrencyCode
	b.Tiers = lo.Map(a.Tiers, func(t *v12.PriceCurrencyOptionTierParams, i int) biz.PriceCurrencyOptionTier {
		r := &biz.PriceCurrencyOptionTier{}
		mapPbPriceCurrencyOptionTier2Biz(a.CurrencyCode, t, r)
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

func mapPbPriceCurrencyOptionTier2Biz(currencyCode string, a *v12.PriceCurrencyOptionTierParams, b *biz.PriceCurrencyOptionTier) {
	b.FlatAmount = price.MustNew(a.FlatAmountDecimal, currencyCode).Amount
	b.UnitAmount = price.MustNew(a.UnitAmountDecimal, currencyCode).Amount
	b.UpTo = a.UpTo
}

func mapPbPriceTier2Biz(currencyCode string, a *v12.PriceTierParams, b *biz.PriceTier) {
	b.FlatAmount = price.MustNew(a.FlatAmountDecimal, currencyCode).Amount
	b.UnitAmount = price.MustNew(a.UnitAmountDecimal, currencyCode).Amount
	b.UpTo = a.UpTo
}
func mapPbPriceTransformQuantity2Biz(a *v12.PriceTransformQuantity, b *biz.PriceTransformQuantity) {
	b.DivideBy = a.DivideBy
	b.Round = biz.PriceTransformQuantityRound(a.Round)
}
