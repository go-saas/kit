package main

import (
	"context"
	"flag"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	dtmserver "github.com/go-saas/kit/dtm/server"
	"github.com/go-saas/kit/event"
	eventserver "github.com/go-saas/kit/event/server"
	"github.com/go-saas/kit/examples/monolithic/private/conf"
	"github.com/go-saas/kit/examples/monolithic/private/server"
	kapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/apisix"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/authz/casbin"
	conf2 "github.com/go-saas/kit/pkg/conf"
	kdal "github.com/go-saas/kit/pkg/dal"
	kitdi "github.com/go-saas/kit/pkg/di"
	kitflag "github.com/go-saas/kit/pkg/flag"
	"github.com/go-saas/kit/pkg/job"
	"github.com/go-saas/kit/pkg/logging"
	sserver "github.com/go-saas/kit/pkg/server"
	"github.com/go-saas/kit/pkg/tracers"
	rbiz "github.com/go-saas/kit/realtime/private/biz"
	rdata "github.com/go-saas/kit/realtime/private/data"
	rservice "github.com/go-saas/kit/realtime/private/service"
	sbiz "github.com/go-saas/kit/saas/private/biz"
	sdata "github.com/go-saas/kit/saas/private/data"
	sservice "github.com/go-saas/kit/saas/private/service"
	sysbiz "github.com/go-saas/kit/sys/private/biz"
	sysdata "github.com/go-saas/kit/sys/private/data"
	sysservice "github.com/go-saas/kit/sys/private/service"
	ubiz "github.com/go-saas/kit/user/private/biz"
	uconf "github.com/go-saas/kit/user/private/conf"
	udata "github.com/go-saas/kit/user/private/data"
	uservice "github.com/go-saas/kit/user/private/service"
	"github.com/go-saas/saas/seed"
	"github.com/goava/di"
	"github.com/goxiaoy/vfs"
	"github.com/spf13/afero"
	"os"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
)

// go build -buildvcs=false -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "monolithic"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf     kitflag.ArrayFlags
	ifSyncApisix bool
	ifSeed       bool
	id, _        = os.Hostname()
)

func init() {
	flag.Var(&flagconf, "conf", "config path, eg: -conf config.yaml")
	flag.BoolVar(&ifSyncApisix, "apisix.sync", true, "sync with apisix upstreams")
	flag.BoolVar(&ifSeed, "seed", true, "run seeder or not")
}

func newApp(
	logger log.Logger,
	c *uconf.UserConf,
	srvs []transport.Server,
	seeder seed.Seeder,
	producer event.Producer,
	r registry.Registrar,
	syncAdmin *apisix.WatchSyncAdmin,
) *kratos.App {
	ctx := event.NewProducerContext(context.Background(), producer)
	if ifSeed {
		extra := map[string]interface{}{}
		extra[ubiz.AdminUsernameKey] = c.Admin.GetUsername()
		extra[ubiz.AdminPasswordKey] = c.Admin.GetPassword()
		if err := seeder.Seed(ctx, seed.AddHost(), seed.WithExtra(extra)); err != nil {
			panic(err)
		}
	}

	if ifSyncApisix {
		srvs = append(srvs, syncAdmin)
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
		kitdi.Value(bc.Data),
		kitdi.Value(bc.Sys),
		kitdi.Value(bc.Saas),
		kitdi.Value(bc.User),
		kitdi.Value(bc.App),
		kitdi.Value(logger),
		kitdi.NewSet(authz.ProviderSet, jwt.ProviderSet, sserver.DefaultProviderSet, kapi.DefaultProviderSet, kdal.DefaultProviderSet,
			job.DefaultProviderSet, dtmserver.DtmProviderSet, eventserver.EventProviderSet,
			//saas
			sdata.ProviderSet, sbiz.ProviderSet, sservice.ProviderSet,
			//sys
			sysdata.ProviderSet, sysbiz.ProviderSet, sysservice.ProviderSet,
			//user
			udata.ProviderSet, ubiz.ProviderSet, uservice.ProviderSet,
			//realtime
			rdata.ProviderSet, rbiz.ProviderSet, rservice.ProviderSet,

			casbin.PermissionProviderSet, server.ProviderSet,
			newApp),
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
