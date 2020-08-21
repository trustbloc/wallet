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
		controller, err := GetCommandHandlers(WithBlocDomain("domain"),
			WithSDSServerURL("SomeURL"), WithAgentUsername("agent1"))
		require.NoError(t, err)
		require.NotNil(t, controller)
	})
	t.Run("Fail to instantiate DID client", func(t *testing.T) {
		controller, err := GetCommandHandlers(WithBlocDomain("domain"),
			WithSDSServerURL(""), WithAgentUsername("agent1"))
		require.EqualError(t, err, "failure while creating new DID client: "+
			"failure while preparing SDS communication: SDS server URL cannot be blank")
		require.Nil(t, controller)
	})
}
