/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edge-agent/pkg/restapi/common"
	"github.com/trustbloc/edge-core/pkg/log/mocklogger"
)

func TestWriteErrorResponsef(t *testing.T) {
	t.Run("writes response", func(t *testing.T) {
		w := httptest.NewRecorder()
		expectedCode := http.StatusOK
		expectedMsg := uuid.New().String()

		common.WriteErrorResponsef(w, &mocklogger.MockLogger{}, expectedCode, expectedMsg)

		require.Equal(t, expectedCode, w.Code)
		result := &common.ErrorResponse{}

		err := json.NewDecoder(w.Body).Decode(result)
		require.NoError(t, err)
		require.Equal(t, expectedMsg, result.Message)
	})

	t.Run("logs error when writer fails", func(t *testing.T) {
		logger := &mocklogger.MockLogger{}
		common.WriteErrorResponsef(&mockHTTPResponseWriter{writeErr: errors.New("test")}, logger, http.StatusOK, "test")
		require.Contains(t, logger.ErrorLogContents, "Unable to send error message")
	})
}

func TestWriteResponse(t *testing.T) {
	t.Run("writes response", func(t *testing.T) {
		expected := map[string]interface{}{
			"data": uuid.New().String(),
		}
		result := httptest.NewRecorder()
		common.WriteResponse(result, &mocklogger.MockLogger{}, expected)
		require.Equal(t, http.StatusOK, result.Code)
		body := result.Body.String()
		require.NotEmpty(t, body)
		data := make(map[string]interface{})
		err := json.Unmarshal([]byte(body), &data)
		require.NoError(t, err)
		require.Equal(t, expected, data)
	})

	t.Run("logs error if cannot write response", func(t *testing.T) {
		expected := errors.New("test")
		result := &mocklogger.MockLogger{}
		common.WriteResponse(&mockHTTPResponseWriter{writeErr: expected}, result, expected.Error())
		require.Contains(t, result.ErrorLogContents, expected.Error())
	})
}

type mockHTTPResponseWriter struct {
	writeErr   error
	writtenVal []byte
}

func (m *mockHTTPResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (m *mockHTTPResponseWriter) Write(b []byte) (int, error) {
	m.writtenVal = b

	return 0, m.writeErr
}

func (m *mockHTTPResponseWriter) WriteHeader(_ int) {}
