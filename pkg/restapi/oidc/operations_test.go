/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	oidc2 "github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"
	"github.com/trustbloc/edge-core/pkg/storage"
	"github.com/trustbloc/edge-core/pkg/storage/memstore"
	"github.com/trustbloc/edge-core/pkg/storage/mockstore"
	"golang.org/x/oauth2"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("returns an instance", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		require.NotNil(t, o)
	})

	t.Run("can init if transient store already exists", func(t *testing.T) {
		config := config(t)
		config.Storage.TransientStorage = &mockstore.Provider{
			ErrCreateStore: storage.ErrDuplicateStore,
		}
		o, err := New(config)
		require.NoError(t, err)
		require.NotNil(t, o)
	})

	t.Run("error if cannot create transient store", func(t *testing.T) {
		expected := errors.New("test")
		config := config(t)
		config.Storage.TransientStorage = &mockstore.Provider{
			ErrCreateStore: expected,
		}
		_, err := New(config)
		require.Error(t, err)
		require.True(t, errors.Is(err, expected))
	})

	t.Run("error if cannot open transient store", func(t *testing.T) {
		expected := errors.New("test")
		config := config(t)
		config.Storage.TransientStorage = &mockstore.Provider{
			ErrOpenStoreHandle: expected,
		}
		_, err := New(config)
		require.Error(t, err)
		require.True(t, errors.Is(err, expected))
	})

	t.Run("error if cannot open user store", func(t *testing.T) {
		config := config(t)
		config.Storage.Storage = &mockstore.Provider{
			FailNameSpace:      userStoreName,
		}
		_, err := New(config)
		require.Error(t, err)
	})

	t.Run("error if cannot open token store", func(t *testing.T) {
		config := config(t)
		config.Storage.Storage = &mockstore.Provider{
			FailNameSpace:      tokenStoreName,
		}
		_, err := New(config)
		require.Error(t, err)
	})
}

func TestOperation_GetRESTHandlers(t *testing.T) {
	o, err := New(config(t))
	require.NoError(t, err)

	require.NotEmpty(t, o.GetRESTHandlers())
}

func TestOperation_OIDCLoginHandler(t *testing.T) {
	t.Run("redirects to OIDC provider", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcLoginHandler(w, newOIDCLoginRequest())
		require.Equal(t, http.StatusFound, w.Code)
		require.NotEmpty(t, w.Header().Get("Location"))
	})

	t.Run("internal server error if cannot save to transient store", func(t *testing.T) {
		config := config(t)
		config.Storage.TransientStorage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store:  make(map[string][]byte),
				ErrPut: errors.New("test"),
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcLoginHandler(w, newOIDCLoginRequest())
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestOperation_OIDCCallbackHandler(t *testing.T) {
	t.Run("fetches OIDC tokens and redirects to the UI", func(t *testing.T) {
		uiEndpoint := "http://test.com/wallet/"
		code := uuid.New().String()
		state := uuid.New().String()

		config := config(t)
		config.UIEndpoint = uiEndpoint
		config.Storage.TransientStorage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store: map[string][]byte{
					state: []byte(state),
				},
			},
		}
		config.OIDCClient = &mockOIDCClient{
			oauthToken: &oauth2.Token{
				AccessToken:  uuid.New().String(),
				RefreshToken: uuid.New().String(),
				TokenType:    "Bearer",
			},
			idToken: &mockIDToken{
				claimsFunc: func(i interface{}) error {
					user, ok := i.(*endUser)
					require.True(t, ok)
					user.Sub = uuid.New().String()
					return nil
				},
			},
		}

		o, err := New(config)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest(code, state))

		require.Equal(t, http.StatusFound, w.Code)
		require.Equal(t, uiEndpoint, w.Header().Get("Location"))
	})

	t.Run("error bad request if state query param is missing", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", ""))
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error bad request if state query param is invalid", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", "INVALID"))
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error bad request if code query param is missing", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("", "state"))
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error internal server error if cannot query transient store", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		config.Storage.TransientStorage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store:  map[string][]byte{
					state: []byte(state),
				},
				ErrGet: errors.New("test"),
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error internal server error if cannot delete transient state", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		config.Storage.TransientStorage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store:  map[string][]byte{
					state: []byte(state),
				},
				ErrDelete: errors.New("test"),
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error bad gateway if cannot exchange code for token", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		config.Storage.TransientStorage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store:  map[string][]byte{
					state: []byte(state),
				},
			},
		}
		config.OIDCClient = &mockOIDCClient{
			oauthErr: errors.New("test"),
		}
		o, err := New(config)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusBadGateway, w.Code)
	})

	t.Run("error bad gateway if cannot verify id_token", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		config.Storage.TransientStorage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store:  map[string][]byte{
					state: []byte(state),
				},
			},
		}
		config.OIDCClient = &mockOIDCClient{
			idTokenErr: errors.New("test"),
		}
		o, err := New(config)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusBadGateway, w.Code)
	})

	t.Run("error internal server error if cannot parse id_token", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		config.Storage.TransientStorage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store:  map[string][]byte{
					state: []byte(state),
				},
			},
		}
		config.OIDCClient = &mockOIDCClient{
			idToken: &mockIDToken{
				claimsErr: errors.New("test"),
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error internal server error if cannot save to user store", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		config.Storage.TransientStorage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store:  map[string][]byte{
					state: []byte(state),
				},
			},
		}
		config.Storage.Storage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store:  make(map[string][]byte),
				ErrPut: errors.New("test"),
			},
		}
		config.OIDCClient = &mockOIDCClient{
			idToken: &mockIDToken{
				claimsFunc: func(i interface{}) error {
					user, ok := i.(*endUser)
					require.True(t, ok)
					user.Sub = uuid.New().String()
					return nil
				},
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func newOIDCLoginRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/oidc/login", nil)
}

func newOIDCCallbackRequest(code, state string) *http.Request {
	return httptest.NewRequest(http.MethodGet, fmt.Sprintf("/oidc/callback?code=%s&state=%s", code, state), nil)
}

func config(t *testing.T) *Config {
	t.Helper()

	return &Config{
		OIDCClient: &mockOIDCClient{},
		Storage: &StorageConfig{
			Storage:          memstore.NewProvider(),
			TransientStorage: memstore.NewProvider(),
		},
	}
}

type mockOIDCClient struct {
	authRequest string
	oauthToken  *oauth2.Token
	oauthErr    error
	idToken     oidc2.IDToken
	idTokenErr  error
}

func (m *mockOIDCClient) FormatRequest(_ string) string {
	return m.authRequest
}

func (m *mockOIDCClient) Exchange(_ context.Context, _ string) (*oauth2.Token, error) {
	return m.oauthToken, m.oauthErr
}

func (m *mockOIDCClient) VerifyIDToken(_ context.Context, _ oidc2.OAuth2Token) (oidc2.IDToken, error) {
	return m.idToken, m.idTokenErr
}

type mockIDToken struct {
	claimsErr  error
	claimsFunc func(interface{}) error
}

func (m *mockIDToken) Claims(i interface{}) error {
	if m.claimsFunc != nil {
		return m.claimsFunc(i)
	}

	return m.claimsErr
}
