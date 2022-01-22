// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authorization"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/casbin"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	uow2 "github.com/goxiaoy/go-saas-kit/pkg/uow"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
	conf2 "github.com/goxiaoy/go-saas-kit/user/private/conf"
	"github.com/goxiaoy/go-saas-kit/user/private/data"
	"github.com/goxiaoy/go-saas-kit/user/private/seed"
	"github.com/goxiaoy/go-saas-kit/user/private/server"
	"github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(services *conf.Services, security *conf.Security, userConf *conf2.UserConf, confData *conf2.Data, logger log.Logger, passwordValidatorConfig *biz.PasswordValidatorConfig, config *uow.Config, gormConfig *gorm.Config, webMultiTenancyOption *http.WebMultiTenancyOption) (*kratos.App, func(), error) {
	tokenizerConfig := jwt.NewTokenizerConfig(security)
	tokenizer := jwt.NewTokenizer(tokenizerConfig)
	dbOpener, cleanup := gorm.NewDbOpener()
	manager := uow2.NewUowManager(gormConfig, config, dbOpener)
	saasContributor := api.NewSaasContributor(webMultiTenancyOption)
	userContributor := api.NewUserContributor()
	option := api.NewDefaultOption(saasContributor, userContributor)
	tenantStore := data.NewTenantStore()
	sessionStorer := server.NewSessionStorer(security, userConf)
	cookieStorer := server.NewCookieStorer(security, userConf)
	dbProvider := data.NewProvider(confData, gormConfig, dbOpener, tenantStore, logger)
	dataData, cleanup2, err := data.NewData(confData, dbProvider, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData)
	passwordHasher := biz.NewPasswordHasher()
	userValidator := biz.NewUserValidator()
	passwordValidator := biz.NewPasswordValidator(passwordValidatorConfig)
	lookupNormalizer := biz.NewLookupNormalizer()
	userManager := biz.NewUserManager(userRepo, passwordHasher, userValidator, passwordValidator, lookupNormalizer, logger)
	userTokenRepo := data.NewUserTokenRepo(dataData)
	authbossStoreWrapper := biz.NewAuthbossStoreWrapper(userManager, userTokenRepo)
	authboss, err := server.NewAuthboss(logger, userConf, sessionStorer, cookieStorer, authbossStoreWrapper)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	enforcerProvider := data.NewEnforcerProvider(dbProvider)
	permissionService := casbin.NewPermissionService(enforcerProvider)
	userRoleContributor := service.NewUserRoleContributor(userRepo)
	authorizationOption := service.NewAuthorizationOption(userRoleContributor)
	subjectResolverImpl := authorization.NewSubjectResolver(authorizationOption)
	defaultAuthorizationService := authorization.NewDefaultAuthorizationService(permissionService, subjectResolverImpl, logger)
	userService := service.NewUserService(userManager, defaultAuthorizationService)
	accountService := service.NewAccountService(userManager)
	roleRepo := data.NewRoleRepo(dataData)
	roleManager := biz.NewRoleManager(roleRepo, lookupNormalizer)
	refreshTokenRepo := data.NewRefreshTokenRepo(dataData)
	authService := service.NewAuthService(userManager, roleManager, tokenizer, tokenizerConfig, passwordValidator, refreshTokenRepo, security)
	roleService := service.NewRoleServiceService(roleRepo, defaultAuthorizationService)
	servicePermissionService := service.NewPermissionService(defaultAuthorizationService, permissionService, subjectResolverImpl)
	httpServer := server.NewHTTPServer(services, security, tokenizer, manager, webMultiTenancyOption, option, tenantStore, logger, authboss, userService, accountService, authService, roleService, servicePermissionService)
	grpcServer := server.NewGRPCServer(services, tokenizer, tenantStore, manager, webMultiTenancyOption, option, logger, userService, accountService, authService, roleService, servicePermissionService)
	migrate := data.NewMigrate(dataData)
	roleSeed := biz.NewRoleSeed(roleManager, permissionService)
	userSeed := biz.NewUserSeed(userManager, roleManager)
	fake := seed.NewFake(userManager)
	permissionSeeder := biz.NewPermissionSeeder(permissionService, permissionService, roleManager)
	seeder := server.NewSeeder(userConf, manager, migrate, roleSeed, userSeed, fake, permissionSeeder)
	app := newApp(logger, httpServer, grpcServer, seeder)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
