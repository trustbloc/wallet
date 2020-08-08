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

func TestSDSComm_CreateDIDVault(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sdsSrv := newTestEDVServer(t)
		defer sdsSrv.Close()

		sdsComm := New(fmt.Sprintf("%s/encrypted-data-vaults", sdsSrv.URL), "James")

		err := sdsComm.CreateDIDVault()
		require.NoError(t, err)
	})
	t.Run("SDS server unreachable (unsupported protocol scheme provided)", func(t *testing.T) {
		sdsComm := New("BadURL", "James")

		err := sdsComm.CreateDIDVault()
		require.Contains(t, err.Error(), "unsupported protocol scheme")
	})
}

func TestSDSComm_StoreDIDDocument(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sdsSrv := newTestEDVServer(t)
		defer sdsSrv.Close()

		sdsComm := New(fmt.Sprintf("%s/encrypted-data-vaults", sdsSrv.URL), "James")

		err := sdsComm.CreateDIDVault()
		require.NoError(t, err)

		sampleDIDDocData := DIDDocData{}

		err = sdsComm.StoreDIDDocument(&sampleDIDDocData)
		require.NoError(t, err)
	})
	t.Run("Fail to store - vault not found", func(t *testing.T) {
		sdsSrv := newTestEDVServer(t)
		defer sdsSrv.Close()

		sdsComm := New(fmt.Sprintf("%s/encrypted-data-vaults", sdsSrv.URL), "James")

		sampleDIDDocData := DIDDocData{}

		err := sdsComm.StoreDIDDocument(&sampleDIDDocData)
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
