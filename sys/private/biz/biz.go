package biz

import (
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
)

// ProviderSet is biz providers.
var ProviderSet = kitdi.NewSet(NewMenuSeed, NewApisixSeed, NewApisixMigrationTaskHandler)

const (
	ConnName dal.ConnName = "sys"
)
