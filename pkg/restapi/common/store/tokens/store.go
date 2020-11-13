/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tokens

import (
	"encoding/json"
	"fmt"

	"github.com/trustbloc/edge-agent/pkg/restapi/common/store"
	"github.com/trustbloc/edge-core/pkg/storage"
)

const (
	// StoreName is the name of the token store.
	StoreName = "edgeagent_tks"
)

// UserTokens are the tokens associated to a User.
type UserTokens struct {
	UserSub string
	Access  string
	Refresh string
}

// NewStore returns a new token Store.
func NewStore(p storage.Provider) (*Store, error) {
	s, err := store.Open(p, StoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to open tokens store: %w", err)
	}

	return &Store{s: s}, nil
}

// Store holds UserTokens.
type Store struct {
	s storage.Store
}

// Save the UserTokens to the store.
func (s *Store) Save(ut *UserTokens) error {
	return store.Save(s.s, ut.UserSub, ut)
}

// Get fetches a UserTokens from the underlying storage.
func (s *Store) Get(sub string) (*UserTokens, error) {
	raw, err := s.s.Get(sub)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user tokens from store: %w", err)
	}

	tokens := &UserTokens{}

	return tokens, json.Unmarshal(raw, tokens)
}
