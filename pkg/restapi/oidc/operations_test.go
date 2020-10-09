/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc // nolint:testpackage // changing to different package requires exposing internal REST handlers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	oidc2 "github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/cookie"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/tokens"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
	"github.com/trustbloc/edge-core/pkg/storage"
	"github.com/trustbloc/edge-core/pkg/storage/memstore"
	"github.com/trustbloc/edge-core/pkg/storage/mockstore"
	"golang.org/x/oauth2"
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
			FailNameSpace: user.StoreName,
		}
		_, err := New(config)
		require.Error(t, err)
	})

	t.Run("error if cannot open token store", func(t *testing.T) {
		config := config(t)
		config.Storage.Storage = &mockstore.Provider{
			FailNameSpace: tokens.StoreName,
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

	t.Run("internal server error if cannot fetch session cookie", func(t *testing.T) {
		config := config(t)
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			OpenErr: errors.New("test"),
		}
		w := httptest.NewRecorder()
		o.oidcLoginHandler(w, newOIDCLoginRequest())
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("internal server error if cannot save to cookie store", func(t *testing.T) {
		config := config(t)
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				SaveErr: errors.New("test"),
			},
		}
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
		config.OIDCClient = &oidc2.MockClient{
			OAuthToken: &oauth2.Token{
				AccessToken:  uuid.New().String(),
				RefreshToken: uuid.New().String(),
				TokenType:    "Bearer",
			},
			IDToken: &oidc2.MockIDToken{
				ClaimsFunc: func(i interface{}) error {
					user, ok := i.(*user.User)
					require.True(t, ok)
					user.Sub = uuid.New().String()

					return nil
				},
			},
		}

		o, err := New(config)
		require.NoError(t, err)

		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					stateCookieName: state,
				},
			},
		}

		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest(code, state))

		require.Equal(t, http.StatusFound, w.Code)
		require.Equal(t, uiEndpoint, w.Header().Get("Location"))
	})

	t.Run("error internal server error if cannot fetch the user's session", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			OpenErr: errors.New("test"),
		}
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", ""))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error bad request if state cookie is not present", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", ""))
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error bad request if state query param is missing", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					stateCookieName: "123",
				},
			},
		}
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", ""))
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error bad request if state query param does not match state cookie", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					stateCookieName: "123",
				},
			},
		}
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", "456"))
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error bad request if code query param is missing", func(t *testing.T) {
		state := uuid.New().String()
		o, err := New(config(t))
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					stateCookieName: state,
				},
			},
		}
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("", state))
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error internal server error if cannot fetch session cookie", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			OpenErr: errors.New("test"),
		}
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error internal server error if cannot persist session cookies", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					stateCookieName: state,
				},
				SaveErr: errors.New("test"),
			},
		}
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error bad gateway if cannot exchange code for token", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		config.OIDCClient = &oidc2.MockClient{
			OAuthErr: errors.New("test"),
		}
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					stateCookieName: state,
				},
			},
		}
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusBadGateway, w.Code)
	})

	t.Run("error bad gateway if cannot verify id_token", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		config.OIDCClient = &oidc2.MockClient{
			IDTokenErr: errors.New("test"),
		}
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					stateCookieName: state,
				},
			},
		}
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusBadGateway, w.Code)
	})

	t.Run("error internal server error if cannot parse id_token", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		config.OIDCClient = &oidc2.MockClient{
			IDToken: &oidc2.MockIDToken{
				ClaimsErr: errors.New("test"),
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					stateCookieName: state,
				},
			},
		}
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error internal server error if cannot query user store", func(t *testing.T) {
		userSub := uuid.New().String()
		state := uuid.New().String()
		config := config(t)
		config.Storage.Storage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store: map[string][]byte{
					userSub: []byte(userSub),
				},
				ErrGet: errors.New("test"),
			},
		}
		config.OIDCClient = &oidc2.MockClient{
			IDToken: &oidc2.MockIDToken{
				ClaimsFunc: func(i interface{}) error {
					user, ok := i.(*user.User)
					require.True(t, ok)
					user.Sub = userSub

					return nil
				},
			},
			OAuthToken: &oauth2.Token{
				AccessToken:  uuid.New().String(),
				RefreshToken: uuid.New().String(),
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					stateCookieName: state,
				},
			},
		}
		w := httptest.NewRecorder()
		o.oidcCallbackHandler(w, newOIDCCallbackRequest("code", state))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error internal server error if cannot save to user store", func(t *testing.T) {
		state := uuid.New().String()
		config := config(t)
		config.Storage.Storage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store:  make(map[string][]byte),
				ErrPut: errors.New("test"),
			},
		}
		config.OIDCClient = &oidc2.MockClient{
			IDToken: &oidc2.MockIDToken{
				ClaimsFunc: func(i interface{}) error {
					user, ok := i.(*user.User)
					require.True(t, ok)
					user.Sub = uuid.New().String()

					return nil
				},
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					stateCookieName: state,
				},
			},
		}
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
		OIDCClient: &oidc2.MockClient{},
		Storage: &StorageConfig{
			Storage:          memstore.NewProvider(),
			TransientStorage: memstore.NewProvider(),
		},
		Keys: &KeyConfig{
			Auth: key(t),
			Enc:  key(t),
		},
	}
}

func key(t *testing.T) []byte {
	t.Helper()

	key := make([]byte, 32)

	n, err := rand.Reader.Read(key)
	require.NoError(t, err)
	require.Equal(t, 32, n)

	return key
}
