/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chapibridge // nolint:testpackage // changing to different package requires exposing internal features.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	didexchangesvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/didexchange"
	mediatorsvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/mediator"
	outofbandsvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/outofband"
	mockdidexchange "github.com/hyperledger/aries-framework-go/pkg/mock/didcomm/protocol/didexchange"
	mockmediator "github.com/hyperledger/aries-framework-go/pkg/mock/didcomm/protocol/mediator"
	mockkms "github.com/hyperledger/aries-framework-go/pkg/mock/kms"
	mockstorage "github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	"github.com/stretchr/testify/require"

	"github.com/trustbloc/edge-agent/pkg/restapi/internal/mocks"
	mockprotocol "github.com/trustbloc/edge-agent/pkg/restapi/internal/mocks/protocol"
)

const (
	sampleAppURL = "http://demo.wallet.app/home"
	sampleErr    = "sample-error"
)

func TestNew(t *testing.T) {
	t.Run("create new instance - success", func(t *testing.T) {
		op, err := New(newMockProvider(), mocks.NewMockNotifier(), "test", "demo")

		require.NoError(t, err)
		require.NotEmpty(t, op)
		require.Len(t, op.GetRESTHandlers(), 3)
	})

	t.Run("create new instance - oob client failure", func(t *testing.T) {
		prov := newMockProvider()
		delete(prov.ServiceMap, outofbandsvc.Name)

		op, err := New(prov, mocks.NewMockNotifier(), "test", "demo")

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create out-of-band client")
		require.Empty(t, op)
	})

	t.Run("create new instance - did exchange client failure", func(t *testing.T) {
		prov := newMockProvider()
		delete(prov.ServiceMap, didexchangesvc.DIDExchange)

		op, err := New(prov, mocks.NewMockNotifier(), "test", "demo")

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create did-exchange client")
		require.Empty(t, op)
	})

	t.Run("create new instance - open store failure", func(t *testing.T) {
		prov := newMockProvider()
		storageProv := mockstorage.NewMockStoreProvider()
		storageProv.ErrOpenStoreHandle = fmt.Errorf(sampleErr)
		prov.StoreProvider = storageProv

		op, err := New(prov, mocks.NewMockNotifier(), "test", "demo")

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to open wallet profile store")
		require.Contains(t, err.Error(), "sample-error")
		require.Empty(t, op)
	})

	t.Run("create new instance - register action event failure", func(t *testing.T) {
		prov := newMockProvider()
		didexSvc := &mockdidexchange.MockDIDExchangeSvc{}
		didexSvc.RegisterActionEventErr = fmt.Errorf(sampleErr)
		prov.ServiceMap[didexchangesvc.DIDExchange] = didexSvc

		op, err := New(prov, mocks.NewMockNotifier(), "test", "demo")

		require.Error(t, err)
		require.Empty(t, op)
		require.Contains(t, err.Error(), "failed to register events")
	})

	t.Run("create new instance - register state event failure", func(t *testing.T) {
		prov := newMockProvider()
		didexSvc := &mockdidexchange.MockDIDExchangeSvc{}
		didexSvc.RegisterMsgEventErr = fmt.Errorf(sampleErr)
		prov.ServiceMap[didexchangesvc.DIDExchange] = didexSvc

		op, err := New(prov, mocks.NewMockNotifier(), "test", "demo")

		require.Error(t, err)
		require.Empty(t, op)
		require.Contains(t, err.Error(), "failed to register events")
	})
}

