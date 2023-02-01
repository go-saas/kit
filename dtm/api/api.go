package api

import "github.com/go-saas/kit/pkg/conf"

//
//func NewInit(client *conf.Client, opt *api.Option, tokenMgr api.TokenManager, logger klog.Logger) *Init {
//	//token interceptor
//	m := api.ClientPropagation(client, opt, tokenMgr, logger)
//	dtmgrpc.AddUnaryInterceptor(api.UnaryClientInterceptor([]middleware.Middleware{m}, 0, nil))
//	dtmgrpc.UseDriver(driver.DriverName)
//	return &Init{}
//}
//
//type Init struct {
//}

const ServiceName = "dtmservice"

var ClientConf = &conf.Client{
	ClientId: ServiceName,
}
