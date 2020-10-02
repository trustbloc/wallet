/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"encoding/json"
	"fmt"
	oidc2 "github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"

	"github.com/trustbloc/edge-core/pkg/storage"
)

// standard claims: https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
type endUser struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
}

func (e *endUser) parse(t oidc2.IDToken) error {
	err := t.Claims(e)
	if err != nil {
		return fmt.Errorf("failed to parse claims from id_token: %w", err)
	}

	err = evaluateClaims(e)
	if err != nil {
		return fmt.Errorf("failed to evaluate claims in id_token: %w", err)
	}

	return nil
}

type endUserTokens struct {
	ID      string
	Access  string
	Refresh string
}

type persistedData struct {
	s storage.Store
}

func newPersistedData(s storage.Store) *persistedData {
	return &persistedData{s: s}
}

func (e *persistedData) put(k string, v interface{}) error {
	bits, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return e.s.Put(k, bits)
}

func (e *persistedData) getEndUser(sub string) (*endUser, error) {
	bits, err := e.s.Get(sub)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user from store: %w", err)
	}

	user := &endUser{}

	return user, json.Unmarshal(bits, user)
}
