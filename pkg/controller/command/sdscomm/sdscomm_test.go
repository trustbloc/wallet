/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sdscomm // nolint:testpackage // uses internal implementation details

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edv/pkg/edvprovider/memedvprovider"
	"github.com/trustbloc/edv/pkg/restapi"
)

const exampleUserID = "JamesBond"

func TestNew(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sdsComm := New("SDSUrl")
		require.NotNil(t, sdsComm)
	})
}

func TestSDSComm_StoreDIDDocument(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sdsSrv := newTestEDVServer(t)
		defer sdsSrv.Close()

		sdsComm := New(fmt.Sprintf("%s/encrypted-data-vaults", sdsSrv.URL))

		err := sdsComm.ensureVaultExists(sdsComm.getDIDVaultID(exampleUserID))
		require.NoError(t, err)

		sampleDIDDocData := SaveDIDDocToSDSRequest{}

		err = sdsComm.StoreDIDDocument(&sampleDIDDocData)
		require.NoError(t, err)
	})
	t.Run("SDS server unreachable (bad SDS server URL)", func(t *testing.T) {
		sdsComm := New("BadURL")

		err := sdsComm.StoreDIDDocument(&SaveDIDDocToSDSRequest{})
		require.Error(t, err)
		require.Contains(t, err.Error(), `unexpected error during the "create vault" call to SDS: `+
			`failed to send POST request:`)
	})
	t.Run("Error while ensuring vault exists - SDS server URL cannot be blank", func(t *testing.T) {
		sdsComm := New("")

		err := sdsComm.StoreDIDDocument(&SaveDIDDocToSDSRequest{})
		require.Error(t, err)
		require.EqualError(t, err,
			fmt.Errorf(failureEnsuringDIDVaultExistsErrMsg, errSDSServerURLBlank).Error())
	})
}

func TestSDSComm_StoreCredential(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sdsSrv := newTestEDVServer(t)
		defer sdsSrv.Close()

		sdsComm := New(fmt.Sprintf("%s/encrypted-data-vaults", sdsSrv.URL))

		err := sdsComm.ensureVaultExists(sdsComm.getCredentialVaultID(exampleUserID))
		require.NoError(t, err)

		sampleCredentialData := SaveCredentialToSDSRequest{}

		err = sdsComm.StoreCredential(&sampleCredentialData)
		require.NoError(t, err)
	})
	t.Run("SDS server unreachable (bad SDS server URL)", func(t *testing.T) {
		sdsComm := New("BadURL")

		err := sdsComm.StoreCredential(&SaveCredentialToSDSRequest{})
		require.Error(t, err)
		require.Contains(t, err.Error(), `unexpected error during the "create vault" call to SDS: `+
			`failed to send POST request:`)
	})
}

func TestSDSComm_StorePresentation(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sdsSrv := newTestEDVServer(t)
		defer sdsSrv.Close()

		sdsComm := New(fmt.Sprintf("%s/encrypted-data-vaults", sdsSrv.URL))

		err := sdsComm.ensureVaultExists(sdsComm.getPresentationVaultID(exampleUserID))
		require.NoError(t, err)

		samplePresentationData := SavePresentationToSDSRequest{}

		err = sdsComm.StorePresentation(&samplePresentationData)
		require.NoError(t, err)
	})
	t.Run("SDS server unreachable (bad SDS server URL)", func(t *testing.T) {
		sdsComm := New("BadURL")

		err := sdsComm.StorePresentation(&SavePresentationToSDSRequest{})
		require.Error(t, err)
		require.Contains(t, err.Error(), `unexpected error during the "create vault" call to SDS: `+
			`failed to send POST request:`)
	})
}

func newTestEDVServer(t *testing.T) *httptest.Server {
	edvService, err := restapi.New(memedvprovider.NewProvider())
	require.NoError(t, err)

	handlers := edvService.GetOperations()
	router := mux.NewRouter()
	router.UseEncodedPath()

	for _, handler := range handlers {
		router.HandleFunc(handler.Path(), handler.Handle()).Methods(handler.Method())
	}

	return httptest.NewServer(router)
}
