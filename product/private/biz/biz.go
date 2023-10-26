package biz

import (
	kitdi "github.com/go-saas/kit/pkg/di"
)

// ProviderSet is biz providers.
var ProviderSet = kitdi.NewSet(NewPostSeeder)

var (
	ProductMediaPath = "product/m"
)

type ProductManageProvider string

const (
	ProductManageProviderInternal ProductManageProvider = "internal"
	ProductManageProviderStripe   ProductManageProvider = "stripe"
)

const ConnName = "product"
