/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package wallet_test

import (
	"testing"

	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries/defaults"
	"github.com/stretchr/testify/require"

	"github.com/trustbloc/wallet/pkg/restapi/wallet"
)

func TestGetRESTHandlers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		framework, err := aries.New(defaults.WithInboundHTTPAddr(":26508", "", "", ""))
		require.NoError(t, err)
		require.NotNil(t, framework)

		defer func() { require.NoError(t, framework.Close()) }()

		ctx, err := framework.Context()
		require.NoError(t, err)
		require.NotNil(t, ctx)

		handlers, err := wallet.GetRESTHandlers(ctx)
		require.NoError(t, err)
		require.Empty(t, handlers)
	})

	t.Run("with options", func(t *testing.T) {
		framework, err := aries.New(defaults.WithInboundHTTPAddr(":26508", "", "", ""))
		require.NoError(t, err)
		require.NotNil(t, framework)

		defer func() { require.NoError(t, framework.Close()) }()

		ctx, err := framework.Context()
		require.NoError(t, err)
		require.NotNil(t, ctx)

		handlers, err := wallet.GetRESTHandlers(ctx, wallet.WithWalletAppURL("demoapp"),
			wallet.WithWebhookURLs("demoURL"), wallet.WithNotifier(nil),
			wallet.WithMessageHandler(nil), wallet.WithDefaultLabel("test"))
		require.NoError(t, err)
		require.Empty(t, handlers)
	})
}
