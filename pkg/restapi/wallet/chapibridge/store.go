/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chapibridge

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/storage"
)

const (
	storageNamespace = "walletappprofile"
	invitationKeyFmt = "inv_%s"
	userIDKeyFmt     = "usr_%s"
)

// walletAppProfile is wallet application profile.
type walletAppProfile struct {
	InvitationID string `json:"invitationID"`
	ConnectionID string `json:"connectionID"`
}

// walletAppProfileStore is wallet application profile store.
type walletAppProfileStore struct {
	store storage.Store
}

func newWalletAppProfileStore(p storage.Provider) (*walletAppProfileStore, error) {
	store, err := p.OpenStore(storageNamespace)
	if err != nil {
		return nil, fmt.Errorf("failed to open wallet applcation profile store: %w", err)
	}

	return &walletAppProfileStore{store}, nil
}

// SaveUserProfile saves wallet app profile by user ID.
// Due to very frequent get wallet profile by user ID call data model will be saved in below format.
// UserID --> WalletAppProfile{} & InvitationID --> UserID.
// For now, Invitation will have one-to-one relationship with User Profile.
func (w *walletAppProfileStore) SaveProfile(userID string, profile *walletAppProfile) error {
	err := w.putProfileInStore(getUserIDKeyPrefix, userID, profile)
	if err != nil {
		return fmt.Errorf("failed to save wallet application profile: %w", err)
	}

	err = w.store.Put(getInvitationKeyPrefix(profile.InvitationID), []byte(userID))
	if err != nil {
		return fmt.Errorf("failed to save wallet application profile: %w", err)
	}

	return nil
}

// UpdateProfile updates wallet app profile in underlying store.
// returns error if no existing mapping found with any user profile.
func (w *walletAppProfileStore) UpdateProfile(profile *walletAppProfile) error {
	userIDBytes, err := w.store.Get(getInvitationKeyPrefix(profile.InvitationID))
	if err != nil {
		return fmt.Errorf("failed to get user profile for given wallet profile from store: %w", err)
	}

	err = w.putProfileInStore(getUserIDKeyPrefix, string(userIDBytes), profile)
	if err != nil {
		return fmt.Errorf("failed to update wallet application profile in store: %w", err)
	}

	return nil
}

// GetProfileByInvitationID returns wallet application profile by given invitation ID
// returns error if no existing mapping found with any user profile.
func (w *walletAppProfileStore) GetProfileByInvitationID(invitationID string) (*walletAppProfile, error) {
	userIDBytes, err := w.store.Get(getInvitationKeyPrefix(invitationID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile for given invitationID: %w", err)
	}

	return w.GetProfileByUserID(string(userIDBytes))
}

// GetProfileByUserID returns wallet application profile by given user profile ID.
// returns error if no existing mapping found with any user profile.
func (w *walletAppProfileStore) GetProfileByUserID(userID string) (*walletAppProfile, error) {
	profileBytes, err := w.store.Get(getUserIDKeyPrefix(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet application profile by user ID: %w", err)
	}

	var profile walletAppProfile

	err = json.Unmarshal(profileBytes, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet application profile bytes: %w", err)
	}

	return &profile, nil
}

func (w *walletAppProfileStore) putProfileInStore(prefix func(string) string, key string, profile interface{}) error {
	profileBytes, err := json.Marshal(profile)
	if err != nil {
		return fmt.Errorf("failed to get wallet application profile bytes: %w", err)
	}

	return w.store.Put(prefix(key), profileBytes)
}

// getInvitationKeyPrefix is key prefix for wallet application profile invitation key.
func getInvitationKeyPrefix(invitationID string) string {
	return fmt.Sprintf(invitationKeyFmt, invitationID)
}

// getUserIDKeyPrefix is key prefix for wallet application profile user ID key.
func getUserIDKeyPrefix(userID string) string {
	return fmt.Sprintf(userIDKeyFmt, userID)
}
