package biz

import (
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/goava/di"
)

// ProviderSet is biz providers.
var ProviderSet = kitdi.NewSet(
	NewTenantUserCase,
	kitdi.NewProvider(NewTenantReadyEventHandler, di.As(new(event.ConsumerHandler))),
	kitdi.NewProvider(NewSubscriptionChangedEventHandler, di.As(new(event.ConsumerHandler))),
	kitdi.NewProvider(NewOrderChangedEventHandler, di.As(new(event.ConsumerHandler))),
	NewConfigConnStrGenerator,
)

const ConnName dal.ConnName = "saas"

const TenantLogoPath = "saas/tenant/logo"
