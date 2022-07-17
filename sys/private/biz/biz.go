package biz

import (
	"github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
)

// ProviderSet is biz providers.
var ProviderSet = kitdi.NewSet(NewMenuSeed, NewApisixSeed, NewApisixMigrationTaskHandler)

const (
	SeedPathKey              = "seed.menu.path"
	ConnName    dal.ConnName = "sys"
)
