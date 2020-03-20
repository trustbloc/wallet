/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package controller

import (
	"fmt"
	"testing"

	"github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	"github.com/stretchr/testify/require"
)

func TestGetCommandHandlers(t *testing.T) {
	t.Run("test error from initialized didclient", func(t *testing.T) {
		controller, err := GetCommandHandlers(&storage.MockStoreProvider{
			ErrOpenStoreHandle: fmt.Errorf("error open store")})
		require.Error(t, err)
		require.Contains(t, err.Error(), "error open store")
		require.Nil(t, controller)
	})

	t.Run("test success", func(t *testing.T) {
		controller, err := GetCommandHandlers(&storage.MockStoreProvider{}, WithBlocDomain("domain"))
		require.NoError(t, err)
		require.NotNil(t, controller)
	})
}
