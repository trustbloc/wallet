// +build js,wasm

/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package vc

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/hyperledger/aries-framework-go/pkg/storage"
)

const (
	friendlyNamePreFix = "fr_"
)

type provider interface {
	StorageProvider() storage.Provider
}

// RegisterHandleJSCallback register handle vc store js callback
func RegisterHandleJSCallback(ctx provider) error {
	vcFriendlyNameStore, err := ctx.StorageProvider().OpenStore("vc-friendlyname")
	if err != nil {
		return fmt.Errorf("failed to open vc frindly name store: %w", err)
	}

	vcStore, err := ctx.StorageProvider().OpenStore("verifiable")
	if err != nil {
		return fmt.Errorf("failed to open vc frindly name store: %w", err)
	}

	c := callback{vcStore: vcStore, vcFriendlyNameStore: vcFriendlyNameStore, jsDoc: js.Global().Get("document")}
	js.Global().Set("storeVC", js.FuncOf(c.storeVC))
	js.Global().Set("populateVC", js.FuncOf(c.populateVC))
	js.Global().Set("getVC", js.FuncOf(c.getVC))

	return nil
}

type callback struct {
	vcStore             storage.Store
	vcFriendlyNameStore storage.Store
	jsDoc               js.Value
}

func (c *callback) storeVC(this js.Value, inputs []js.Value) interface{} {
	// https://github.com/golang/go/issues/26382
	go func() {
		m := make(map[string]interface{})
		m["dataType"] = "Response"
		m["data"] = "success"
		vcData := inputs[0].Get("credential").Get("data").String()

		var vcMap map[string]interface{}
		if err := json.Unmarshal([]byte(vcData), &vcMap); err != nil {
			m["data"] = fmt.Sprintf("failed to unmarshal vc json: %s", err.Error())
			inputs[0].Call("respondWith", js.Global().Get("Promise").Call("resolve", m))
			return
		}

		vcID := stringEntry(vcMap["id"])
		if vcID == "" {
			m["data"] = fmt.Sprintf("vc id not exist")
			inputs[0].Call("respondWith", js.Global().Get("Promise").Call("resolve", m))
			return
		}

		if err := c.vcStore.Put(vcID, []byte(vcData)); err != nil {
			m["data"] = fmt.Sprintf("failed to put in vc store: %s", err.Error())
			inputs[0].Call("respondWith", js.Global().Get("Promise").Call("resolve", m))
			return
		}

		if err := c.vcFriendlyNameStore.Put(friendlyNamePreFix+inputs[1].String(), []byte(vcID)); err != nil {
			m["data"] = fmt.Sprintf("failed to put in vc friendly name store: %s", err.Error())
		}

		inputs[0].Call("respondWith", js.Global().Get("Promise").Call("resolve", m))
	}()

	return nil
}

func (c *callback) populateVC(this js.Value, inputs []js.Value) interface{} {
	// https://github.com/golang/go/issues/26382
	go func() {
		itr := c.vcFriendlyNameStore.Iterator(friendlyNamePreFix, "")
		for itr.Next() {
			optionEL := c.jsDoc.Call("createElement", "option")
			val := strings.TrimPrefix(string(itr.Key()), friendlyNamePreFix)
			optionEL.Set("textContent", val)
			optionEL.Set("value", val)
			inputs[0].Call("appendChild", optionEL)
		}
	}()

	return nil
}

func (c *callback) getVC(this js.Value, inputs []js.Value) interface{} {
	// https://github.com/golang/go/issues/26382
	go func() {
		m := make(map[string]interface{})
		m["dataType"] = "Response"

		vcID, err := c.vcFriendlyNameStore.Get(friendlyNamePreFix + inputs[1].String())
		if err != nil {
			m["data"] = fmt.Sprintf("failed to get in vc friendly name store: %s", err.Error())
			inputs[0].Call("respondWith", js.Global().Get("Promise").Call("resolve", m))
			return
		}

		vcBytes, err := c.vcStore.Get(string(vcID))
		if err != nil {
			m["data"] = fmt.Sprintf("failed to get in vc store: %s", err.Error())
			inputs[0].Call("respondWith", js.Global().Get("Promise").Call("resolve", m))
			return
		}

		m["data"] = string(vcBytes)
		inputs[0].Call("respondWith", js.Global().Get("Promise").Call("resolve", m))
	}()

	return nil
}

// stringEntry
func stringEntry(entry interface{}) string {
	if entry == nil {
		return ""
	}

	return entry.(string)
}
