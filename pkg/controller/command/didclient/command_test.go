/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package didclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edv/pkg/edvprovider/memedvprovider"
	"github.com/trustbloc/edv/pkg/restapi"
	didclient "github.com/trustbloc/trustbloc-did-method/pkg/did"

	"github.com/trustbloc/edge-agent/pkg/controller/command"
	"github.com/trustbloc/edge-agent/pkg/controller/command/sdscomm"
)

func TestNew(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		c := New("domain", "", "")
		require.NotNil(t, c)
		require.NotNil(t, c.GetHandlers())
	})
}

func TestCommand_CreateDID(t *testing.T) {
	t.Run("test error from request", func(t *testing.T) {
		c := New("domain", "", "")
		require.NotNil(t, c)

		c.didClient = &mockDIDClient{createDIDErr: fmt.Errorf("error create did")}

		var b bytes.Buffer

		cmdErr := c.CreateDID(&b, bytes.NewBufferString("--"))
		require.Error(t, cmdErr)
		require.Equal(t, InvalidRequestErrorCode, cmdErr.Code())
		require.Equal(t, command.ValidationError, cmdErr.Type())
	})

	t.Run("test error from create did", func(t *testing.T) {
		c := New("domain", "", "")
		require.NotNil(t, c)

		c.didClient = &mockDIDClient{createDIDErr: fmt.Errorf("error create did")}

		var b bytes.Buffer

		req, err := json.Marshal(CreateDIDRequest{PublicKeys: []PublicKey{{ID: "key1", Type: "key1",
			Value: base64.RawURLEncoding.EncodeToString([]byte("value"))}}})
		require.NoError(t, err)

		cmdErr := c.CreateDID(&b, bytes.NewBuffer(req))
		require.Error(t, cmdErr)
		require.Equal(t, CreateDIDErrorCode, cmdErr.Code())
		require.Equal(t, command.ExecuteError, cmdErr.Type())
		require.Contains(t, cmdErr.Error(), "error create did")
	})

	t.Run("test error from did base64 decode", func(t *testing.T) {
		c := New("domain", "", "")
		require.NotNil(t, c)

		c.didClient = &mockDIDClient{createDIDErr: fmt.Errorf("error create did")}

		var b bytes.Buffer

		req, err := json.Marshal(CreateDIDRequest{PublicKeys: []PublicKey{{ID: "key1", Type: "key1",
			Value: "value"}}})
		require.NoError(t, err)

		cmdErr := c.CreateDID(&b, bytes.NewBuffer(req))
		require.Error(t, cmdErr)
		require.Equal(t, CreateDIDErrorCode, cmdErr.Code())
		require.Equal(t, command.ExecuteError, cmdErr.Type())
		require.Contains(t, cmdErr.Error(), "illegal base64 data")
	})

	c := New("domain", "", "")
	require.NotNil(t, c)

	c.didClient = &mockDIDClient{createDIDValue: &did.Doc{ID: "1"}}

	var b bytes.Buffer

	t.Run("test success create did with Ed25519 key", func(t *testing.T) {
		// ED key
		r, err := json.Marshal(CreateDIDRequest{PublicKeys: []PublicKey{{ID: "key1", Type: "key1", KeyType: "Ed25519",
			Value: base64.RawURLEncoding.EncodeToString([]byte("value"))}}})
		require.NoError(t, err)

		cmdErr := c.CreateDID(&b, bytes.NewBuffer(r))
		require.NoError(t, cmdErr)

		resp := &CreateDIDResponse{}
		err = json.NewDecoder(&b).Decode(&resp)
		require.NoError(t, err)

		require.Equal(t, "1", resp.DID["id"])
	})
}

func TestCommand_SaveDID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		sdsSrv := newTestEDVServer(t)
		defer sdsSrv.Close()

		sampleDIDDocData := sdscomm.DIDDocData{}

		didDocDataBytes, err := json.Marshal(sampleDIDDocData)
		require.NoError(t, err)

		cmd := New("", fmt.Sprintf("%s/encrypted-data-vaults", sdsSrv.URL), "James")
		cmdErr := cmd.SaveDID(nil, bytes.NewBuffer(didDocDataBytes))
		require.NoError(t, cmdErr)
	})
	t.Run("Fail to unmarshal - invalid DIDDocData", func(t *testing.T) {
		cmd := New("", "", "")
		cmdErr := cmd.SaveDID(nil, bytes.NewBuffer([]byte("")))
		require.Contains(t, cmdErr.Error(), failDecodeDIDDocDataErrMsg)
	})
	t.Run("Fail to save DID - SDS server unreachable", func(t *testing.T) {
		cmd := New("", "BadURL", "agentUsername")

		sampleDIDDocData := sdscomm.DIDDocData{}

		didDocDataBytes, err := json.Marshal(sampleDIDDocData)
		require.NoError(t, err)

		cmdErr := cmd.SaveDID(nil, bytes.NewBuffer(didDocDataBytes))
		require.Contains(t, cmdErr.Error(), failCreateDIDVaultErrMsg)
	})
	t.Run("Fail to save DID - failed to initialize sdscomm", func(t *testing.T) {
		cmd := New("", "", "")

		sampleDIDDocData := sdscomm.DIDDocData{}

		didDocDataBytes, err := json.Marshal(sampleDIDDocData)
		require.NoError(t, err)

		cmdErr := cmd.SaveDID(nil, bytes.NewBuffer(didDocDataBytes))
		require.Contains(t, cmdErr.Error(), failCreateSDSCommErrMsg)
	})
}

type mockDIDClient struct {
	createDIDValue *did.Doc
	createDIDErr   error
}

func (m *mockDIDClient) CreateDID(domain string, opts ...didclient.CreateDIDOption) (*did.Doc, error) {
	return m.createDIDValue, m.createDIDErr
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
