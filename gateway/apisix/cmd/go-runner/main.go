/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-saas/kit/gateway/apisix/internal/conf"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/authz/authz"
	conf2 "github.com/go-saas/kit/pkg/conf"
	kitdi "github.com/go-saas/kit/pkg/di"
	"github.com/go-saas/kit/pkg/logging"
	"github.com/go-saas/kit/pkg/tracers"
	sapi "github.com/go-saas/kit/saas/api"
	uapi "github.com/go-saas/kit/user/api"
	"github.com/goava/di"
	"github.com/goxiaoy/vfs"
	"github.com/spf13/afero"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"

	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/runner"

	klog "github.com/go-kratos/kratos/v2/log"

	_ "github.com/go-saas/kit/gateway/apisix/cmd/go-runner/plugins"

	_ "github.com/go-saas/kit/pkg/registry/etcd"
)

var (
	InfoOut io.Writer = os.Stdout
)

func newVersionCommand() *cobra.Command {
	var long bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "version",
		Run: func(cmd *cobra.Command, _ []string) {
			if long {
				fmt.Fprint(InfoOut, longVersion())
			} else {
				fmt.Fprintf(InfoOut, "version %s\n", shortVersion())
			}
		},
	}

	cmd.PersistentFlags().BoolVar(&long, "long", false, "show long mode version information")
	return cmd
}

type RunMode enumflag.Flag

const (
	Dev  RunMode = iota // Development
	Prod                // Product
	Prof                // Profile

	ProfileFilePath = "./logs/profile."
	LogFilePath     = "./logs/runner.log"
)

var RunModeIds = map[RunMode][]string{
	Prod: {"prod"},
	Dev:  {"dev"},
	Prof: {"prof"},
}

func openFileToWrite(name string) (*os.File, error) {
	dir := filepath.Dir(name)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func newRunCommand() *cobra.Command {
	var mode RunMode
	var clientName string
	var flagconf []string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "run",
		Run: func(cmd *cobra.Command, _ []string) {
			cfg := runner.RunnerConfig{}
			if mode == Prod {
				cfg.LogLevel = zapcore.WarnLevel
				f, err := openFileToWrite(LogFilePath)
				if err != nil {
					log.Fatalf("failed to open log: %s", err)
				}
				cfg.LogOutput = f
			} else if mode == Prof {
				cfg.LogLevel = zapcore.WarnLevel

				cpuProfileFile := ProfileFilePath + "cpu"
				f, err := os.Create(cpuProfileFile)
				if err != nil {
					log.Fatalf("could not create CPU profile: %s", err)
				}
				defer f.Close()
				if err := pprof.StartCPUProfile(f); err != nil {
					log.Fatalf("could not start CPU profile: %s", err)
				}
				defer pprof.StopCPUProfile()

				defer func() {
					memProfileFile := ProfileFilePath + "mem"
					f, err := os.Create(memProfileFile)
					if err != nil {
						log.Fatalf("could not create memory profile: %s", err)
					}
					defer f.Close()

					runtime.GC()
					if err := pprof.WriteHeapProfile(f); err != nil {
						log.Fatalf("could not write memory profile: %s", err)
					}
				}()
			}

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
			c := config.New(config.WithSource(source...))
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
			logger := klog.With(l,
				"ts", klog.DefaultTimestamp,
				"caller", klog.DefaultCaller,
			)

			shutdown, err := tracers.SetTracerProvider(context.Background(), bc.Tracing, "apisix")
			if err != nil {
				logger.Log(klog.LevelError, err)
			}
			defer shutdown()

			di.SetTracer(&di.StdTracer{})
			builder, err := di.New(
				kitdi.Value(bc.Services),
				kitdi.Value(bc.Security),
				kitdi.Value(bc.App),
				kitdi.Value(api.ClientName(clientName)),
				kitdi.Value([]grpc.ClientOption{}),
				kitdi.Value(logger),
				kitdi.NewSet(
					ProviderSet, authz.ProviderSet, sapi.GrpcProviderSet, uapi.GrpcProviderSet, jwt.ProviderSet, newApp),
			)

			if err != nil {
				panic(err)
			}

			defer builder.Cleanup()
			err = builder.Invoke(func(app *App) error {
				// start and wait for stop signal
				return app.load()
			})
			if err != nil {
				panic(err)
			}
			runner.Run(cfg)
		},
	}

	cmd.PersistentFlags().VarP(
		enumflag.New(&mode, "mode", RunModeIds, enumflag.EnumCaseInsensitive),
		"mode", "m",
		"the runner's run mode; can be 'prod' or 'dev', default to 'dev'")
	cmd.PersistentFlags().StringVarP(&clientName, "client", "n", "apisix", "client name")
	cmd.PersistentFlags().StringSliceVarP(&flagconf, "conf", "c", []string{"./configs"}, "config path, eg: -conf config.yaml")

	return cmd
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apisix-go-plugin-runner [command]",
		Long:    "The Plugin runner to run Go plugins",
		Version: shortVersion(),
	}

	cmd.AddCommand(newRunCommand())
	cmd.AddCommand(newVersionCommand())
	return cmd
}

func main() {
	root := NewCommand()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
