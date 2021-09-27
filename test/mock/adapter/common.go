/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"net/http"
	"text/template"
)

func loadTemplate(w http.ResponseWriter, fileName string, data map[string]interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := t.Execute(w, data); err != nil {
		return err
	}

	return nil
}

func handleError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(ErrorResponse{
		Message: msg,
	})

	if err != nil {
		logger.Errorf("Unable to send error message, %s", err)
	}
}

// ErrorResponse to send error message in the response.
type ErrorResponse struct {
	Message string `json:"errMessage,omitempty"`
}
