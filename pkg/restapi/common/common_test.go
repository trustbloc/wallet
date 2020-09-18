/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edge-core/pkg/log/mocklogger"
)

func TestWriteErrorResponsef(t *testing.T) {
	t.Run("writes response", func(t *testing.T) {
		w := httptest.NewRecorder()
		expectedCode := http.StatusOK
		expectedMsg := uuid.New().String()

		WriteErrorResponsef(w, &mocklogger.MockLogger{}, expectedCode, expectedMsg)

		require.Equal(t, expectedCode, w.Code)
		result := &ErrorResponse{}

		err := json.NewDecoder(w.Body).Decode(result)
		require.NoError(t, err)
		require.Equal(t, expectedMsg, result.Message)
	})

	t.Run("logs error when writer fails", func(t *testing.T) {
		logger := &mocklogger.MockLogger{}
		WriteErrorResponsef(&mockHttpResponseWriter{writeErr: errors.New("test")}, logger, http.StatusOK, "test")
		require.Contains(t, logger.ErrorLogContents, "Unable to send error message")
	})
}

type mockHttpResponseWriter struct {
	writeErr error
}

func (m *mockHttpResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (m *mockHttpResponseWriter) Write(_ []byte) (int, error) {
	return 0, m.writeErr
}

func (m *mockHttpResponseWriter) WriteHeader(_ int) {

}
