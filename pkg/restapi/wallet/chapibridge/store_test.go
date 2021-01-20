/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chapibridge // nolint:testpackage // changing to different package requires exposing internal features.

import (
	"encoding/json"
	"fmt"
	"testing"

	mockstorage "github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	"github.com/stretchr/testify/require"
)

const (
	sampleStoreErr = "sample error"
	sampleUserID   = "userID-001"
	sampleInvID    = "invID-001"
	sampleConnID   = "connID-001"
)

func TestNewWalletAppProfileStore(t *testing.T) {
	t.Run("create new wallet app profile store - success", func(t *testing.T) {
		store, err := newWalletAppProfileStore(&mockstorage.MockStoreProvider{})

		require.NoError(t, err)
		require.NotEmpty(t, store)
	})

	t.Run("create new wallet app profile store - failure", func(t *testing.T) {
		store, err := newWalletAppProfileStore(&mockstorage.MockStoreProvider{
			ErrOpenStoreHandle: fmt.Errorf(sampleStoreErr),
		})

		require.Error(t, err)
		require.Empty(t, store)
		require.Contains(t, err.Error(), "sample error")
	})
}

func TestWalletAppProfileStore_SaveProfile(t *testing.T) {
	t.Run("save wallet app profile - success", func(t *testing.T) {
		appProfileStore, err := newWalletAppProfileStore(mockstorage.NewMockStoreProvider())

		require.NoError(t, err)
		require.NotEmpty(t, appProfileStore)

		err = appProfileStore.SaveProfile(sampleUserID, &walletAppProfile{InvitationID: sampleInvID})
		require.NoError(t, err)

		userIDBytes, err := appProfileStore.store.Get(getInvitationKeyPrefix(sampleInvID))
		require.NoError(t, err)
		require.Equal(t, string(userIDBytes), sampleUserID)

		profileIDBytes, err := appProfileStore.store.Get(getUserIDKeyPrefix(sampleUserID))
		require.NoError(t, err)

		var profile walletAppProfile
		err = json.Unmarshal(profileIDBytes, &profile)
		require.NoError(t, err)
		require.Equal(t, profile.InvitationID, sampleInvID)
	})

	t.Run("save wallet app profile - failure", func(t *testing.T) {
		provider := mockstorage.NewMockStoreProvider()
		provider.Store = &mockstorage.MockStore{
			ErrPut: fmt.Errorf(sampleStoreErr),
		}

		appProfileStore, err := newWalletAppProfileStore(provider)

		require.NoError(t, err)
		require.NotEmpty(t, appProfileStore)

		err = appProfileStore.SaveProfile(sampleUserID, &walletAppProfile{InvitationID: sampleInvID})
		require.Error(t, err)
		require.Contains(t, err.Error(), sampleStoreErr)
	})
}

