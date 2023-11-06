package stripe

import (
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/stripe/stripe-go/v76"
	stripeclient "github.com/stripe/stripe-go/v76/client"
)

var ProviderSet = kitdi.NewSet(NewStripeClient)

func NewStripeClient(c *Conf) *stripeclient.API {
	if c == nil {
		return nil
	}
	sc := &stripeclient.API{}
	sc.Init(c.PrivateKey, nil)
	return sc
}

// String returns a pointer to the string value passed in.
func String(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func Int64(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

const (
	ProviderName = "stripe"
)

func MapStripeSubscription(a *stripe.Subscription, b *Subscription) {
	b.Id = a.ID
	if a.LatestInvoice != nil {
		b.LatestInvoice = &Invoice{}
		MapStripeInvoice(a.LatestInvoice, b.LatestInvoice)
	}
}

func MapStripeInvoice(a *stripe.Invoice, b *Invoice) {
	b.Id = a.ID
	if a.PaymentIntent != nil {
		b.PaymentIntent = &PaymentIntent{}
		MapStripePaymentIntent(a.PaymentIntent, b.PaymentIntent)
	}
}

func MapStripePaymentIntent(a *stripe.PaymentIntent, b *PaymentIntent) {
	b.Id = a.ID
	b.ClientSecret = a.ClientSecret
	b.Status = string(a.Status)
}

func MapStripeEphemeralKey(a *stripe.EphemeralKey, b *EphemeralKey) {
	b.Secret = a.Secret
}