func TestOperation_CreateInvitation(t *testing.T) {
	const sampleRequest = `{"userID": "1234"}`

	const sampleInvalidRequest = `{"userID": ""}`

	t.Run("create new invitation - success", func(t *testing.T) {
		op, err := New(newMockProvider(), mocks.NewMockNotifier(), "sample-agent", sampleAppURL)

		require.NoError(t, err)
		require.NotEmpty(t, op)

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodPost, commandName+CreateInvitationPath,
			bytes.NewBufferString(sampleRequest))

		op.CreateInvitation(rw, rq)

		require.Equal(t, rw.Code, http.StatusOK)
		require.Contains(t, rw.Body.String(), `{"url":"http://demo.wallet.app/home?oob=eyJAaWQiO`)
	})

	t.Run("create new invitation - failure - oob error", func(t *testing.T) {
		prov := newMockProvider()
		prov.ServiceMap[outofbandsvc.Name] = &mockprotocol.MockOobService{
			SaveInvitationErr: fmt.Errorf(sampleErr),
		}

		op, err := New(prov, mocks.NewMockNotifier(), "sample-agent", sampleAppURL)

		require.NoError(t, err)
		require.NotEmpty(t, op)

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodPost, commandName+CreateInvitationPath,
			bytes.NewBufferString(sampleRequest))

		op.CreateInvitation(rw, rq)
		require.Equal(t, rw.Code, http.StatusInternalServerError)

		require.Contains(t, rw.Body.String(), `{"errMessage":"failed to save outofband invitation : sample-error"}`)
	})

	t.Run("create new invitation - failure - validation error", func(t *testing.T) {
		op, err := New(newMockProvider(), mocks.NewMockNotifier(), "sample-agent", sampleAppURL)

		require.NoError(t, err)
		require.NotEmpty(t, op)

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodPost, commandName+CreateInvitationPath,
			bytes.NewBufferString(sampleInvalidRequest))

		op.CreateInvitation(rw, rq)
		require.Equal(t, rw.Code, http.StatusBadRequest)

		require.Contains(t, rw.Body.String(), invalidIDErr)
	})

	t.Run("create new invitation - failure - invalid request", func(t *testing.T) {
		op, err := New(newMockProvider(), mocks.NewMockNotifier(), "sample-agent", sampleAppURL)

		require.NoError(t, err)
		require.NotEmpty(t, op)

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodPost, commandName+CreateInvitationPath,
			bytes.NewBufferString("-----"))

		op.CreateInvitation(rw, rq)
		require.Equal(t, rw.Code, http.StatusBadRequest)

		require.Contains(t, rw.Body.String(), "invalid character")
	})

	t.Run("create new invitation - failure - save profile error", func(t *testing.T) {
		prov := newMockProvider()
		storageProv := mockstorage.NewMockStoreProvider()
		storageProv.Store.ErrPut = fmt.Errorf(sampleErr)
		prov.StoreProvider = storageProv

		op, err := New(prov, mocks.NewMockNotifier(), "sample-agent", sampleAppURL)

		require.NoError(t, err)
		require.NotEmpty(t, op)

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodPost, commandName+CreateInvitationPath,
			bytes.NewBufferString(sampleRequest))

		op.CreateInvitation(rw, rq)
		require.Equal(t, rw.Code, http.StatusInternalServerError)

		require.Contains(t, rw.Body.String(), `{"errMessage":"failed to save wallet application profile: sample-error"}`)
	})
}

