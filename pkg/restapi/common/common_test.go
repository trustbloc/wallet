/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/edge-core/pkg/log/mocklogger"

	"github.com/trustbloc/edge-agent/pkg/restapi/common"
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

func TestSendHTTPRequest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body, header, err := common.SendHTTPRequest(
			httptest.NewRequest(http.MethodGet, "http://example.com", nil),
			MockHTTPClient{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Cat": {"dog"}},
					Body:       ioutil.NopCloser(strings.NewReader("hello there")),
				},
				returnErr: nil,
			}, http.StatusOK, nil)

		require.NoError(t, err)
		require.NotNil(t, body)
		require.Contains(t, string(body), "hello there")
		require.NotNil(t, header)
		require.Contains(t, header.Values("cat"), "dog")
	})

	t.Run("wrong status code", func(t *testing.T) {
		body, header, err := common.SendHTTPRequest(
			httptest.NewRequest(http.MethodGet, "http://example.com", nil),
			MockHTTPClient{
				response: &http.Response{
					StatusCode: http.StatusNotFound,
					Header:     http.Header{"Cat": {"dog"}},
					Body:       ioutil.NopCloser(strings.NewReader("hello there")),
				},
				returnErr: nil,
			}, http.StatusOK, nil)

		require.Error(t, err)
		require.Contains(t, err.Error(), "expected=200 actual=404")
		require.Nil(t, body)
		require.Nil(t, header)
	})

	t.Run("response body errors on closing", func(t *testing.T) {
		var mockLog mocklogger.MockLogger

		log.Initialize(&mocklogger.Provider{MockLogger: &mockLog})

		_, _, err := common.SendHTTPRequest(
			httptest.NewRequest(http.MethodGet, "http://example.com", nil),
			MockHTTPClient{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Cat": {"dog"}},
					Body:       &brokenReadCloser{readErr: io.EOF, closeErr: fmt.Errorf("oh no")},
				},
				returnErr: nil,
			}, http.StatusOK, log.New("mock"))

		require.NoError(t, err)

		require.Contains(t, mockLog.AllLogContents, "failed to close response body")
	})

	t.Run("response body errors on read", func(t *testing.T) {
		body, header, err := common.SendHTTPRequest(
			httptest.NewRequest(http.MethodGet, "http://example.com", nil),
			MockHTTPClient{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Cat": {"dog"}},
					Body:       &brokenReadCloser{readErr: fmt.Errorf("oh no"), closeErr: nil},
				},
				returnErr: nil,
			}, http.StatusOK, nil)

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to read resp body")
		require.Nil(t, body)
		require.Nil(t, header)
	})

	t.Run("client errors on invocation", func(t *testing.T) {
		body, header, err := common.SendHTTPRequest(
			httptest.NewRequest(http.MethodGet, "http://example.com", nil),
			MockHTTPClient{
				returnErr: fmt.Errorf("oh no"),
			}, http.StatusOK, nil)

		require.Error(t, err)
		require.Contains(t, err.Error(), "http request")
		require.Nil(t, body)
		require.Nil(t, header)
	})
}

type MockHTTPClient struct {
	response  *http.Response
	returnErr error
}

func (m MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.returnErr
}

type brokenReadCloser struct {
	readErr  error
	closeErr error
}

func (bc *brokenReadCloser) Read(p []byte) (int, error) {
	return 0, bc.readErr
}

func (bc *brokenReadCloser) Close() error {
	return bc.closeErr
}
