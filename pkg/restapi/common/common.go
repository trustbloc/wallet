/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"encoding/json"
	"fmt"
	"net/http"
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
