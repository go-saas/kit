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
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/goxiaoy/go-saas-kit/gateway/apisix/internal/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"

	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/runner"

	_ "github.com/goxiaoy/go-saas-kit/gateway/apisix/cmd/go-runner/plugins"
	_ "github.com/goxiaoy/go-saas/gateway/apisix"
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
	var flagconf string
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
			c := config.New(
				config.WithSource(
					file.NewSource(strings.TrimSpace(flagconf)),
				),
			)
			if err := c.Load(); err != nil {
				panic(err)
			}

			var bc conf.Bootstrap
			if err := c.Scan(&bc); err != nil {
				panic(err)
			}

			//init all
			_, clean, err := initApp(bc.Services, bc.Security, api.ClientName(clientName))
			if err != nil {
				panic(err)
			}
			defer clean()
			runner.Run(cfg)
		},
	}

	cmd.PersistentFlags().VarP(
		enumflag.New(&mode, "mode", RunModeIds, enumflag.EnumCaseInsensitive),
		"mode", "m",
		"the runner's run mode; can be 'prod' or 'dev', default to 'dev'")
	cmd.PersistentFlags().StringVarP(&clientName, "client", "n", "apisix", "client name")
	cmd.PersistentFlags().StringVarP(&flagconf, "conf", "c", "../../configs", "config path, eg: -conf config.yaml")
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
