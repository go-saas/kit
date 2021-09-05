//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/user/internal_/biz"
	"github.com/goxiaoy/go-saas-kit/user/internal_/conf"
	"github.com/goxiaoy/go-saas-kit/user/internal_/data"
	"github.com/goxiaoy/go-saas-kit/user/internal_/server"
	"github.com/goxiaoy/go-saas-kit/user/internal_/service"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
)

// initApp init kratos application.
func initApp(*conf2.Services, *conf.UserConf, *conf.Data, log.Logger, *biz.PasswordValidatorConfig, *jwt.TokenizerConfig, *uow.Config, *gorm.Config) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
