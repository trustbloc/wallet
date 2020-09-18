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
	"github.com/trustbloc/edge-core/pkg/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/trustbloc/edge-core/pkg/storage/mockstore"

	"github.com/coreos/go-oidc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edge-core/pkg/storage/memstore"
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
}

func TestOperation_GetRESTHandlers(t *testing.T) {
	o, err := New(config(t))
	require.NoError(t, err)

	require.NotEmpty(t, o.GetRESTHandlers())
}

func TestOperation_OIDCRequestHandler(t *testing.T) {
	t.Run("redirects to OIDC provider", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		w := httptest.NewRecorder()
		o.oidcLoginHandler(w, newOIDCHTTPRequest())
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
		o.oidcLoginHandler(w, newOIDCHTTPRequest())
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func newOIDCHTTPRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/oidc/login", nil)
}

func config(t *testing.T) *Config {
	oidcProvider, err := oidc.NewProvider(context.Background(), newTestOIDCProvider(t))
	require.NoError(t, err)

	return &Config{
		OIDC: &OIDCConfig{
			Provider:     &oidcProviderImpl{op: oidcProvider},
			ClientID:     uuid.New().String(),
			ClientSecret: uuid.New().String(),
			Scopes:       []string{oidc.ScopeOpenID, "test"},
			CallbackURL:  "http://test.com/callback",
		},
		Storage: &StorageConfig{
			Storage:          memstore.NewProvider(),
			TransientStorage: memstore.NewProvider(),
		},
	}
}

func newTestOIDCProvider(t *testing.T) string {
	h := &testOIDCProvider{}
	srv := httptest.NewServer(h)
	h.baseURL = srv.URL
	t.Cleanup(srv.Close)

	return srv.URL
}

type testOIDCProvider struct {
	baseURL string
}

func (t *testOIDCProvider) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	response, err := json.Marshal(map[string]interface{}{
		"issuer":                                t.baseURL,
		"authorization_endpoint":                fmt.Sprintf("%s/oauth2/auth", t.baseURL),
		"token_endpoint":                        fmt.Sprintf("%s/oauth2/token", t.baseURL),
		"jwks_uri":                              fmt.Sprintf("%s/oauth2/certs", t.baseURL),
		"userinfo_endpoint":                     fmt.Sprintf("%s/oauth2/userinfo", t.baseURL),
		"id_token_signing_alg_values_supported": []string{"RS256"},
	})
	if err != nil {
		panic(err)
	}

	_, err = w.Write(response)
	if err != nil {
		panic(err)
	}
}
