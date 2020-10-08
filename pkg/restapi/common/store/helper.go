/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package store

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/trustbloc/edge-core/pkg/storage"
)

// Open the store with the given name. The store will be created if it does not exist.
func Open(p storage.Provider, name string) (storage.Store, error) {
	err := p.CreateStore(name)
	if err != nil && !errors.Is(err, storage.ErrDuplicateStore) {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return p.OpenStore(name)
}

// Save the value mapped to the key in the given store.
func Save(s storage.Store, k string, v interface{}) error {
	bits, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %s", err)
	}

	return s.Put(k, bits)
}
