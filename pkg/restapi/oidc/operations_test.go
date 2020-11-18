/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc // nolint:testpackage // changing to different package requires exposing internal REST handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	"github.com/trustbloc/edv/pkg/client"
	"github.com/trustbloc/edv/pkg/restapi/models"
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
	uiEndpoint := "http://test.com/wallet/"

	t.Run("fetches OIDC tokens and redirects to the UI", func(t *testing.T) {
		code := uuid.New().String()
		state := uuid.New().String()

		config := config(t)
		config.WalletDashboard = uiEndpoint
		config.OIDCClient = &oidc2.MockClient{
			OAuthToken: &oauth2.Token{
				AccessToken:  uuid.New().String(),
				RefreshToken: uuid.New().String(),
				TokenType:    "Bearer",
			},
			IDToken: &oidc2.MockClaimer{
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

		o.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusCreated, Body: ioutil.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
		}
		o.keySDSClient = &mockSDSClient{}

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
			IDToken: &oidc2.MockClaimer{
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
			IDToken: &oidc2.MockClaimer{
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
			IDToken: &oidc2.MockClaimer{
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

	t.Run("failure to create authz keystore", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
		}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "create authz keystore")
	})

	t.Run("failure to split secret", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.secretSplitter = &mockSplitter{SplitErr: errors.New("secret split error")}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "split user secret key")
	})

	t.Run("failure to create key data vault", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
		}
		ops.keySDSClient = &mockSDSClient{
			CreateErr: errors.New("vault creation error"),
		}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "create key sds vault")
	})

	t.Run("failure to create ops keystore", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				var request createKeystoreReq
				err := json.NewDecoder(req.Body).Decode(&request)
				require.NoError(t, err)

				if request.OperationalVaultID != "" {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
		}
		ops.keySDSClient = &mockSDSClient{}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "create operational keystore")
	})
}

func TestOperation_UserProfileHandler(t *testing.T) {
	t.Run("returns the user profile", func(t *testing.T) {
		sub := uuid.New().String()
		config := config(t)
		config.Storage.Storage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store: map[string][]byte{
					sub: marshal(t, &tokens.UserTokens{}),
				},
			},
		}
		config.OIDCClient = &oidc2.MockClient{
			UserInfoVal: &oidc2.MockClaimer{
				ClaimsFunc: func(v interface{}) error {
					m, ok := v.(*map[string]interface{})
					require.True(t, ok)
					(*m)["sub"] = sub

					return nil
				},
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: sub,
				},
			},
		}
		result := httptest.NewRecorder()
		o.userProfileHandler(result, newUserProfileRequest())
		require.Equal(t, http.StatusOK, result.Code)
		resultData := make(map[string]interface{})
		err = json.NewDecoder(result.Body).Decode(&resultData)
		require.NoError(t, err)
		require.Contains(t, resultData, "sub")
		require.Equal(t, sub, resultData["sub"])
	})

	t.Run("err badrequest if cannot open cookies", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			OpenErr: errors.New("test"),
		}
		result := httptest.NewRecorder()
		o.userProfileHandler(result, newUserProfileRequest())
		require.Equal(t, http.StatusBadRequest, result.Code)
		require.Contains(t, result.Body.String(), "cannot open cookies")
	})

	t.Run("err forbidden if user cookie is not set", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		result := httptest.NewRecorder()
		o.userProfileHandler(result, newUserProfileRequest())
		require.Equal(t, http.StatusForbidden, result.Code)
		require.Contains(t, result.Body.String(), "not logged in")
	})

	t.Run("err internalservererror if cookie is not a string (should not happen)", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: struct{}{},
				},
			},
		}
		result := httptest.NewRecorder()
		o.userProfileHandler(result, newUserProfileRequest())
		require.Equal(t, http.StatusInternalServerError, result.Code)
		require.Contains(t, result.Body.String(), "invalid user sub cookie format")
	})

	t.Run("err internal server error if cannot fetch user tokens from storage", func(t *testing.T) {
		config := config(t)
		config.Storage.Storage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				ErrGet: errors.New("test"),
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: uuid.New().String(),
				},
			},
		}
		result := httptest.NewRecorder()
		o.userProfileHandler(result, newUserProfileRequest())
		require.Equal(t, http.StatusInternalServerError, result.Code)
		require.Contains(t, result.Body.String(), "failed to fetch user tokens from store")
	})

	t.Run("err badgateway error if cannot fetch userinfo from oidc provider", func(t *testing.T) {
		sub := uuid.New().String()
		config := config(t)
		config.OIDCClient = &oidc2.MockClient{UserInfoErr: errors.New("test")}
		config.Storage.Storage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store: map[string][]byte{
					sub: marshal(t, &tokens.UserTokens{}),
				},
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: sub,
				},
			},
		}
		result := httptest.NewRecorder()
		o.userProfileHandler(result, newUserProfileRequest())
		require.Equal(t, http.StatusBadGateway, result.Code)
		require.Contains(t, result.Body.String(), "failed to fetch user info")
	})

	t.Run("err internalservererror if cannot extract claims from userinfo", func(t *testing.T) {
		sub := uuid.New().String()
		config := config(t)
		config.OIDCClient = &oidc2.MockClient{UserInfoVal: &oidc2.MockClaimer{
			ClaimsErr: errors.New("test"),
		}}
		config.Storage.Storage = &mockstore.Provider{
			Store: &mockstore.MockStore{
				Store: map[string][]byte{
					sub: marshal(t, &tokens.UserTokens{}),
				},
			},
		}
		o, err := New(config)
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: sub,
				},
			},
		}
		result := httptest.NewRecorder()
		o.userProfileHandler(result, newUserProfileRequest())
		require.Equal(t, http.StatusInternalServerError, result.Code)
		require.Contains(t, result.Body.String(), "failed to extract claims from user info")
	})
}

