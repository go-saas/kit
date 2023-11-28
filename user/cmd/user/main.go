package main

import (
	"context"
	"flag"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	dtmserver "github.com/go-saas/kit/dtm/server"
	"github.com/go-saas/kit/event"
	eventserver "github.com/go-saas/kit/event/server"
	oidcdata "github.com/go-saas/kit/oidc/data"
	kapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/authz/casbin"
	kitconf "github.com/go-saas/kit/pkg/conf"
	kdal "github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	kitflag "github.com/go-saas/kit/pkg/flag"
	"github.com/go-saas/kit/pkg/job"
	"github.com/go-saas/kit/pkg/logging"
	kitserver "github.com/go-saas/kit/pkg/server"
	"github.com/go-saas/kit/pkg/tracers"
	sapi "github.com/go-saas/kit/saas/api"
	"github.com/go-saas/kit/user/private/biz"
	"github.com/go-saas/kit/user/private/data"
	"github.com/go-saas/kit/user/private/server"
	"github.com/go-saas/kit/user/private/service"
	"github.com/go-saas/saas/seed"
	"github.com/goava/di"
	"github.com/goxiaoy/vfs"
	"github.com/spf13/afero"
	"os"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-saas/kit/user/private/conf"

	//this helps to load permission definitions in following services
	_ "github.com/go-saas/kit/realtime/api"
	_ "github.com/go-saas/kit/saas/api"
	_ "github.com/go-saas/kit/sys/api"

	_ "github.com/go-saas/kit/order/api"
	_ "github.com/go-saas/kit/payment/api"
	_ "github.com/go-saas/kit/product/api"
)

// go build -buildvcs=false -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "user"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf kitflag.ArrayFlags
	ifSeed   bool

	id, _ = os.Hostname()
)

func init() {
	flag.Var(&flagconf, "conf", "config path, eg: -conf config.yaml")
	flag.BoolVar(&ifSeed, "seed", true, "run seeder or not")
}

func newApp(
	c *conf.UserConf,
	logger log.Logger,
	seeder seed.Seeder,
	producer event.Producer,
	_ dtmserver.Init,
	r registry.Registrar,
	srvs []transport.Server,
) *kratos.App {
	ctx := event.NewProducerContext(context.Background(), producer)
	if ifSeed {
		if err := seeder.Seed(ctx,
			seed.AddHost(),
			seed.SetExtra(biz.AdminUsernameKey, c.Admin.GetUsername()),
			seed.SetExtra(biz.AdminPasswordKey, c.Admin.GetPassword()),
		); err != nil {
			panic(err)
		}
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
		source = append(source, kitconf.NewVfs(v, "/"))
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
		kitdi.Value(kitserver.Name(Name)),
		kitdi.Value(bc.Services),
		kitdi.Value(bc.Security),
		kitdi.Value(bc.App),
		kitdi.Value(bc.Dev),
		kitdi.Value(bc.User),
		kitdi.Value(bc.User.Idp),
		kitdi.Value(bc.Stripe),
		kitdi.Value(bc.Data),
		kitdi.Value(logger),
		kitdi.Value([]grpc.ClientOption{}),
		authz.ProviderSet, kitserver.DefaultProviderSet, jwt.ProviderSet, kapi.DefaultProviderSet, kdal.DefaultProviderSet, oidcdata.ProviderSet,
		job.DefaultProviderSet, dtmserver.DtmProviderSet, eventserver.EventProviderSet,
		sapi.GrpcProviderSet,
		casbin.PermissionProviderSet, server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet,
		kitdi.NewSet(newApp),
	)
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
