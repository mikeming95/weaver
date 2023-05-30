// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package weavertest

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/ServiceWeaver/weaver/internal/private"
	"github.com/ServiceWeaver/weaver/runtime"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"github.com/ServiceWeaver/weaver/runtime/protos"
	"github.com/google/uuid"
)

const matchNothingRE = "a^" // Regular expression that never matches

// initMultiProcess initializes a brand new multi-process execution environment
// that places every component in its own collocation group. It returns a
// function that can be used to stop the execution.
//
// logWriter is used to handle log entries generated by the execution.
//
// config contains configuration identical to what might be found in a file passed
// when deploying an application. It can contain application level as well as
// component level configs. config is allowed to be empty.
//
// Future extension: allow options so the user can control collocation/replication/etc.
func initMultiProcess(ctx context.Context, t testing.TB, isBench bool, runner Runner, logWriter func(*protos.LogEntry)) (context.Context, func() error, error) {
	t.Helper()
	bootstrap, err := runtime.GetBootstrap(ctx)
	if err != nil {
		return nil, nil, err
	}
	if bootstrap.HasPipes() {
		// This is a child process, so just start the application and wait.
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "panic in Service Weaver sub-process: %v\n", r)
			} else {
				fmt.Fprintf(os.Stderr, "Service Weaver sub-process exiting\n")
			}
			os.Exit(1)
		}()

		err := runWeaver(ctx, t, runner, func(context.Context, private.App) error {
			<-ctx.Done() // Wait for parent process
			return ctx.Err()
		})
		if err != nil {
			panic(err)
		}
		return nil, nil, nil
	}

	// Construct AppConfig and EnvelopeInfo.
	appConfig := &protos.AppConfig{}
	if runner.Config != "" {
		var err error
		appConfig, err = runtime.ParseConfig("[testconfig]", runner.Config, codegen.ComponentConfigValidator)
		if err != nil {
			return nil, nil, err
		}
	}
	exe, err := os.Executable()
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching binary path: %v", err)
	}
	name := t.Name()
	appConfig.Name = strings.ReplaceAll(name, "/", "_")
	appConfig.Binary = exe
	nameRE := "^" + regexp.QuoteMeta(name) + "$"
	if isBench {
		appConfig.Args = []string{"-test.run", matchNothingRE, "-test.bench", nameRE}
	} else {
		appConfig.Args = []string{"-test.run", nameRE}
	}

	wlet := &protos.EnvelopeInfo{
		App:           appConfig.Name,
		DeploymentId:  uuid.New().String(),
		Sections:      appConfig.Sections,
		SingleProcess: false,
		SingleMachine: true,
	}

	// Launch the deployer.
	d := newDeployer(ctx, wlet, appConfig, runner, logWriter)
	if err := d.start(); err != nil {
		return nil, nil, err
	}
	return d.ctx, d.cleanup, nil
}