func TestOperation_UserLogoutHandler(t *testing.T) {
	t.Run("logs out user", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: uuid.New().String(),
				},
			},
		}
		result := httptest.NewRecorder()
		o.userLogoutHandler(result, newUserLogoutRequest())
		require.Equal(t, http.StatusOK, result.Code)
	})

	t.Run("err badrequest if cannot open cookies", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			OpenErr: errors.New("test"),
		}
		result := httptest.NewRecorder()
		o.userLogoutHandler(result, newUserLogoutRequest())
		require.Equal(t, http.StatusBadRequest, result.Code)
		require.Contains(t, result.Body.String(), "cannot open cookies")
	})

	t.Run("no-op if user sub cookie is not found", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		result := httptest.NewRecorder()
		o.userLogoutHandler(result, newUserLogoutRequest())
		require.Equal(t, http.StatusOK, result.Code)
	})

	t.Run("err internal server error if cannot delete cookie", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: uuid.New().String(),
				},
				SaveErr: errors.New("test"),
			},
		}
		result := httptest.NewRecorder()
		o.userLogoutHandler(result, newUserLogoutRequest())
		require.Equal(t, http.StatusInternalServerError, result.Code)
		require.Contains(t, result.Body.String(), "failed to delete user sub cookie")
	})
}

func newOIDCLoginRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/oidc/login", nil)
}

func newOIDCCallbackRequest(code, state string) *http.Request {
	return httptest.NewRequest(http.MethodGet, fmt.Sprintf("/oidc/callback?code=%s&state=%s", code, state), nil)
}

func newUserProfileRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/oidc/userinfo", nil)
}

func newUserLogoutRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/oidc/logout", nil)
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
		KeyServer: &KeyServerConfig{
			AuthzKMSURL: "",
			KeySDSURL:   "",
			OpsKMSURL:   "",
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

func marshal(t *testing.T, v interface{}) []byte {
	t.Helper()

	bits, err := json.Marshal(v)
	require.NoError(t, err)

	return bits
}

func setupOnboardingTest(t *testing.T, state string) *Operation {
	config := config(t)
	config.WalletDashboard = "http://test.com/wallet/"
	config.OIDCClient = &oidc2.MockClient{
		OAuthToken: &oauth2.Token{
			AccessToken:  uuid.New().String(),
			RefreshToken: uuid.New().String(),
			TokenType:    "Bearer",
		},
		IDToken: &oidc2.MockClaimer{
			ClaimsFunc: func(i interface{}) error {
				user, ok := i.(*user.User)
				require.True(t, ok)
				user.Sub = uuid.New().String()

				return nil
			},
		},
	}

	ops, err := New(config)
	require.NoError(t, err)

	ops.store.cookies = &cookie.MockStore{
		Jar: &cookie.MockJar{
			Cookies: map[interface{}]interface{}{
				stateCookieName: state,
			},
		},
	}

	return ops
}

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

type mockSplitter struct {
	SplitErr   error
	CombineErr error
}

func (m *mockSplitter) Split(secret []byte, numParts, threshold int) ([][]byte, error) {
	return nil, m.SplitErr
}

func (m *mockSplitter) Combine(secretParts [][]byte) ([]byte, error) {
	return nil, m.CombineErr
}

type mockSDSClient struct {
	CreateErr error
}

func (m *mockSDSClient) CreateDataVault(_ *models.DataVaultConfiguration, _ ...client.EDVOption) (string, error) {
	if m.CreateErr != nil {
		return "", m.CreateErr
	}

	return "http://sds.example.com" + uuid.New().String(), nil
}
