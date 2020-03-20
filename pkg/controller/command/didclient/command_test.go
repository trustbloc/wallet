/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package didclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edge-agent/pkg/controller/command"
)

func TestNew(t *testing.T) {
	t.Run("test error from initialized kms", func(t *testing.T) {
		c, err := New(&storage.MockStoreProvider{
			ErrOpenStoreHandle: fmt.Errorf("error open store")}, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "error open store")
		require.Nil(t, c)
	})

	t.Run("test success", func(t *testing.T) {
		c, err := New(&storage.MockStoreProvider{}, "domain")
		require.NoError(t, err)
		require.NotNil(t, c)
		require.NotNil(t, c.GetHandlers())
	})
}

func TestCommand_CreateDID(t *testing.T) {
	t.Run("test error from create did", func(t *testing.T) {
		c, err := New(&storage.MockStoreProvider{}, "domain")
		require.NoError(t, err)
		require.NotNil(t, c)

		c.client = &mockDIDClient{createDIDErr: fmt.Errorf("error create did")}

		var b bytes.Buffer

		cmdErr := c.CreateDID(&b, nil)
		require.Error(t, cmdErr)
		require.Equal(t, CreateDIDErrorCode, cmdErr.Code())
		require.Equal(t, command.ExecuteError, cmdErr.Type())
	})
	t.Run("test success", func(t *testing.T) {
		c, err := New(&storage.MockStoreProvider{}, "domain")
		require.NoError(t, err)
		require.NotNil(t, c)

		c.client = &mockDIDClient{createDIDValue: &did.Doc{ID: "1"}}

		var b bytes.Buffer

		cmdErr := c.CreateDID(&b, nil)
		require.NoError(t, cmdErr)

		didDoc := &did.Doc{}
		err = json.NewDecoder(&b).Decode(&didDoc)
		require.NoError(t, err)

		require.Equal(t, "1", didDoc.ID)
	})
}

type mockDIDClient struct {
	createDIDValue *did.Doc
	createDIDErr   error
}

func (m *mockDIDClient) CreateDID(domain string) (*did.Doc, error) {
	return m.createDIDValue, m.createDIDErr
}
