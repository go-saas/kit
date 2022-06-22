package biz

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/pkg/dal"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewMenuSeed)

const (
	SeedPathKey              = "seed.menu.path"
	ConnName    dal.ConnName = "sys"
)
