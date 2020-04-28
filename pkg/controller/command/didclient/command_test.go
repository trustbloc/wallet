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
	"testing"

	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/stretchr/testify/require"
	didclient "github.com/trustbloc/trustbloc-did-method/pkg/did"

	"github.com/trustbloc/edge-agent/pkg/controller/command"
)

func TestNew(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		c := New("domain")
		require.NotNil(t, c)
		require.NotNil(t, c.GetHandlers())
	})
}

func TestCommand_CreateDID(t *testing.T) {
	t.Run("test error from request", func(t *testing.T) {
		c := New("domain")
		require.NotNil(t, c)

		c.client = &mockDIDClient{createDIDErr: fmt.Errorf("error create did")}

		var b bytes.Buffer

		cmdErr := c.CreateDID(&b, bytes.NewBufferString("--"))
		require.Error(t, cmdErr)
		require.Equal(t, InvalidRequestErrorCode, cmdErr.Code())
		require.Equal(t, command.ValidationError, cmdErr.Type())
	})

	t.Run("test error from create did", func(t *testing.T) {
		c := New("domain")
		require.NotNil(t, c)

		c.client = &mockDIDClient{createDIDErr: fmt.Errorf("error create did")}

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
		c := New("domain")
		require.NotNil(t, c)

		c.client = &mockDIDClient{createDIDErr: fmt.Errorf("error create did")}

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

	c := New("domain")
	require.NotNil(t, c)

	c.client = &mockDIDClient{createDIDValue: &did.Doc{ID: "1"}}

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

type mockDIDClient struct {
	createDIDValue *did.Doc
	createDIDErr   error
}

func (m *mockDIDClient) CreateDID(domain string, opts ...didclient.CreateDIDOption) (*did.Doc, error) {
	return m.createDIDValue, m.createDIDErr
}
