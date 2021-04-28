/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package healthcheck

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cucumber/godog"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/trustbloc/edge-core/pkg/log"

	"github.com/trustbloc/edge-agent/test/bdd/pkg/bddcontext"
)

var logger = log.New("bdd-test-healthcheck")

// Steps is steps for VC BDD tests.
type Steps struct {
	bddContext *bddcontext.BDDContext
	queryValue string
}

// NewSteps returns new agent from client SDK.
func NewSteps(ctx *bddcontext.BDDContext) *Steps {
	return &Steps{bddContext: ctx}
}

// Register registers agent steps.
func (s *Steps) Register(gs *godog.Suite) {
	gs.Step(`^HTTP GET is sent to "([^"]*)"$`, s.httpGet)
	gs.Step(`^The "([^"]*)" field in the response has the value "([^"]*)"$`, s.validateJSONResField)
}

// httpGet sends a GET request to the given URL.
func (s *Steps) httpGet(url string) error {
	s.queryValue = ""

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: s.bddContext.TLSConfig}}
	defer client.CloseIdleConnections()

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}

	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			logger.Warnf("Error closing HTTP response from [%s]: %s", url, errClose)
		}
	}()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body failed: %w", err)
	}

	s.queryValue = string(payload)

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("received status code %d", resp.StatusCode)
	}

	return nil
}

func (s *Steps) validateJSONResField(path, expected string) error {
	r := gjson.Get(s.queryValue, path)

	logger.Infof("Path [%s] of JSON %s resolves to %s", path, s.queryValue, r.Str)

	if r.Str == expected {
		return nil
	}

	return fmt.Errorf("JSON path resolves to [%s] which is not the expected value [%s]", r.Str, expected)
}
