/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package user

import (
	"encoding/json"
	"fmt"

	"github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store"
	"github.com/trustbloc/edge-core/pkg/storage"
)

const (
	// StoreName is the name of the user store.
	StoreName = "edgeagent_users"
)

// User is a user of the wallet.
// The user attributes are based on standard OIDC claims:
// https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims.
type User struct {
	Sub         string `json:"sub"`
	Name        string `json:"name"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	Email       string `json:"email"`
	SecretShare string `json:"secretShare"`
}

// ParseIDToken parses a User from an IDToken.
func ParseIDToken(t oidc.Claimer) (*User, error) {
	user := &User{}

	err := t.Claims(user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse claims from id_token: %w", err)
	}

	err = evaluateClaims(user)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate claims in id_token: %w", err)
	}

	return user, nil
}

// validation rules on the received user claims from the OIDC provider go here.
func evaluateClaims(u *User) error {
	if u.Sub == "" {
		return fmt.Errorf("empty 'sub' in end user claims")
	}

	return nil
}

// NewStore returns a new user Store.
func NewStore(p storage.Provider) (*Store, error) {
	s, err := store.Open(p, StoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to open users store: %w", err)
	}

	return &Store{s: s}, nil
}

// Store stores Users.
type Store struct {
	s storage.Store
}

// Save this user with the user's 'sub' as the key.
func (s *Store) Save(u *User) error {
	return store.Save(s.s, u.Sub, u)
}

// Get the User with the given 'sub'.
func (s *Store) Get(sub string) (*User, error) {
	bits, err := s.s.Get(sub)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user from store: %w", err)
	}

	user := &User{}

	return user, json.Unmarshal(bits, user)
}
