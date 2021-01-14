/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package wallet_test

import (
	"testing"

	"github.com/hyperledger/aries-framework-go/pkg/framework/context"
	"github.com/stretchr/testify/require"

	"github.com/trustbloc/edge-agent/pkg/restapi/wallet"
)

func TestGetRESTHandlers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		_, err := wallet.GetRESTHandlers(&context.Provider{})
		require.NoError(t, err)
	})
}
