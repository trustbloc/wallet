/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package controller_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edge-agent/pkg/controller"
)

func TestGetCommandHandlers(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		ctrl, err := controller.GetCommandHandlers(controller.WithBlocDomain("domain"))
		require.NoError(t, err)
		require.NotNil(t, ctrl)
	})
}
