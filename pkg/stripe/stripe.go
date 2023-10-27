package stripe

import (
	kitdi "github.com/go-saas/kit/pkg/di"
	stripeclient "github.com/stripe/stripe-go/v76/client"
)

var ProviderSet = kitdi.NewSet(NewStripeClient)

func NewStripeClient(c *StripeConf) *stripeclient.API {
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
