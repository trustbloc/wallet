/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package credentialclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/trustbloc/edge-agent/pkg/controller/command/sdscomm"
	"github.com/trustbloc/edv/pkg/edvprovider/memedvprovider"
	"github.com/trustbloc/edv/pkg/restapi"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		c, err := New("EDVURL", "AgentUsername")
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NotNil(t, c.GetHandlers())
	})
	t.Run("Fail to instantiate sdscomm", func(t *testing.T) {
		c, err := New("", "AgentUsername")
		require.EqualError(t, err, "failure while preparing SDS communication: "+
			"SDS server URL cannot be blank")
		require.Nil(t, c)
	})
}

func TestCommand_SaveCredential(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sdsSrv := newTestEDVServer(t)
		defer sdsSrv.Close()

		sampleCredentialData := sdscomm.CredentialData{}

		credentialDataBytes, err := json.Marshal(sampleCredentialData)
		require.NoError(t, err)

		cmd, err := New(fmt.Sprintf("%s/encrypted-data-vaults", sdsSrv.URL), "AgentUsername")
		require.NoError(t, err)
		cmdErr := cmd.SaveCredential(nil, bytes.NewBuffer(credentialDataBytes))
		require.NoError(t, cmdErr)
	})
	t.Run("Fail to unmarshal - invalid DIDDocData", func(t *testing.T) {
		cmd, err := New("EDVURL", "AgentUsername")
		require.NoError(t, err)
		cmdErr := cmd.SaveCredential(nil, bytes.NewBuffer([]byte("")))
		require.Contains(t, cmdErr.Error(), failDecodeCredentialDocDataErrMsg)
	})
	t.Run("Fail to save credential - bad SDS server URL", func(t *testing.T) {
		cmd, err := New("BadURL", "AgentUsername")
		require.NoError(t, err)

		sampleCredentialData := sdscomm.CredentialData{}

		credentialDataBytes, err := json.Marshal(sampleCredentialData)
		require.NoError(t, err)

		cmdErr := cmd.SaveCredential(nil, bytes.NewBuffer(credentialDataBytes))
		require.Contains(t, cmdErr.Error(), `failure while storing credential in SDS: `+
			`failure while ensuring that the user's credential vault exists: unexpected error during `+
			`the "create vault" call to SDS: failed to send POST request:`)
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
