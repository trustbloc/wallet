/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sdscomm

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/trustbloc/edv/pkg/restapi/messages"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edv/pkg/edvprovider/memedvprovider"
	"github.com/trustbloc/edv/pkg/restapi"
)

func TestNew(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sdsComm, err := New("SDSUrl", "AgentUsername")
		require.NoError(t, err)
		require.NotNil(t, sdsComm)
	})
	t.Run("Failure: blank SDS URL", func(t *testing.T) {
		sdsComm, err := New("", "AgentUsername")
		require.EqualError(t, err, errBlankSDSURL.Error())
		require.Nil(t, sdsComm)
	})
	t.Run("Failure: blank agent username URL", func(t *testing.T) {
		sdsComm, err := New("SDSUrl", "")
		require.EqualError(t, err, errBlankAgentUsername.Error())
		require.Nil(t, sdsComm)
	})
}

func TestSDSComm_CreateDIDVault(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sdsSrv := newTestEDVServer(t)
		defer sdsSrv.Close()

		sdsComm, err := New(fmt.Sprintf("%s/encrypted-data-vaults", sdsSrv.URL), "James")
		require.NoError(t, err)

		err = sdsComm.CreateDIDVault()
		require.NoError(t, err)
	})
	t.Run("SDS server unreachable (unsupported protocol scheme provided)", func(t *testing.T) {
		sdsComm, err := New("BadURL", "James")
		require.NoError(t, err)

		err = sdsComm.CreateDIDVault()
		require.Contains(t, err.Error(), "unsupported protocol scheme")
	})
}

func TestSDSComm_StoreDIDDocument(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sdsSrv := newTestEDVServer(t)
		defer sdsSrv.Close()

		sdsComm, err := New(fmt.Sprintf("%s/encrypted-data-vaults", sdsSrv.URL), "James")
		require.NoError(t, err)

		err = sdsComm.CreateDIDVault()
		require.NoError(t, err)

		sampleDIDDocData := DIDDocData{}

		err = sdsComm.StoreDIDDocument(&sampleDIDDocData)
		require.NoError(t, err)
	})
	t.Run("Fail to store - vault not found", func(t *testing.T) {
		sdsSrv := newTestEDVServer(t)
		defer sdsSrv.Close()

		sdsComm, err := New(fmt.Sprintf("%s/encrypted-data-vaults", sdsSrv.URL), "James")
		require.NoError(t, err)

		sampleDIDDocData := DIDDocData{}

		err = sdsComm.StoreDIDDocument(&sampleDIDDocData)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), messages.ErrVaultNotFound.Error())
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
