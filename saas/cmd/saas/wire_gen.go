// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/auth/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	uow2 "github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/biz"
	conf2 "github.com/goxiaoy/go-saas-kit/saas/internal_/conf"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/data"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/server"
	"github.com/goxiaoy/go-saas-kit/saas/internal_/service"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(services *conf.Services, confData *conf2.Data, logger log.Logger, tokenizerConfig *jwt.TokenizerConfig, config *uow.Config, gormConfig *gorm.Config) (*kratos.App, func(), error) {
	tokenizer := jwt.NewTokenizer(tokenizerConfig)
	tenantRepo := data.NewTenantRepo()
	tenantStore := data.NewTenantStore(tenantRepo)
	dbOpener, cleanup := gorm.NewDbOpener()
	manager := uow2.NewUowManager(gormConfig, config, dbOpener)
	tenantUseCase := biz.NewTenantUserCase(tenantRepo)
	tenantService := service.NewTenantService(tenantUseCase)
	httpServer := server.NewHTTPServer(services, tokenizer, tenantStore, manager, tenantService, logger)
	grpcServer := server.NewGRPCServer(services, tokenizer, tenantStore, manager, tenantService, logger)
	dbProvider := data.NewProvider(confData, gormConfig, dbOpener, manager, tenantStore, logger)
	dataData, cleanup2, err := data.NewData(confData, dbProvider, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	migrate := data.NewMigrate(dataData)
	seeder := server.NewSeeder(confData, manager, migrate)
	app := newApp(logger, httpServer, grpcServer, seeder)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
