/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chapibridge // nolint:testpackage // changing to different package requires exposing internal features.

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mediatorsvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/mediator"
	outofbandsvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/outofband"
	mockmediator "github.com/hyperledger/aries-framework-go/pkg/mock/didcomm/protocol/mediator"
	"github.com/stretchr/testify/require"

	mockprotocol "github.com/trustbloc/edge-agent/pkg/restapi/internal/mocks/protocol"
)

const (
	sampleAppURL = "http://demo.wallet.app/home"
	sampleErr    = "sample-error"
)

func TestNew(t *testing.T) {
	t.Run("create new instance - success", func(t *testing.T) {
		op, err := New(newMockProvider(), "test", "demo")

		require.NoError(t, err)
		require.NotEmpty(t, op)
		require.Len(t, op.GetRESTHandlers(), 3)
	})

	t.Run("create new instance - failure", func(t *testing.T) {
		op, err := New(mockprotocol.NewMockProvider(), "test", "demo")

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create out-of-band client")
		require.Empty(t, op)
	})
}

func TestOperation_CreateInvitation(t *testing.T) {
	t.Run("create new invitation - success", func(t *testing.T) {
		op, err := New(newMockProvider(), "sample-agent", sampleAppURL)

		require.NoError(t, err)
		require.NotEmpty(t, op)

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodGet, commandName+CreateInvitationPath, nil)

		op.CreateInvitation(rw, rq)
		require.Equal(t, rw.Code, http.StatusOK)
		require.Contains(t, rw.Body.String(), `{"url":"http://demo.wallet.app/home?oob=eyJAaWQiO`)
	})

	t.Run("create new invitation - failure", func(t *testing.T) {
		prov := newMockProvider()
		prov.ServiceMap[outofbandsvc.Name] = &mockprotocol.MockOobService{
			SaveInvitationErr: fmt.Errorf(sampleErr),
		}

		op, err := New(prov, "sample-agent", sampleAppURL)

		require.NoError(t, err)
		require.NotEmpty(t, op)

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodGet, commandName+CreateInvitationPath, nil)

		op.CreateInvitation(rw, rq)
		require.Equal(t, rw.Code, http.StatusInternalServerError)
		fmt.Println(rw.Body.String())
		require.Contains(t, rw.Body.String(), `{"errMessage":"failed to save outofband invitation : sample-error"}`)
	})
}

func TestOperation_RequestApplicationProfile(t *testing.T) {
	op, err := New(newMockProvider(), "test", "demo")

	require.NoError(t, err)
	require.NotEmpty(t, op)

	rw := httptest.NewRecorder()
	rq := httptest.NewRequest(http.MethodGet, commandName+RequestCHAPIAppProfile, nil)
	op.RequestApplicationProfile(rw, rq)
	require.Contains(t, rw.Body.String(), `{"errMessage":"To be implemented !"}`)
	require.Equal(t, rw.Code, http.StatusNotImplemented)
}

func TestOperation_SendCHAPIRequest(t *testing.T) {
	op, err := New(newMockProvider(), "test", "demo")

	require.NoError(t, err)
	require.NotEmpty(t, op)

	rw := httptest.NewRecorder()
	rq := httptest.NewRequest(http.MethodGet, commandName+SendCHAPIRequest, nil)
	op.SendCHAPIRequest(rw, rq)
	require.Contains(t, rw.Body.String(), `{"errMessage":"To be implemented !"}`)
	require.Equal(t, rw.Code, http.StatusNotImplemented)
}

func newMockProvider() *mockprotocol.MockProvider {
	prov := mockprotocol.NewMockProvider()
	prov.ServiceMap = map[string]interface{}{
		outofbandsvc.Name:        &mockprotocol.MockOobService{},
		mediatorsvc.Coordination: &mockmediator.MockMediatorSvc{},
	}

	return prov
}
