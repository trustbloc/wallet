/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package bdd_test

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/trustbloc/edge-agent/test/bdd/dockerutil"
	"github.com/trustbloc/edge-agent/test/bdd/pkg/bddcontext"
	"github.com/trustbloc/edge-agent/test/bdd/pkg/healthcheck"
	"github.com/trustbloc/edge-agent/test/bdd/pkg/login"
)

func TestMain(m *testing.M) {
	// default is to run all tests with tag @all
	tags := "all"

	if os.Getenv("TAGS") != "" {
		tags = os.Getenv("TAGS")
	}

	flag.Parse()

	format := "progress"
	if getCmdArg("test.v") == "true" {
		format = "pretty"
	}

	runArg := getCmdArg("test.run")
	if runArg != "" {
		tags = runArg
	}

	status := runBDDTests(tags, format)
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func runBDDTests(tags, format string) int {
	return godog.RunWithOptions("godogs", func(s *godog.Suite) {
		var composition []*dockerutil.Composition
		composeFiles := []string{"./fixtures/agent-wasm"}
		s.BeforeSuite(beforeSuite(composition, composeFiles))
		s.AfterSuite(afterSuite(composition))
		featureContext(s)
	}, godog.Options{
		Tags:          tags,
		Format:        format,
		Paths:         []string{"features"},
		Randomize:     time.Now().UTC().UnixNano(), // randomize scenario execution order
		Strict:        true,
		StopOnFailure: true,
	})
}

func beforeSuite(composition []*dockerutil.Composition, composeFiles []string) func() {
	return func() {
		if os.Getenv("DISABLE_COMPOSITION") != "true" {
			composeProjectName := strings.ReplaceAll(uuid.New().String(), "-", "")

			for _, v := range composeFiles {
				newComposition, err := dockerutil.NewComposition(composeProjectName, "docker-compose.yml", v)
				if err != nil {
					panic(fmt.Sprintf("Error composing system in BDD context: %s", err))
				}

				composition = append(composition, newComposition)
			}

			fmt.Println("docker-compose up ... waiting for containers to start ...")

			testSleep := 30

			if os.Getenv("TEST_SLEEP") != "" {
				var e error

				testSleep, e = strconv.Atoi(os.Getenv("TEST_SLEEP"))
				if e != nil {
					panic(fmt.Sprintf("Invalid value found in 'TEST_SLEEP': %s", e))
				}
			}

			fmt.Printf("*** testSleep=%d", testSleep)
			println()
			time.Sleep(time.Second * time.Duration(testSleep))
		}
	}
}

func afterSuite(composition []*dockerutil.Composition) func() {
	return func() {
		for _, c := range composition {
			if c != nil {
				if err := c.GenerateLogs(c.Dir, "docker-compose.log"); err != nil {
					panic(err)
				}

				if _, err := c.Decompose(c.Dir); err != nil {
					panic(err)
				}
			}
		}
	}
}

func getCmdArg(argName string) string {
	cmdTags := flag.CommandLine.Lookup(argName)
	if cmdTags != nil && cmdTags.Value != nil && cmdTags.Value.String() != "" {
		return cmdTags.Value.String()
	}

	return ""
}

func featureContext(s *godog.Suite) {
	ctx, err := bddcontext.NewBDDContext("fixtures/keys/tls/ec-cacert.pem")
	if err != nil {
		panic(fmt.Sprintf("Error returned from NewBDDContext: %s", err))
	}

	login.NewSteps(ctx).Register(s)
	healthcheck.NewSteps(ctx).Register(s)
}
