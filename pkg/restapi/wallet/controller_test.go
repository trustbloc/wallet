/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package wallet_test

import (
	"testing"

	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries/api"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries/defaults"
	"github.com/hyperledger/aries-framework-go/pkg/framework/context"
	mockstorage "github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	"github.com/stretchr/testify/require"

	"github.com/trustbloc/edge-agent/pkg/restapi/wallet"
)

func TestGetRESTHandlers(t *testing.T) {
	t.Run("test failure", func(t *testing.T) {
		ctx, err := context.New(context.WithStorageProvider(mockstorage.NewMockStoreProvider()),
			context.WithProtocolStateStorageProvider(mockstorage.NewMockStoreProvider()))
		require.NoError(t, err)

		ctrl, err := wallet.GetRESTHandlers(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), api.ErrSvcNotFound.Error())
		require.Nil(t, ctrl)
	})

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
		require.NotEmpty(t, handlers)
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
		require.NotEmpty(t, handlers)
	})
}
