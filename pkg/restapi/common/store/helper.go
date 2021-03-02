/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package store

import (
	"encoding/json"
	"fmt"

	ariesstorage "github.com/hyperledger/aries-framework-go/spi/storage"
)

// Open the store with the given name. The store will be created if it does not exist.
func Open(p ariesstorage.Provider, name string) (ariesstorage.Store, error) {
	return p.OpenStore(name)
}

// Save the value mapped to the key in the given store.
func Save(s ariesstorage.Store, k string, v interface{}) error {
	bits, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %s", err)
	}

	return s.Put(k, bits)
}
