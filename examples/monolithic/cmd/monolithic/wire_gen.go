// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/goxiaoy/go-eventbus"
	server2 "github.com/goxiaoy/go-saas-kit/examples/monolithic/private/server"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/casbin"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/dal"
	"github.com/goxiaoy/go-saas-kit/pkg/gorm"
	"github.com/goxiaoy/go-saas-kit/pkg/job"
	"github.com/goxiaoy/go-saas-kit/pkg/redis"
	"github.com/goxiaoy/go-saas-kit/pkg/server"
	"github.com/goxiaoy/go-saas-kit/pkg/uow"
	api2 "github.com/goxiaoy/go-saas-kit/saas/api"
	"github.com/goxiaoy/go-saas-kit/saas/private/biz"
	conf2 "github.com/goxiaoy/go-saas-kit/saas/private/conf"
	"github.com/goxiaoy/go-saas-kit/saas/private/data"
	server5 "github.com/goxiaoy/go-saas-kit/saas/private/server"
	"github.com/goxiaoy/go-saas-kit/saas/private/service"
	biz3 "github.com/goxiaoy/go-saas-kit/sys/private/biz"
	data3 "github.com/goxiaoy/go-saas-kit/sys/private/data"
	server4 "github.com/goxiaoy/go-saas-kit/sys/private/server"
	service3 "github.com/goxiaoy/go-saas-kit/sys/private/service"
	api3 "github.com/goxiaoy/go-saas-kit/user/api"
	biz2 "github.com/goxiaoy/go-saas-kit/user/private/biz"
	conf3 "github.com/goxiaoy/go-saas-kit/user/private/conf"
	data2 "github.com/goxiaoy/go-saas-kit/user/private/data"
	server3 "github.com/goxiaoy/go-saas-kit/user/private/server"
	service2 "github.com/goxiaoy/go-saas-kit/user/private/service"
	"github.com/goxiaoy/go-saas-kit/user/private/service/http"
)

