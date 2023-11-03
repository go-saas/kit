package service

import (
	"github.com/go-saas/kit/payment/private/biz"
	"github.com/stripe/stripe-go/v76"
	"strings"
	"time"
)

func MapStripeSubscription2Biz(a *stripe.Subscription, b *biz.Subscription) {
	b.Status = biz.SubscriptionStatus(a.Status)
	b.CancelAtPeriodEnd = a.CancelAtPeriodEnd
	b.CurrencyCode = strings.ToUpper(string(a.Currency))
	startTime := time.Unix(a.CurrentPeriodStart, 0)
	b.CurrentPeriodStart = &startTime
	endTime := time.Unix(a.CurrentPeriodEnd, 0)
	b.CurrentPeriodEnd = &endTime

}
