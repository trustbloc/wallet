/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"encoding/json"
	"fmt"
	"github.com/trustbloc/edge-core/pkg/storage"
)

type transientData struct {
	s storage.Store
}

func newTransientData(store storage.Store) *transientData {
	return &transientData{s: store}
}

func (t *transientData) Put(k string, v interface{}) error {
	bits, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return t.s.Put(k, bits)
}
