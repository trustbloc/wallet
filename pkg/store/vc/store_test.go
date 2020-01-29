// +build js,wasm

/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package vc

import (
	"fmt"
	"syscall/js"
	"testing"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/hyperledger/aries-framework-go/pkg/storage"
	"github.com/stretchr/testify/require"
	mockstorage "github.com/trustbloc/edge-agent/pkg/internal/mock/storage"
)

//nolint:gochecknoglobals,lll
var credential = `
{
  "@context": [
    "https://www.w3.org/2018/credentials/v1",
    "https://www.w3.org/2018/credentials/examples/v1"
  ],
  "credentialSchema": [],
  "credentialSubject": {
    "degree": {
      "type": "BachelorDegree",
      "university": "MIT"
    },
    "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
    "name": "Jayden Doe",
    "spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
  },
  "expirationDate": "2020-01-01T19:23:24Z",
  "id": "http://example.edu/credentials/1872",
  "issuanceDate": "2009-01-01T19:23:24Z",
  "issuer": {
    "id": "did:example:76e12ec712ebc6f1c221ebfeb1f",
    "name": "Example University"
  },
  "referenceNumber": 83294849,
  "type": [
    "VerifiableCredential",
    "UniversityDegreeCredential"
  ]
}
`

func TestRegisterHandleInvitationJSCallback(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		require.NoError(t, RegisterHandleJSCallback(&mockProvider{
			storageProviderValue: &mockstorage.MockStoreProvider{}}))
		require.True(t, js.Global().Get("storeVC").Truthy())
	})

	t.Run("test error from open store", func(t *testing.T) {
		err := RegisterHandleJSCallback(&mockProvider{
			storageProviderValue: &mockstorage.MockStoreProvider{
				ErrOpenStoreHandle: fmt.Errorf("open store error")}})
		require.Error(t, err)
		require.Contains(t, err.Error(), "open store error")
	})
}

func TestStoreVC(t *testing.T) {
	t.Run("test error from new credential", func(t *testing.T) {
		c := callback{store: &mockVCStore{}}
		m := make(map[string]interface{})
		m1 := make(map[string]interface{})
		m1["data"] = "wrongVC"
		m["credential"] = m1
		vcEvent := js.ValueOf(m)
		vcEvent.Set("respondWith", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return nil
		}))

		m2 := make(map[string]interface{})
		valueFlag := make(chan js.Value, 2)
		m2["resolve"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			valueFlag <- args[0]
			return nil
		})
		js.Global().Set("Promise", js.ValueOf(m2))
		c.storeVC(js.Value{}, []js.Value{vcEvent})

		time.Sleep(500 * time.Millisecond)
		select {
		case value := <-valueFlag:
			require.Equal(t, "Response", value.Get("dataType").String())
			require.Contains(t, value.Get("data").String(), "failed to create new credential")
		case <-time.After(2 * time.Second):
			require.Fail(t, "timeout waiting for response")
		}
	})

	t.Run("test error from vc store", func(t *testing.T) {
		c := callback{store: &mockVCStore{saveVCError: fmt.Errorf("put error")}}
		m := make(map[string]interface{})
		m1 := make(map[string]interface{})
		m1["data"] = credential
		m["credential"] = m1
		vcEvent := js.ValueOf(m)
		vcEvent.Set("respondWith", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return nil
		}))

		m2 := make(map[string]interface{})
		valueFlag := make(chan js.Value, 2)
		m2["resolve"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			valueFlag <- args[0]
			return nil
		})
		js.Global().Set("Promise", js.ValueOf(m2))
		c.storeVC(js.Value{}, []js.Value{vcEvent})

		time.Sleep(500 * time.Millisecond)
		select {
		case value := <-valueFlag:
			require.Equal(t, "Response", value.Get("dataType").String())
			require.Contains(t, value.Get("data").String(), "failed to put in vc store")
		case <-time.After(2 * time.Second):
			require.Fail(t, "timeout waiting for response")
		}
	})

	t.Run("test error from vc friendly name store", func(t *testing.T) {
		c := callback{store: &mockVCStore{}, vcFriendlyNameStore: &mockstorage.MockStore{ErrPut: fmt.Errorf("error put")}}
		m := make(map[string]interface{})
		m1 := make(map[string]interface{})
		m1["data"] = credential
		m["credential"] = m1
		vcEvent := js.ValueOf(m)
		vcEvent.Set("respondWith", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return nil
		}))

		m2 := make(map[string]interface{})
		valueFlag := make(chan js.Value, 2)
		m2["resolve"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			valueFlag <- args[0]
			return nil
		})
		js.Global().Set("Promise", js.ValueOf(m2))
		c.storeVC(js.Value{}, []js.Value{vcEvent, js.ValueOf("vc1")})

		time.Sleep(500 * time.Millisecond)
		select {
		case value := <-valueFlag:
			require.Equal(t, "Response", value.Get("dataType").String())
			require.Contains(t, value.Get("data").String(), "failed to put in vc friendly name store")
		case <-time.After(2 * time.Second):
			require.Fail(t, "timeout waiting for response")
		}
	})

	t.Run("test success", func(t *testing.T) {
		c := callback{store: &mockVCStore{}, vcFriendlyNameStore: &mockstorage.MockStore{
			Store: make(map[string][]byte)}}
		m := make(map[string]interface{})
		m1 := make(map[string]interface{})
		m1["data"] = credential
		m["credential"] = m1
		vcEvent := js.ValueOf(m)
		vcEvent.Set("respondWith", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return nil
		}))

		m2 := make(map[string]interface{})
		valueFlag := make(chan js.Value, 2)
		m2["resolve"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			valueFlag <- args[0]
			return nil
		})
		js.Global().Set("Promise", js.ValueOf(m2))
		c.storeVC(js.Value{}, []js.Value{vcEvent, js.ValueOf("vc1")})

		time.Sleep(500 * time.Millisecond)
		select {
		case value := <-valueFlag:
			require.Equal(t, "Response", value.Get("dataType").String())
			require.Equal(t, "success", value.Get("data").String())
		case <-time.After(2 * time.Second):
			require.Fail(t, "timeout waiting for response")
		}
	})
}

