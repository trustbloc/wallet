/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"context"

	"golang.org/x/oauth2"
)

// MockClient is a mock OIDC client.
type MockClient struct {
	AuthRequest string
	OAuthToken  *oauth2.Token
	OAuthErr    error
	IDToken     IDToken
	IDTokenErr  error
}

// FormatRequest formats the OIDC authorization request.
func (m *MockClient) FormatRequest(_ string) string {
	return m.AuthRequest
}

// Exchange exchanges the code for an oauth token.
func (m *MockClient) Exchange(_ context.Context, _ string) (*oauth2.Token, error) {
	return m.OAuthToken, m.OAuthErr
}

// VerifyIDToken verifies the id_token inside the OAuth2 token.
func (m *MockClient) VerifyIDToken(_ context.Context, _ OAuth2Token) (IDToken, error) {
	return m.IDToken, m.IDTokenErr
}

// MockIDToken is a mock OIDC id_token.
type MockIDToken struct {
	ClaimsErr  error
	ClaimsFunc func(interface{}) error
}

// Claims scans the claims into 'i'.
func (m *MockIDToken) Claims(i interface{}) error {
	if m.ClaimsFunc != nil {
		return m.ClaimsFunc(i)
	}

	return m.ClaimsErr
}
