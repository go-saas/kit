package server

import (
	"github.com/dtm-labs/client/dtmgrpc"
	driver "github.com/dtm-labs/dtmdriver-kratos"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"
	v1 "github.com/go-saas/kit/dtm/api/dtm/v1"
	"github.com/go-saas/kit/dtm/data"
	"github.com/go-saas/kit/dtm/service"
	sapi "github.com/go-saas/kit/pkg/api"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/goava/di"
	"google.golang.org/grpc"
	"sync"
)

var DtmProviderSet = kitdi.NewSet(
	NewInit,
	kitdi.NewProvider(service.NewMsgService, di.As(new(v1.MsgServiceServer))),
	data.NewBarrierMigrator, data.NewStorageMigrator, data.NewMigrator,
	service.NewHttpServerRegister, service.NewGrpcServerRegister,
)

var (
	once sync.Once
)

type Init any

func NewInit(dis registry.Discovery, opt *sapi.Option) (Init, error) {
	var opts = []discovery.Option{
		discovery.WithInsecure(opt.Insecure),
	}
	//TODO subset
	//if opt.Subset != nil {
	//	opts = append(opts, discovery.WithSubset(*opt.Subset))
	//}
	dtmgrpc.AddDailOption(grpc.WithResolvers(discovery.NewBuilder(dis, opts...)))
	dtmgrpc.UseDriver(driver.DriverName)
	return "", nil
}
