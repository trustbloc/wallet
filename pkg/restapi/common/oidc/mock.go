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
	IDToken     Claimer
	IDTokenErr  error
	UserInfoVal Claimer
	UserInfoErr error
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
func (m *MockClient) VerifyIDToken(_ context.Context, _ OAuth2Token) (Claimer, error) {
	return m.IDToken, m.IDTokenErr
}

// UserInfo returns the user's info.
func (m *MockClient) UserInfo(_ context.Context, _ *oauth2.Token) (Claimer, error) {
	return m.UserInfoVal, m.UserInfoErr
}

// MockClaimer can be a mock id_token or a mock UserInfo.
type MockClaimer struct {
	ClaimsErr  error
	ClaimsFunc func(interface{}) error
}

// Claims scans the claims into 'i'.
func (m *MockClaimer) Claims(i interface{}) error {
	if m.ClaimsFunc != nil {
		return m.ClaimsFunc(i)
	}

	return m.ClaimsErr
}