import (
	_ "github.com/goxiaoy/go-saas-kit/pkg/event/kafka"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(services *conf.Services, security *conf.Security, confData *conf.Data, saasConf *conf2.SaasConf, userConf *conf3.UserConf, logger log.Logger, appConfig *conf.AppConfig, arg ...grpc.ClientOption) (*kratos.App, func(), error) {
	tokenizerConfig := jwt.NewTokenizerConfig(security)
	tokenizer := jwt.NewTokenizer(tokenizerConfig)
	config := _wireConfigValue
	dbCache, cleanup := gorm.NewDbCache(confData, logger)
	manager := uow.NewUowManager(config, dbCache)
	webMultiTenancyOption := server.NewWebMultiTenancyOption(appConfig)
	option := api.NewDefaultOption(logger)
	trustedContextValidator := api.NewClientTrustedContextValidator()
	eventBus := _wireEventBusValue
	connStrings := dal.NewConstantConnStrResolver(confData)
	constDbProvider := dal.NewConstDbProvider(dbCache, connStrings, confData)
	dataData, cleanup2, err := data.NewData(confData, constDbProvider, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	tenantRepo := data.NewTenantRepo(eventBus, dataData)
	connStrGenerator := biz.NewConfigConnStrGenerator(saasConf)
	connName := _wireConnNameValue
	producer, cleanup3, err := dal.NewEventSender(confData, connName)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	tenantUseCase := biz.NewTenantUserCase(tenantRepo, connStrGenerator, producer)
	factory := dal.NewBlobFactory(confData)
	tenantInternalService := &service.TenantInternalService{
		Trusted: trustedContextValidator,
		UseCase: tenantUseCase,
		App:     appConfig,
		Blob:    factory,
	}
	tenantStore := api2.NewTenantStore(tenantInternalService)
	decodeRequestFunc := _wireDecodeRequestFuncValue
	encodeResponseFunc := _wireEncodeResponseFuncValue
	encodeErrorFunc := _wireEncodeErrorFuncValue
	connStrResolver := dal.NewConnStrResolver(confData, tenantStore)
	dbProvider := gorm.NewDbProvider(dbCache, connStrResolver, confData)
	data4, cleanup4, err := data2.NewData(confData, dbProvider, logger)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	userRepo := data2.NewUserRepo(data4)
	passwordHasher := biz2.NewPasswordHasher()
	userValidator := biz2.NewUserValidator()
	passwordValidator := biz2.NewPasswordValidator(userConf)
	lookupNormalizer := biz2.NewLookupNormalizer()
	userTokenRepo := data2.NewUserTokenRepo(data4)
	refreshTokenRepo := data2.NewRefreshTokenRepo(data4)
	userTenantRepo := data2.NewUserTenantRepo(data4)
	client, err := dal.NewRedis(confData, connName)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	emailTokenProvider := biz2.NewEmailTokenProvider(client)
	phoneTokenProvider := biz2.NewPhoneTokenProvider(client)
	cache := redis.NewCache(client)
	twoStepTokenProvider := biz2.NewTwoStepTokenProvider(cache)
	userManager := biz2.NewUserManager(userConf, userRepo, passwordHasher, userValidator, passwordValidator, lookupNormalizer, userTokenRepo, refreshTokenRepo, userTenantRepo, emailTokenProvider, phoneTokenProvider, twoStepTokenProvider, logger)
	roleRepo := data2.NewRoleRepo(data4, eventBus)
	roleManager := biz2.NewRoleManager(roleRepo, lookupNormalizer)
	enforcerProvider, err := data2.NewEnforcerProvider(logger, dbProvider)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	permissionService := casbin.NewPermissionService(enforcerProvider)
	userRoleContrib := service2.NewUserRoleContrib(userRepo)
	authzOption := server2.NewAuthorizationOption(userRoleContrib)
	subjectResolverImpl := authz.NewSubjectResolver(authzOption)
	defaultAuthorizationService := authz.NewDefaultAuthorizationService(permissionService, subjectResolverImpl, logger)
	userService := service2.NewUserService(userManager, roleManager, defaultAuthorizationService, factory, trustedContextValidator, logger)
	userTenantContrib := api3.NewUserTenantContrib(userService)
	lazyClient := dal.NewEmailer(confData)
	emailSender := biz2.NewEmailSender(lazyClient, confData)
	authService := service2.NewAuthService(userManager, roleManager, tokenizer, tokenizerConfig, passwordValidator, refreshTokenRepo, emailSender, security, defaultAuthorizationService, trustedContextValidator, logger)
	refreshTokenProvider := api3.NewRefreshProvider(authService, logger)
	tenantService := service.NewTenantService(tenantUseCase, defaultAuthorizationService, trustedContextValidator, factory, appConfig)
	userSettingRepo := data2.NewUserSettingRepo(data4, eventBus)
	userAddressRepo := data2.NewUserAddrRepo(data4, eventBus)
	accountService := service2.NewAccountService(userManager, factory, tenantService, userSettingRepo, userAddressRepo, lookupNormalizer)
	roleService := service2.NewRoleServiceService(roleManager, defaultAuthorizationService, permissionService)
	servicePermissionService := service2.NewPermissionService(defaultAuthorizationService, permissionService, subjectResolverImpl, trustedContextValidator)
	signInManager := biz2.NewSignInManager(userManager, security)
	apiClient := service2.NewHydra(security)
	auth := http.NewAuth(decodeRequestFunc, userManager, logger, signInManager, apiClient)
	httpServerRegister := service2.NewHttpServerRegister(userService, encodeResponseFunc, encodeErrorFunc, accountService, authService, roleService, servicePermissionService, auth, confData, defaultAuthorizationService, factory)
	menuRepo := data3.NewMenuRepo(dbProvider, eventBus)
	menuService := service3.NewMenuService(defaultAuthorizationService, menuRepo, logger)
	redisConnOpt := job.NewAsynqClientOpt(client)
	serviceHttpServerRegister := service3.NewHttpServerRegister(menuService, defaultAuthorizationService, encodeErrorFunc, factory, confData, redisConnOpt)
	httpServerRegister2 := service.NewHttpServerRegister(tenantService, factory, defaultAuthorizationService, encodeErrorFunc, tenantInternalService, confData)
	serverHttpServerRegister := server2.NewHttpServiceRegister(httpServerRegister, serviceHttpServerRegister, httpServerRegister2)
	httpServer := server2.NewHTTPServer(services, security, tokenizer, manager, webMultiTenancyOption, option, tenantStore, decodeRequestFunc, encodeResponseFunc, encodeErrorFunc, logger, userTenantContrib, trustedContextValidator, refreshTokenProvider, serverHttpServerRegister)
	grpcServerRegister := service2.NewGrpcServerRegister(userService, accountService, authService, roleService, servicePermissionService)
	serviceGrpcServerRegister := service3.NewGrpcServerRegister(menuService)
	grpcServerRegister2 := service.NewGrpcServerRegister(tenantService, tenantInternalService)
	serverGrpcServerRegister := server2.NewGrpcServiceRegister(grpcServerRegister, serviceGrpcServerRegister, grpcServerRegister2)
	grpcServer := server2.NewGRPCServer(services, tokenizer, tenantStore, manager, webMultiTenancyOption, option, logger, trustedContextValidator, userTenantContrib, serverGrpcServerRegister)
	migrate := data2.NewMigrate(data4)
	roleSeed := biz2.NewRoleSeed(roleManager, permissionService)
	userSeed := biz2.NewUserSeed(userManager, roleManager)
	permissionSeeder := biz2.NewPermissionSeeder(permissionService, permissionService, roleManager)
	seeding := server3.NewSeeding(manager, migrate, roleSeed, userSeed, permissionSeeder)
	data5, cleanup5, err := data3.NewData(confData, dbProvider, logger)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	dataMigrate := data3.NewMigrate(data5)
	menuSeed := biz3.NewMenuSeed(dbProvider, menuRepo)
	serverSeeding := server4.NewSeeding(manager, dataMigrate, menuSeed)
	migrate2 := data.NewMigrate(dataData)
	seeding2 := server5.NewSeeding(manager, migrate2)
	seeder := server2.NewSeeder(tenantStore, seeding, serverSeeding, seeding2)
	userMigrationTaskHandler := biz2.NewUserMigrationTaskHandler(seeder, producer)
	jobServer := server2.NewJobServer(redisConnOpt, logger, userMigrationTaskHandler)
	tenantReadyEventHandler := biz.NewTenantReadyEventHandler(tenantUseCase)
	asynqClient, cleanup6 := job.NewAsynqClient(redisConnOpt)
	tenantSeedEventHandler := biz2.NewTenantSeedEventHandler(asynqClient)
	consumerFactoryServer := server2.NewEventServer(confData, connName, logger, manager, tenantReadyEventHandler, tenantSeedEventHandler)
	app := newApp(logger, userConf, httpServer, grpcServer, jobServer, consumerFactoryServer, seeder)
	return app, func() {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wireConfigValue             = dal.UowCfg
	_wireEventBusValue           = eventbus.Default
	_wireConnNameValue           = dal.ConnName("default")
	_wireDecodeRequestFuncValue  = server.ReqDecode
	_wireEncodeResponseFuncValue = server.ResEncoder
	_wireEncodeErrorFuncValue    = server.ErrEncoder
)