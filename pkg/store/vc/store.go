// +build js,wasm

/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package vc

import (
	"fmt"
	"syscall/js"

	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"

	"github.com/hyperledger/aries-framework-go/pkg/storage"
	verifiablestore "github.com/hyperledger/aries-framework-go/pkg/store/verifiable"
)

type provider interface {
	StorageProvider() storage.Provider
}

type vcStore interface {
	SaveVC(vc *verifiable.Credential) error
}

// RegisterHandleJSCallback register handle vc store js callback
func RegisterHandleJSCallback(ctx provider) error {
	store, err := verifiablestore.New(ctx)
	if err != nil {
		return fmt.Errorf("failed to create new instance of verifiable store: %w", err)
	}

	c := callback{store: store}
	js.Global().Set("storeVC", js.FuncOf(c.storeVC))

	return nil
}

type callback struct {
	store vcStore
}

func (c *callback) storeVC(this js.Value, inputs []js.Value) interface{} {
	// https://github.com/golang/go/issues/26382
	go func() {
		m := make(map[string]interface{})
		m["dataType"] = "Response"
		m["data"] = "success"
		vcData := inputs[0].Get("credential").Get("data").String()
		vc, _, err := verifiable.NewCredential([]byte(vcData))
		if err != nil {
			m["data"] = fmt.Sprintf("failed to create new credential: %s", err.Error())
		}

		if err := c.store.SaveVC(vc); err != nil {
			m["data"] = fmt.Sprintf("failed to put in vc store: %s", err.Error())
		}

		inputs[0].Call("respondWith", js.Global().Get("Promise").Call("resolve", m))
	}()

	return nil
}
