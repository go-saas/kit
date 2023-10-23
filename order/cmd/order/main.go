package main

import (
	"context"
	"flag"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/order/api"
	"github.com/go-saas/kit/order/private/biz"
	"github.com/go-saas/kit/order/private/conf"
	"github.com/go-saas/kit/order/private/data"
	"github.com/go-saas/kit/order/private/server"
	"github.com/go-saas/kit/order/private/service"
	kapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/authz/authz"
	conf2 "github.com/go-saas/kit/pkg/conf"
	kdal "github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	kitflag "github.com/go-saas/kit/pkg/flag"
	"github.com/go-saas/kit/pkg/logging"
	kitserver "github.com/go-saas/kit/pkg/server"
	"github.com/go-saas/kit/pkg/tracers"
	sapi "github.com/go-saas/kit/saas/api"
	uapi "github.com/go-saas/kit/user/api"
	"github.com/go-saas/saas/seed"
	"github.com/goava/di"
	"github.com/goxiaoy/vfs"
	"github.com/spf13/afero"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-saas/kit/event/kafka"
	_ "github.com/go-saas/kit/event/pulsar"

	_ "github.com/go-saas/kit/pkg/registry/etcd"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = api.ServiceName
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf kitflag.ArrayFlags

	id, _ = os.Hostname()
)

func init() {
	flag.Var(&flagconf, "conf", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, srvs []transport.Server, r registry.Registrar, seeder seed.Seeder, producer event.Producer) *kratos.App {
	ctx := event.NewProducerContext(context.Background(), producer)
	if err := seeder.Seed(context.Background(), seed.AddHost()); err != nil {
		panic(err)
	}
	return kratos.New(
		kratos.Context(ctx),
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Registrar(r),
		kratos.Server(
			srvs...,
		),
	)
}

func main() {
	flag.Parse()

	source := []config.Source{
		env.NewSource("KRATOS_"),
	}
	if flagconf == nil {
		flagconf = append(flagconf, "./configs")
	}
	for _, s := range flagconf {
		v := vfs.New()
		v.Mount("/", afero.NewRegexpFs(afero.NewBasePathFs(afero.NewOsFs(), strings.TrimSpace(s)), regexp.MustCompile(`\.(json|proto|xml|yaml)$`)))
		source = append(source, conf2.NewVfs(v, "/"))
	}

	c := config.New(
		config.WithSource(
			source...,
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	l, lc, err := logging.NewLogger(bc.Logging)
	if err != nil {
		panic(err)
	}
	defer lc()

	logger := log.With(l,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)

	shutdown, err := tracers.SetTracerProvider(context.Background(), bc.Tracing, Name)
	if err != nil {
		log.Error(err)
	}
	defer shutdown()

	di.SetTracer(&di.StdTracer{})
	builder, err := di.New(
		kitdi.Value(bc.Services),
		kitdi.Value(bc.Security),
		kitdi.Value(bc.App),
		kitdi.Value(bc.Data),
		kitdi.Value(logger),
		kitdi.Value([]grpc.ClientOption{}),
		authz.ProviderSet, kitserver.DefaultProviderSet, jwt.ProviderSet, kapi.DefaultProviderSet, kdal.DefaultProviderSet,
		uapi.GrpcProviderSet,
		sapi.GrpcProviderSet,
		server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, kitdi.NewSet(newApp))
	if err != nil {
		panic(err)
	}
	defer builder.Cleanup()
	err = builder.Invoke(func(app *kratos.App) error {
		// start and wait for stop signal
		return app.Run()
	})

	if err != nil {
		panic(err)
	}

}
