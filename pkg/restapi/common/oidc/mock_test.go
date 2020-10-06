/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestMockClient_FormatRequest(t *testing.T) {
	t.Run("returns mock request", func(t *testing.T) {
		expected := uuid.New().String()
		m := &MockClient{AuthRequest: expected}
		require.Equal(t, expected, m.FormatRequest(""))
	})
}

func TestMockClient_Exchange(t *testing.T) {
	t.Run("returns token", func(t *testing.T) {
		expected := &oauth2.Token{
			AccessToken:  uuid.New().String(),
			RefreshToken: uuid.New().String(),
		}
		m := &MockClient{OAuthToken: expected}
		result, err := m.Exchange(nil, "")
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("returns error", func(t *testing.T) {
		expected := errors.New("test")
		m := &MockClient{OAuthErr: expected}
		_, err := m.Exchange(nil, "")
		require.Equal(t, expected, err)
	})
}

func TestMockClient_VerifyIDToken(t *testing.T) {
	t.Run("returns id_token", func(t *testing.T) {
		expected := &MockIDToken{}
		m := &MockClient{IDToken: expected}
		result, err := m.VerifyIDToken(nil, nil)
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("returns error", func(t *testing.T) {
		expected := errors.New("test")
		m := &MockClient{IDTokenErr: expected}
		_, err := m.VerifyIDToken(nil, nil)
		require.Equal(t, expected, err)
	})
}

func TestMockIDToken_Claims(t *testing.T) {
	t.Run("runs ClaimsFunc", func(t *testing.T) {
		executed := false
		m := &MockIDToken{ClaimsFunc: func(interface{}) error {
			executed = true
			return nil
		}}
		err := m.Claims(nil)
		require.NoError(t, err)
		require.True(t, executed)
	})

	t.Run("returns error", func(t *testing.T) {
		expected := errors.New("test")
		m := &MockIDToken{ClaimsErr: expected}
		result := m.Claims(nil)
		require.Equal(t, expected, result)
	})
}