func TestWalletAppProfileStore_UpdateProfile(t *testing.T) {
	t.Run("update wallet app profile - success", func(t *testing.T) {
		appProfileStore, err := newWalletAppProfileStore(mockstorage.NewMockStoreProvider())

		require.NoError(t, err)
		require.NotEmpty(t, appProfileStore)

		err = appProfileStore.SaveProfile(sampleUserID, &walletAppProfile{InvitationID: sampleInvID})
		require.NoError(t, err)

		err = appProfileStore.UpdateProfile(&walletAppProfile{
			InvitationID: sampleInvID,
			ConnectionID: sampleConnID,
		})
		require.NoError(t, err)

		userIDBytes, err := appProfileStore.store.Get(getInvitationKeyPrefix(sampleInvID))
		require.NoError(t, err)
		require.Equal(t, string(userIDBytes), sampleUserID)

		profileIDBytes, err := appProfileStore.store.Get(getUserIDKeyPrefix(sampleUserID))
		require.NoError(t, err)

		var profile walletAppProfile
		err = json.Unmarshal(profileIDBytes, &profile)
		require.NoError(t, err)
		require.Equal(t, profile.InvitationID, sampleInvID)
		require.Equal(t, profile.ConnectionID, sampleConnID)
	})

	t.Run("update wallet app profile - failure - user profile missing", func(t *testing.T) {
		appProfileStore, err := newWalletAppProfileStore(mockstorage.NewMockStoreProvider())

		require.NoError(t, err)
		require.NotEmpty(t, appProfileStore)

		err = appProfileStore.UpdateProfile(&walletAppProfile{
			InvitationID: sampleInvID,
			ConnectionID: sampleConnID,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get user profile for given wallet profile")
	})

	t.Run("update wallet app profile - failure - store error", func(t *testing.T) {
		provider := mockstorage.NewMockStoreProvider()
		provider.Store = &mockstorage.MockStore{
			ErrPut: fmt.Errorf(sampleStoreErr),
			Store: map[string][]byte{
				getInvitationKeyPrefix(sampleInvID): []byte(sampleUserID),
			},
		}

		appProfileStore, err := newWalletAppProfileStore(provider)

		require.NoError(t, err)
		require.NotEmpty(t, appProfileStore)

		err = appProfileStore.UpdateProfile(&walletAppProfile{
			InvitationID: sampleInvID,
			ConnectionID: sampleConnID,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), sampleStoreErr)
		require.Contains(t, err.Error(), "failed to update wallet application profile in store")
	})
}

func TestWalletAppProfileStore_Get(t *testing.T) {
	t.Run("get wallet app profile - success", func(t *testing.T) {
		appProfileStore, err := newWalletAppProfileStore(mockstorage.NewMockStoreProvider())

		require.NoError(t, err)
		require.NotEmpty(t, appProfileStore)

		storedProfile := &walletAppProfile{InvitationID: sampleInvID}
		err = appProfileStore.SaveProfile(sampleUserID, storedProfile)
		require.NoError(t, err)

		profile, err := appProfileStore.GetProfileByUserID(sampleUserID)
		require.NoError(t, err)
		require.Equal(t, profile, storedProfile)

		profile, err = appProfileStore.GetProfileByInvitationID(sampleInvID)
		require.NoError(t, err)
		require.Equal(t, profile, storedProfile)

		storedProfile.ConnectionID = sampleConnID
		err = appProfileStore.UpdateProfile(storedProfile)
		require.NoError(t, err)

		profile, err = appProfileStore.GetProfileByUserID(sampleUserID)
		require.NoError(t, err)
		require.Equal(t, profile, storedProfile)

		profile, err = appProfileStore.GetProfileByInvitationID(sampleInvID)
		require.NoError(t, err)
		require.Equal(t, profile, storedProfile)
	})

	t.Run("get wallet app profile  - failure - store error", func(t *testing.T) {
		provider := mockstorage.NewMockStoreProvider()
		provider.Store = &mockstorage.MockStore{
			ErrGet: fmt.Errorf(sampleStoreErr),
		}

		appProfileStore, err := newWalletAppProfileStore(provider)
		require.NoError(t, err)

		profile, err := appProfileStore.GetProfileByUserID(sampleUserID)
		require.Error(t, err)
		require.Empty(t, profile)
		require.Contains(t, err.Error(), sampleStoreErr)

		profile, err = appProfileStore.GetProfileByInvitationID(sampleInvID)
		require.Error(t, err)
		require.Empty(t, profile)
		require.Contains(t, err.Error(), sampleStoreErr)
	})

	t.Run("get wallet app profile  - failure - invalid data", func(t *testing.T) {
		provider := mockstorage.NewMockStoreProvider()
		provider.Store = &mockstorage.MockStore{
			Store: map[string][]byte{
				getUserIDKeyPrefix(sampleUserID): []byte("--"),
			},
		}

		appProfileStore, err := newWalletAppProfileStore(provider)
		require.NoError(t, err)

		profile, err := appProfileStore.GetProfileByUserID(sampleUserID)
		require.Error(t, err)
		require.Empty(t, profile)
		require.Contains(t, err.Error(), "failed to get wallet application profile bytes")
	})
}

func TestWalletAppProfileStore_putProfileInStore(t *testing.T) {
	t.Run("get wallet app profile  - failure - invalid data", func(t *testing.T) {
		appProfileStore, err := newWalletAppProfileStore(mockstorage.NewMockStoreProvider())
		require.NoError(t, err)

		err = appProfileStore.putProfileInStore(getUserIDKeyPrefix, sampleUserID, make(chan int))
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get wallet application profile bytes")
	})
}
