/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/trustbloc/edge-core/pkg/log"
)

// Handler http handler for each controller API endpoint.
type Handler interface {
	Path() string
	Method() string
	Handle() http.HandlerFunc
}

type logger interface {
	Errorf(string, ...interface{})
}

// ErrorResponse to send error message in the response.
type ErrorResponse struct {
	Message string `json:"errMessage,omitempty"`
}

// WriteErrorResponsef write error resp.
func WriteErrorResponsef(rw http.ResponseWriter, logger logger, status int, msg string, args ...interface{}) {
	logger.Errorf(msg, args...)

	rw.WriteHeader(status)

	err := json.NewEncoder(rw).Encode(ErrorResponse{
		Message: fmt.Sprintf(msg, args...),
	})
	if err != nil {
		logger.Errorf("Unable to send error message: %s", err)
	}
}

// WriteResponse writes interface value to response.
func WriteResponse(rw http.ResponseWriter, l logger, v interface{}) {
	rw.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(rw).Encode(v)
	if err != nil {
		l.Errorf("Unable to send error response, %s", err.Error())
	}
}

// HTTPClient http client interface.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// SendHTTPRequest sends an http request to the client, expecting a given status, and returning the result.
func SendHTTPRequest(req *http.Request, httpClient HTTPClient, status int,
	logger *log.Log) ([]byte, http.Header, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("http request : %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			if logger != nil {
				logger.Errorf("failed to close response body")
			}
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("http request: failed to read resp body %d : %w", resp.StatusCode, err)
	}

	if resp.StatusCode != status {
		return nil, nil, fmt.Errorf("http request: expected=%d actual=%d body=%s", status, resp.StatusCode, string(body))
	}

	return body, resp.Header, nil
}
