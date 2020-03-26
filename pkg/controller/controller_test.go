/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCommandHandlers(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		controller, err := GetCommandHandlers(WithBlocDomain("domain"))
		require.NoError(t, err)
		require.NotNil(t, controller)
	})
}
