package biz

import (
	"github.com/go-saas/kit/pkg/dal"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewMenuSeed, wire.Struct(new(ApisixSeed), "*"))

const (
	SeedPathKey              = "seed.menu.path"
	ConnName    dal.ConnName = "sys"
)