func TestOperation_RequestApplicationProfile(t *testing.T) {
	const pathFmt = "/wallet/%s/request-app-profile"

	t.Run("create application profile - success", func(t *testing.T) {
		op, err := New(newMockProvider(), mocks.NewMockNotifier(), "sample-agent", sampleAppURL)

		require.NoError(t, err)
		require.NotEmpty(t, op)

		sampleProfile := &walletAppProfile{InvitationID: sampleInvID, ConnectionID: sampleConnID}
		err = op.store.SaveProfile(sampleUserID, sampleProfile)
		require.NoError(t, err)

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodGet, fmt.Sprintf(pathFmt, sampleUserID), nil)
		rq = mux.SetURLVars(rq, map[string]string{
			"id": sampleUserID,
		})

		op.RequestApplicationProfile(rw, rq)

		require.Equal(t, rw.Code, http.StatusOK)

		response := applicationProfileResponse{}
		err = json.Unmarshal(rw.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, response.InvitationID, sampleProfile.InvitationID)
		require.Equal(t, response.ConnectionStatus, didexchangesvc.StateIDCompleted)
	})

	t.Run("create application profile - success but status not completed", func(t *testing.T) {
		op, err := New(newMockProvider(), mocks.NewMockNotifier(), "sample-agent", sampleAppURL)

		require.NoError(t, err)
		require.NotEmpty(t, op)

		sampleProfile := &walletAppProfile{InvitationID: sampleInvID}
		err = op.store.SaveProfile(sampleUserID, sampleProfile)
		require.NoError(t, err)

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodGet, fmt.Sprintf(pathFmt, sampleUserID), nil)
		rq = mux.SetURLVars(rq, map[string]string{
			"id": sampleUserID,
		})

		op.RequestApplicationProfile(rw, rq)

		require.Equal(t, rw.Code, http.StatusOK)

		response := applicationProfileResponse{}
		err = json.Unmarshal(rw.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, response.InvitationID, sampleProfile.InvitationID)
		require.Empty(t, response.ConnectionStatus)
	})

	t.Run("create application profile - invalid id", func(t *testing.T) {
		op, err := New(newMockProvider(), mocks.NewMockNotifier(), "sample-agent", sampleAppURL)

		require.NoError(t, err)
		require.NotEmpty(t, op)

		sampleProfile := &walletAppProfile{InvitationID: sampleInvID, ConnectionID: sampleConnID}
		err = op.store.SaveProfile(sampleUserID, sampleProfile)
		require.NoError(t, err)

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodGet, fmt.Sprintf(pathFmt, sampleUserID), nil)

		op.RequestApplicationProfile(rw, rq)

		require.Equal(t, rw.Code, http.StatusBadRequest)
		require.Contains(t, rw.Body.String(), invalidIDErr)
	})

	t.Run("create application profile - profile not found", func(t *testing.T) {
		op, err := New(newMockProvider(), mocks.NewMockNotifier(), "sample-agent", sampleAppURL)
		require.NoError(t, err)
		require.NotEmpty(t, op)

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodGet, fmt.Sprintf(pathFmt, sampleUserID), nil)
		rq = mux.SetURLVars(rq, map[string]string{
			"id": sampleUserID,
		})

		op.RequestApplicationProfile(rw, rq)

		require.Equal(t, rw.Code, http.StatusInternalServerError)
		require.Contains(t, rw.Body.String(), "failed to get wallet application profile by user ID: data not found")
	})

	t.Run("test didexchange completed", func(t *testing.T) {
		op, err := New(newMockProvider(), mocks.NewMockNotifier(), "sample-agent", sampleAppURL)
		require.NoError(t, err)
		require.NotEmpty(t, op)

		sampleProfile := &walletAppProfile{InvitationID: sampleInvID}
		err = op.store.SaveProfile(sampleUserID, sampleProfile)
		require.NoError(t, err)

		ch := make(chan service.StateMsg)

		go op.stateMsgListener(ch)

		ch <- service.StateMsg{
			Type:    service.PostState,
			StateID: didexchangesvc.StateIDCompleted,
			Properties: &mockdidexchange.MockEventProperties{
				InvID:  sampleInvID,
				ConnID: sampleConnID,
			},
		}

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodGet, fmt.Sprintf(pathFmt, sampleUserID), nil)
		rq = mux.SetURLVars(rq, map[string]string{
			"id": sampleUserID,
		})

		op.RequestApplicationProfile(rw, rq)

		require.Equal(t, rw.Code, http.StatusOK)

		response := applicationProfileResponse{}
		err = json.Unmarshal(rw.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, response.InvitationID, sampleProfile.InvitationID)
		require.Equal(t, response.ConnectionStatus, didexchangesvc.StateIDCompleted)
	})

	t.Run("test didexchange completed - but update profile failed", func(t *testing.T) {
		prov := newMockProvider()
		storageProv := mockstorage.NewMockStoreProvider()
		storageProv.Store.ErrPut = fmt.Errorf(sampleErr)
		prov.StoreProvider = storageProv

		op, err := New(prov, mocks.NewMockNotifier(), "sample-agent", sampleAppURL)
		require.NoError(t, err)
		require.NotEmpty(t, op)

		ch := make(chan service.StateMsg)

		go op.stateMsgListener(ch)

		ch <- service.StateMsg{
			Type:    service.PostState,
			StateID: didexchangesvc.StateIDCompleted,
			Properties: &mockdidexchange.MockEventProperties{
				InvID:  sampleInvID,
				ConnID: sampleConnID,
			},
		}

		profile, err := op.store.GetProfileByUserID(sampleUserID)
		require.Error(t, err)
		require.Empty(t, profile)
	})

	t.Run("test didexchange not completed", func(t *testing.T) {
		op, err := New(newMockProvider(), mocks.NewMockNotifier(), "sample-agent", sampleAppURL)
		require.NoError(t, err)
		require.NotEmpty(t, op)

		sampleProfile := &walletAppProfile{InvitationID: sampleInvID}
		err = op.store.SaveProfile(sampleUserID, sampleProfile)
		require.NoError(t, err)

		ch := make(chan service.StateMsg)

		go op.stateMsgListener(ch)

		ch <- service.StateMsg{
			Type:    service.PostState,
			StateID: didexchangesvc.StateIDCompleted,
		}

		ch <- service.StateMsg{
			Type:    service.PostState,
			StateID: didexchangesvc.StateIDRequested,
			Properties: &mockdidexchange.MockEventProperties{
				InvID:  sampleInvID,
				ConnID: sampleConnID,
			},
		}

		rw := httptest.NewRecorder()
		rq := httptest.NewRequest(http.MethodGet, fmt.Sprintf(pathFmt, sampleUserID), nil)
		rq = mux.SetURLVars(rq, map[string]string{
			"id": sampleUserID,
		})

		op.RequestApplicationProfile(rw, rq)

		require.Equal(t, rw.Code, http.StatusOK)

		response := applicationProfileResponse{}
		err = json.Unmarshal(rw.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, response.InvitationID, sampleProfile.InvitationID)
		require.Empty(t, response.ConnectionStatus)
	})
}

func TestOperation_SendCHAPIRequest(t *testing.T) {
	op, err := New(newMockProvider(), mocks.NewMockNotifier(), "test", "demo")

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
		outofbandsvc.Name:          &mockprotocol.MockOobService{},
		mediatorsvc.Coordination:   &mockmediator.MockMediatorSvc{},
		didexchangesvc.DIDExchange: &mockdidexchange.MockDIDExchangeSvc{},
	}

	prov.StoreProvider = mockstorage.NewMockStoreProvider()
	prov.ProtocolStateStoreProvider = mockstorage.NewMockStoreProvider()
	prov.CustomKMS = &mockkms.KeyManager{}

	return prov
}
