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
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
	"github.com/stretchr/testify/require"
	oidc2 "github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/cookie"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/tokens"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
	"github.com/trustbloc/edge-core/pkg/storage"
	"github.com/trustbloc/edge-core/pkg/storage/memstore"
	"github.com/trustbloc/edge-core/pkg/storage/mockstore"
	"github.com/trustbloc/edge-core/pkg/zcapld"
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

	t.Run("user already logged in", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		result := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: uuid.New().String(),
				},
			},
		}
		o.oidcLoginHandler(result, newOIDCLoginRequest())
		require.Equal(t, http.StatusMovedPermanently, result.Code)
	})
}

func TestKmsSigner_Sign(t *testing.T) {
	t.Run("failed to sign", func(t *testing.T) {
		_, err := newKMSSigner("", "", "",
			&mockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
					}, nil
				},
			}).Sign([]byte("data"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to sign from kms")
	})

	t.Run("failed to unmarshal sign resp", func(t *testing.T) {
		_, err := newKMSSigner("", "", "",
			&mockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				},
			}).Sign([]byte("data"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to unmarshal sign resp")
	})

	t.Run("failed to unmarshal sign resp", func(t *testing.T) {
		_, err := newKMSSigner("", "", "",
			&mockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"signature":"1"}`))),
					}, nil
				},
			}).Sign([]byte("data"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "illegal base64 data")
	})
}

func TestOperation_OIDCCallbackHandler(t *testing.T) { //nolint: gocritic,gocognit,gocyclo // test
	uiEndpoint := "http://test.com/dashboard"

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
				if req.URL.Path == hubAuthSecretPath {
					return &http.Response{
						StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				} else if req.URL.Path == hubAuthBootstrapDataPath {
					return &http.Response{
						StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				}

				statusCode := http.StatusCreated

				if strings.Contains(req.URL.Path, "/export") ||
					strings.Contains(req.URL.Path, "/sign") ||
					strings.Contains(req.URL.Path, "/capability") {
					statusCode = http.StatusOK
				}

				return &http.Response{
					StatusCode: statusCode, Body: ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
				}, nil
			},
		}
		o.keyEDVClient = &mockEDVClient{}
		o.userEDVClient = &mockEDVClient{}

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

	t.Run("failure to create authz key", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.Path == hubKMSCreateKeyStorePath {
					return &http.Response{
						StatusCode: http.StatusCreated, Body: ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
		}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "failed create authz key")
	})

	t.Run("failure to export authz key", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.Path == hubKMSCreateKeyStorePath || strings.Contains(req.URL.Path, "/keys") {
					return &http.Response{
						StatusCode: http.StatusCreated, Body: ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
		}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "failed export public key")
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

	t.Run("failure to post secret", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
		}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "post half secret to hub-auth")
	})

	t.Run("failure to create key data vault", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.Path == hubAuthSecretPath {
					return &http.Response{
						StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				}

				statusCode := http.StatusCreated

				if strings.Contains(req.URL.Path, "/export") {
					statusCode = http.StatusOK
				}

				return &http.Response{
					StatusCode: statusCode,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
				}, nil
			},
		}
		ops.keyEDVClient = &mockEDVClient{
			CreateErr: errors.New("vault creation error"),
		}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "vault creation error")
	})

	t.Run("failure to create ops keystore", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.Path == hubAuthSecretPath {
					return &http.Response{
						StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				}

				var request createKeystoreReq

				if req.Body != nil {
					err := json.NewDecoder(req.Body).Decode(&request)
					require.NoError(t, err)

					if request.VaultID != "" {
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
						}, nil
					}
				}

				statusCode := http.StatusCreated

				if req.Method == http.MethodGet {
					statusCode = http.StatusOK
				}

				return &http.Response{
					StatusCode: statusCode,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
				}, nil
			},
		}
		ops.keyEDVClient = &mockEDVClient{}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "create operational keystore")
	})

	t.Run("failure to create user edv vault", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.Path == hubAuthSecretPath {
					return &http.Response{
						StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				}

				statusCode := http.StatusCreated

				if strings.Contains(req.URL.Path, "/export") ||
					strings.Contains(req.URL.Path, "/sign") ||
					strings.Contains(req.URL.Path, "/capability") {
					statusCode = http.StatusOK
				}

				return &http.Response{
					StatusCode: statusCode,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
				}, nil
			},
		}
		ops.keyEDVClient = &mockEDVClient{}
		ops.userEDVClient = &mockEDVClient{CreateErr: errors.New("create error")}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "create user edv vault")
	})

	t.Run("failure to update edv capability in keystore", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.Path == hubAuthSecretPath {
					return &http.Response{
						StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				}

				statusCode := http.StatusCreated

				if strings.Contains(req.URL.Path, "/export") ||
					strings.Contains(req.URL.Path, "/sign") {
					statusCode = http.StatusOK
				}

				if strings.Contains(req.URL.Path, "/capability") {
					statusCode = http.StatusInternalServerError
				}

				return &http.Response{
					StatusCode: statusCode,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
				}, nil
			},
		}
		ops.keyEDVClient = &mockEDVClient{}
		ops.userEDVClient = &mockEDVClient{CreateErr: errors.New("create error")}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "failed to update edv capability keystore")
	})

	t.Run("failure to post bootstrap data", func(t *testing.T) {
		state := uuid.New().String()
		ops := setupOnboardingTest(t, state)
		ops.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.Path == hubAuthSecretPath {
					return &http.Response{
						StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				}

				statusCode := http.StatusCreated

				if strings.Contains(req.URL.Path, "/export") ||
					strings.Contains(req.URL.Path, "/sign") ||
					strings.Contains(req.URL.Path, "/capability") {
					statusCode = http.StatusOK
				}

				return &http.Response{
					StatusCode: statusCode,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
				}, nil
			},
		}
		ops.keyEDVClient = &mockEDVClient{}
		ops.userEDVClient = &mockEDVClient{}

		w := httptest.NewRecorder()
		ops.oidcCallbackHandler(w, newOIDCCallbackRequest(uuid.New().String(), state))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "update user bootstrap data")
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

		originalZcap, err := zcapld.NewCapability(&zcapld.Signer{
			SignatureSuite:     ed25519signature2018.New(suite.WithSigner(&mockSigner{})),
			SuiteType:          ed25519signature2018.SignatureType,
			VerificationMethod: "test:123",
		}, zcapld.WithParent(uuid.New().URN()))
		require.NoError(t, err)

		originalZcapBytes, err := json.Marshal(originalZcap)
		require.NoError(t, err)

		d := &BootstrapData{
			AuthzKeyStoreURL:  "http://localhost/authz/kms/" + uuid.New().String(),
			OpsKeyStoreURL:    "http://localhost/ops/kms/" + uuid.New().String(),
			UserEDVVaultURL:   "http://localhost/user/vault/" + uuid.New().String(),
			OpsEDVVaultURL:    "http://localhost/ops/vault/" + uuid.New().String(),
			UserEDVCapability: string(originalZcapBytes),
		}
		o.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				data := userBootstrapData{
					Data: d,
				}

				resp, respErr := json.Marshal(data)
				require.NoError(t, respErr)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(resp)),
				}, nil
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

		b, err := json.Marshal(resultData["bootstrap"])
		require.NoError(t, err)

		respData := BootstrapData{}

		err = json.Unmarshal(b, &respData)
		require.NoError(t, err)

		require.Equal(t, d.AuthzKeyStoreURL, respData.AuthzKeyStoreURL)
		require.Equal(t, d.OpsEDVVaultURL, respData.OpsEDVVaultURL)
		require.Equal(t, d.OpsKeyStoreURL, respData.OpsKeyStoreURL)
		require.Equal(t, d.UserEDVVaultURL, respData.UserEDVVaultURL)

		zCapResp := &zcapld.Capability{}

		err = json.Unmarshal([]byte(respData.UserEDVCapability), zCapResp)
		require.NoError(t, err)

		require.Equal(t, originalZcap.Controller, zCapResp.Controller)
		require.Equal(t, originalZcap.ID, zCapResp.ID)
		require.Equal(t, originalZcap.Parent, zCapResp.Parent)
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

	t.Run("err internalserver error if cannot fetch temporary bootstrap data", func(t *testing.T) {
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
		o.httpClient = &mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
				}, nil
			},
		}

		result := httptest.NewRecorder()
		o.userProfileHandler(result, newUserProfileRequest())
		require.Equal(t, http.StatusInternalServerError, result.Code)
		require.Contains(t, result.Body.String(), "failed to fetch bootstrap data")
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
			KeyEDVURL:   "",
			OpsKMSURL:   "",
		},
		UserEDVURL: "http://example.com",
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
	config.WalletDashboard = "http://test.com/dashboard"
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

type mockEDVClient struct {
	CreateErr error
}

func (m *mockEDVClient) CreateDataVault(_ *models.DataVaultConfiguration,
	_ ...client.ReqOption) (string, []byte, error) {
	if m.CreateErr != nil {
		return "", nil, m.CreateErr
	}

	c, err := zcapld.NewCapability(&zcapld.Signer{
		SignatureSuite:     ed25519signature2018.New(suite.WithSigner(&mockSigner{})),
		SuiteType:          ed25519signature2018.SignatureType,
		VerificationMethod: "test:123",
	}, zcapld.WithParent(uuid.New().URN()))
	if err != nil {
		return "", nil, err
	}

	b, err := json.Marshal(c)
	if err != nil {
		return "", nil, err
	}

	return "http://edv.example.com" + uuid.New().String(), b, nil
}

type mockSigner struct {
}

func (m *mockSigner) Sign(data []byte) ([]byte, error) {
	return nil, nil
}
