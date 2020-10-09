/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"
	"golang.org/x/oauth2"
)

func TestMockClient_FormatRequest(t *testing.T) {
	t.Run("returns mock request", func(t *testing.T) {
		expected := uuid.New().String()
		m := &oidc.MockClient{AuthRequest: expected}
		require.Equal(t, expected, m.FormatRequest(""))
	})
}

func TestMockClient_Exchange(t *testing.T) {
	t.Run("returns token", func(t *testing.T) {
		expected := &oauth2.Token{
			AccessToken:  uuid.New().String(),
			RefreshToken: uuid.New().String(),
		}
		m := &oidc.MockClient{OAuthToken: expected}
		result, err := m.Exchange(context.TODO(), "")
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("returns error", func(t *testing.T) {
		expected := errors.New("test")
		m := &oidc.MockClient{OAuthErr: expected}
		_, err := m.Exchange(context.TODO(), "")
		require.Equal(t, expected, err)
	})
}

func TestMockClient_VerifyIDToken(t *testing.T) {
	t.Run("returns id_token", func(t *testing.T) {
		expected := &oidc.MockIDToken{}
		m := &oidc.MockClient{IDToken: expected}
		result, err := m.VerifyIDToken(context.TODO(), nil)
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("returns error", func(t *testing.T) {
		expected := errors.New("test")
		m := &oidc.MockClient{IDTokenErr: expected}
		_, err := m.VerifyIDToken(context.TODO(), nil)
		require.Equal(t, expected, err)
	})
}

func TestMockIDToken_Claims(t *testing.T) {
	t.Run("runs ClaimsFunc", func(t *testing.T) {
		executed := false
		m := &oidc.MockIDToken{ClaimsFunc: func(interface{}) error {
			executed = true

			return nil
		}}
		err := m.Claims(nil)
		require.NoError(t, err)
		require.True(t, executed)
	})

	t.Run("returns error", func(t *testing.T) {
		expected := errors.New("test")
		m := &oidc.MockIDToken{ClaimsErr: expected}
		result := m.Claims(nil)
		require.Equal(t, expected, result)
	})
}
