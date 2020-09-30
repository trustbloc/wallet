/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestClient_FormatRequest(t *testing.T) {
	t.Run("formats request", func(t *testing.T) {
		state := uuid.New().String()
		clientID := uuid.New().String()
		callbackURL := "http://test.com/callback"
		scopes := []string{"scopeA", "scopeB", "scopeC"}
		endpoint := oauth2.Endpoint{
			AuthURL:  "http://test.com/oauth2/authorize",
			TokenURL: "http://test.com/oauth2/token",
		}
		expected := (&oauth2.Config{
			ClientID:    clientID,
			Endpoint:    endpoint,
			RedirectURL: callbackURL,
			Scopes:      scopes,
		}).AuthCodeURL(state)
		result := NewClient(&Config{
			Provider:    &mockOIDCProvider{endpoint: endpoint},
			ClientID:    clientID,
			CallbackURL: callbackURL,
			Scopes:      scopes,
		}).FormatRequest(state)
		require.Equal(t, expected, result)
	})
}

func TestClient_Exchange(t *testing.T) {
	t.Run("exchanges code for token", func(t *testing.T) {
		expected := &oauth2.Token{
			AccessToken:  uuid.New().String(),
			RefreshToken: uuid.New().String(),
		}
		c := NewClient(&Config{
			Provider:     &mockOIDCProvider{},
			CallbackURL:  "http://test.com/callback",
			ClientID:     uuid.New().String(),
			ClientSecret: uuid.New().String(),
			Scopes:       []string{"scope1", "scope2"},
		})
		c.oauth2ConfigSupplier = func() oauth2Config {
			return &mockOAuth2Config{
				token: expected,
			}
		}
		result, err := c.Exchange(context.Background(), "code")
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("error if cannot exchange code for token", func(t *testing.T) {
		expected := errors.New("test")
		c := NewClient(&Config{
			Provider:     &mockOIDCProvider{},
			CallbackURL:  "http://test.com/callback",
			ClientID:     uuid.New().String(),
			ClientSecret: uuid.New().String(),
			Scopes:       []string{"scope1", "scope2"},
		})
		c.oauth2ConfigSupplier = func() oauth2Config {
			return &mockOAuth2Config{
				tokenErr: expected,
			}
		}
		_, err := c.Exchange(context.Background(), "code")
		require.True(t, errors.Is(err, expected))
	})

	t.Run("error if returned token is invalid", func(t *testing.T) {
		c := NewClient(&Config{
			Provider:     &mockOIDCProvider{},
			CallbackURL:  "http://test.com/callback",
			ClientID:     uuid.New().String(),
			ClientSecret: uuid.New().String(),
			Scopes:       []string{"scope1", "scope2"},
		})
		c.oauth2ConfigSupplier = func() oauth2Config {
			return &mockOAuth2Config{}
		}
		_, err := c.Exchange(context.Background(), "code")
		require.Error(t, err)
	})
}

func TestClient_VerifyIDToken(t *testing.T) {
	t.Run("verifies token", func(t *testing.T) {
		expected := &oidc.IDToken{
			Issuer:   "http://test.issuer.com",
			Subject:  uuid.New().String(),
			IssuedAt: time.Now(),
			Nonce:    uuid.New().String(),
		}
		c := NewClient(&Config{
			Provider: &mockOIDCProvider{
				verifier: &mockOIDCVerifier{
					token: expected,
				},
			},
			CallbackURL:  "http://test.com/callback",
			ClientID:     uuid.New().String(),
			ClientSecret: uuid.New().String(),
			Scopes:       []string{"scope1", "scope2"},
		})
		result, err := c.VerifyIDToken(context.Background(), &mockOAuthToken{extra: uuid.New().String()})
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("error if id_token is missing", func(t *testing.T) {
		c := NewClient(&Config{
			Provider: &mockOIDCProvider{
				verifier: &mockOIDCVerifier{},
			},
			CallbackURL:  "http://test.com/callback",
			ClientID:     uuid.New().String(),
			ClientSecret: uuid.New().String(),
			Scopes:       []string{"scope1", "scope2"},
		})
		_, err := c.VerifyIDToken(context.Background(), &mockOAuthToken{})
		require.Error(t, err)
	})

	t.Run("error if cannot verify id_token", func(t *testing.T) {
		expected := errors.New("test")
		c := NewClient(&Config{
			Provider: &mockOIDCProvider{
				verifier: &mockOIDCVerifier{err: expected},
			},
			CallbackURL:  "http://test.com/callback",
			ClientID:     uuid.New().String(),
			ClientSecret: uuid.New().String(),
			Scopes:       []string{"scope1", "scope2"},
		})
		_, err := c.VerifyIDToken(context.Background(), &mockOAuthToken{extra: uuid.New().String()})
		require.True(t, errors.Is(err, expected))
	})
}

type mockOIDCProvider struct {
	endpoint oauth2.Endpoint
	verifier OIDCVerifier
}

func (m *mockOIDCProvider) Endpoint() oauth2.Endpoint {
	return m.endpoint
}

func (m *mockOIDCProvider) Verifier(config *oidc.Config) OIDCVerifier {
	return m.verifier
}

type mockOAuth2Config struct {
	token    *oauth2.Token
	tokenErr error
}

func (m *mockOAuth2Config) AuthCodeURL(_ string, _ ...oauth2.AuthCodeOption) string {
	panic("implement me")
}

func (m *mockOAuth2Config) Exchange(_ context.Context, _ string, _ ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return m.token, m.tokenErr
}

type mockOAuthToken struct {
	extra interface{}
	valid bool
}

func (m *mockOAuthToken) Extra(_ string) interface{} {
	return m.extra
}

func (m *mockOAuthToken) Valid() bool {
	return m.valid
}

type mockOIDCVerifier struct {
	token *oidc.IDToken
	err   error
}

func (m *mockOIDCVerifier) Verify(_ context.Context, _ string) (*oidc.IDToken, error) {
	return m.token, m.err
}

func mockOIDCProviderService(t *testing.T) string {
	h := &testOIDCProvider{}
	srv := httptest.NewServer(h)
	h.baseURL = srv.URL

	t.Cleanup(srv.Close)

	return srv.URL
}

type oidcConfigJSON struct {
	Issuer      string   `json:"issuer"`
	AuthURL     string   `json:"authorization_endpoint"`
	TokenURL    string   `json:"token_endpoint"`
	JWKSURL     string   `json:"jwks_uri"`
	UserInfoURL string   `json:"userinfo_endpoint"`
	Algorithms  []string `json:"id_token_signing_alg_values_supported"`
}

type testOIDCProvider struct {
	baseURL string
}

func (t *testOIDCProvider) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	response, err := json.Marshal(&oidcConfigJSON{
		Issuer:      t.baseURL,
		AuthURL:     fmt.Sprintf("%s/oauth2/auth", t.baseURL),
		TokenURL:    fmt.Sprintf("%s/oauth2/token", t.baseURL),
		JWKSURL:     fmt.Sprintf("%s/oauth2/certs", t.baseURL),
		UserInfoURL: fmt.Sprintf("%s/oauth2/userinfo", t.baseURL),
		Algorithms:  []string{"RS256"},
	})
	if err != nil {
		panic(err)
	}

	_, err = w.Write(response)
	if err != nil {
		panic(err)
	}
}
