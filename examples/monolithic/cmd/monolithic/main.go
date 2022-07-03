package main

import (
	"context"
	"flag"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/examples/monolithic/private/conf"
	"github.com/go-saas/kit/pkg/job"
	"github.com/go-saas/kit/pkg/logging"
	"github.com/go-saas/kit/pkg/tracers"
	sysbiz "github.com/go-saas/kit/sys/private/biz"
	ubiz "github.com/go-saas/kit/user/private/biz"
	uconf "github.com/go-saas/kit/user/private/conf"
	"github.com/go-saas/saas/seed"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "github.com/go-saas/kit/event/kafka"
	_ "github.com/go-saas/kit/event/pulsar"
	_ "github.com/go-saas/kit/pkg/registry/etcd"
)

// go build -buildvcs=false -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "monolithic"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf arrayFlags
	ifSeed   bool
	id, _    = os.Hostname()

	seedPath string
)

func init() {
	flag.Var(&flagconf, "conf", "config path, eg: -conf config.yaml")
	flag.BoolVar(&ifSeed, "seed", true, "run seeder or not")
	flag.StringVar(&seedPath, sysbiz.SeedPathKey, "", "menu seed file path")
}

func newApp(logger log.Logger, c *uconf.UserConf, hs *http.Server, gs *grpc.Server, js *job.Server, es *event.ConsumerFactoryServer, seeder seed.Seeder, r registry.Registrar) *kratos.App {
	if ifSeed {
		extra := map[string]interface{}{}
		if len(seedPath) > 0 {
			extra[sysbiz.SeedPathKey] = seedPath
		}
		extra[ubiz.AdminUsernameKey] = c.Admin.GetUsername()
		extra[ubiz.AdminPasswordKey] = c.Admin.GetPassword()
		if err := seeder.Seed(context.Background(), seed.AddHost(), seed.WithExtra(extra)); err != nil {
			panic(err)
		}
	}
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Registrar(r),
		kratos.Server(
			hs,
			gs,
			js,
			es,
		),
	)
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	flag.Parse()

	source := []config.Source{
		env.NewSource("KRATOS_"),
	}
	if flagconf != nil {
		for _, s := range flagconf {
			source = append(source, file.NewSource(strings.TrimSpace(s)))
		}
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
	app, cleanup, err := initApp(bc.Services, bc.Security, bc.Data, bc.Saas, bc.User, logger, bc.App)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