func TestPopulateVC(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		s := make(map[string][]byte)
		js.Global().Set("document1", make(map[string]interface{}))
		c := callback{store: &mockVCStore{}, vcFriendlyNameStore: &mockstorage.MockStore{
			Store: s}, jsDoc: js.Global().Get("document")}
		c.jsDoc.Set("createElement", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			m := make(map[string]interface{})
			return m
		}))

		s[friendlyNamePreFix+"key"] = []byte("vcKey")
		valueFlag := make(chan js.Value, 1)
		selectEL := make(map[string]interface{})
		selectEL["appendChild"] = js.ValueOf(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			valueFlag <- args[0]
			return nil
		}))

		c.populateVC(js.Value{}, []js.Value{js.ValueOf(selectEL)})

		time.Sleep(500 * time.Millisecond)
		select {
		case value := <-valueFlag:
			require.Equal(t, "key", value.Get("value").String())
			require.Equal(t, "key", value.Get("textContent").String())
		case <-time.After(2 * time.Second):
			require.Fail(t, "timeout waiting for response")
		}
	})
}

func TestGetVC(t *testing.T) {
	t.Run("test error from vc friendly name store", func(t *testing.T) {
		c := callback{store: &mockVCStore{}, vcFriendlyNameStore: &mockstorage.MockStore{ErrGet: fmt.Errorf("error get")}}

		m := make(map[string]interface{})
		vcEvent := js.ValueOf(m)
		vcEvent.Set("respondWith", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return nil
		}))

		m2 := make(map[string]interface{})
		valueFlag := make(chan js.Value, 1)
		m2["resolve"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			valueFlag <- args[0]
			return nil
		})
		js.Global().Set("Promise", js.ValueOf(m2))
		c.getVC(js.Value{}, []js.Value{vcEvent, js.ValueOf("vc1")})

		time.Sleep(500 * time.Millisecond)
		select {
		case value := <-valueFlag:
			require.Equal(t, "Response", value.Get("dataType").String())
			require.Contains(t, value.Get("data").String(), "failed to get in vc friendly name store")
		case <-time.After(2 * time.Second):
			require.Fail(t, "timeout waiting for response")
		}
	})

	t.Run("test error from vc store", func(t *testing.T) {
		s := make(map[string][]byte)
		c := callback{store: &mockVCStore{getVCFunc: func(vcID string) (v *verifiable.Credential, err error) {
			return nil, fmt.Errorf("vc store error")
		}},
			vcFriendlyNameStore: &mockstorage.MockStore{Store: s}}

		s[friendlyNamePreFix+"key"] = []byte("vcKey")

		m := make(map[string]interface{})
		vcEvent := js.ValueOf(m)
		vcEvent.Set("respondWith", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return nil
		}))

		m2 := make(map[string]interface{})
		valueFlag := make(chan js.Value, 1)
		m2["resolve"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			valueFlag <- args[0]
			return nil
		})
		js.Global().Set("Promise", js.ValueOf(m2))
		c.getVC(js.Value{}, []js.Value{vcEvent, js.ValueOf("key")})

		time.Sleep(500 * time.Millisecond)
		select {
		case value := <-valueFlag:
			require.Equal(t, "Response", value.Get("dataType").String())
			require.Contains(t, value.Get("data").String(), "vc store error")
		case <-time.After(2 * time.Second):
			require.Fail(t, "timeout waiting for response")
		}
	})

	t.Run("test success", func(t *testing.T) {
		s := make(map[string][]byte)
		c := callback{store: &mockVCStore{getVCFunc: func(vcID string) (v *verifiable.Credential, err error) {
			return &verifiable.Credential{ID: "vcID"}, nil
		}},
			vcFriendlyNameStore: &mockstorage.MockStore{Store: s}}

		s[friendlyNamePreFix+"key"] = []byte("vcKey")

		m := make(map[string]interface{})
		vcEvent := js.ValueOf(m)
		vcEvent.Set("respondWith", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return nil
		}))

		m2 := make(map[string]interface{})
		valueFlag := make(chan js.Value, 1)
		m2["resolve"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			valueFlag <- args[0]
			return nil
		})
		js.Global().Set("Promise", js.ValueOf(m2))
		c.getVC(js.Value{}, []js.Value{vcEvent, js.ValueOf("key")})

		time.Sleep(500 * time.Millisecond)
		select {
		case value := <-valueFlag:
			require.Equal(t, "Response", value.Get("dataType").String())
			require.Contains(t, value.Get("data").String(), "vcID")
		case <-time.After(2 * time.Second):
			require.Fail(t, "timeout waiting for response")
		}
	})
}

type mockProvider struct {
	storageProviderValue storage.Provider
}

func (m *mockProvider) StorageProvider() storage.Provider {
	return m.storageProviderValue
}

type mockVCStore struct {
	saveVCError error
	getVCFunc   func(vcID string) (*verifiable.Credential, error)
}

func (m *mockVCStore) SaveVC(vc *verifiable.Credential) error {
	return m.saveVCError
}

func (m *mockVCStore) GetVC(vcID string) (*verifiable.Credential, error) {
	return m.getVCFunc(vcID)
}
