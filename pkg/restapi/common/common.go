/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/trustbloc/edge-core/pkg/log"
)

// Handler http handler for each controller API endpoint.
type Handler interface {
	Path() string
	Method() string
	Handle() http.HandlerFunc
}

// ErrorResponse to send error message in the response.
type ErrorResponse struct {
	Message string `json:"errMessage,omitempty"`
}

// WriteErrorResponsef write error resp.
func WriteErrorResponsef(rw http.ResponseWriter, logger log.Logger, status int, msg string, args ...interface{}) {
	rw.WriteHeader(status)

	err := json.NewEncoder(rw).Encode(ErrorResponse{
		Message: fmt.Sprintf(msg, args...),
	})

	if err != nil {
		logger.Errorf("Unable to send error message: %s", err)
	}
}
