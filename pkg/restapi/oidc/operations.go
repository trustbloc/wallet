/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/google/uuid"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/edge-core/pkg/storage"
	"golang.org/x/oauth2"

	"github.com/trustbloc/edge-agent/pkg/restapi/common"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"
)

// Endpoints.
const (
	oidcLoginPath    = "/login"
	oidcCallbackPath = "/callback"
)

// Stores.
const (
	transientStoreName = "edgeagent_trx"
	userStoreName      = "edgeagent_users"
	tokenStoreName     = "edgeagent_tks"
	cookieStoreName    = "edgeagent_wallet"
	stateCookieName    = "oauth2_state"
)

var logger = log.New("hub-auth/oidc")

// Config holds all configuration for an Operation.
type Config struct {
	OIDCClient OIDCClient
	Storage    *StorageConfig
	UIEndpoint string
	TLSConfig  *tls.Config
	Keys       *KeyConfig
}

type KeyConfig struct {
	Auth []byte
	Enc  []byte
}

// StorageConfig holds storage config.
type StorageConfig struct {
	Storage          storage.Provider
	TransientStorage storage.Provider
}

type store struct {
	users     storage.Store
	tokens    storage.Store
	transient storage.Store
	cookies   CookieStore
}

// Operation implements OIDC operations.
type Operation struct {
	store      *store
	oidcClient OIDCClient
	uiEndpoint string
	tlsConfig  *tls.Config
}

// New returns a new Operation.
func New(config *Config) (*Operation, error) {
	// TODO make session cookies max age configurable: https://github.com/trustbloc/edge-agent/issues/388
	cookieStore := sessions.NewCookieStore(config.Keys.Auth, config.Keys.Enc)
	cookieStore.MaxAge(900) // 15 mins

	op := &Operation{
		oidcClient: config.OIDCClient,
		store: &store{
			cookies: &cookieStoreAdapter{
				cookieStore: cookieStore,
			},
		},
		uiEndpoint: config.UIEndpoint,
		tlsConfig:  config.TLSConfig,
	}

	var err error

	op.store.transient, err = openStore(config.Storage.TransientStorage, transientStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to open transient store: %w", err)
	}

	op.store.users, err = openStore(config.Storage.Storage, userStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to open users store: %w", err)
	}

	op.store.tokens, err = openStore(config.Storage.Storage, tokenStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to open tokens store: %w", err)
	}

	return op, nil
}

// GetRESTHandlers get all controller API handler available for this service.
func (o *Operation) GetRESTHandlers() []common.Handler {
	return []common.Handler{
		common.NewHTTPHandler(oidcLoginPath, http.MethodGet, o.oidcLoginHandler),
		common.NewHTTPHandler(oidcCallbackPath, http.MethodGet, o.oidcCallbackHandler),
	}
}

func (o *Operation) oidcLoginHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("handling login request: %s", r.URL.String())

	session, err := o.store.cookies.Get(r, cookieStoreName)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to create or decode session cookie: %s", err.Error())

		return
	}

	state := uuid.New().String()
	session.Set(stateCookieName, state)
	redirectURL := o.oidcClient.FormatRequest(state)

	err = session.Save(r, w)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save session cookie: %s", err.Error())

		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
	logger.Debugf("redirected to login url: %s", redirectURL)
}

// TODO encrypt data before storing: https://github.com/trustbloc/edge-agent/issues/380
func (o *Operation) oidcCallbackHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("handling oidc callback: %s", r.URL.String())

	oauthToken, oidcToken, canProceed := o.fetchTokens(w, r)
	if !canProceed {
		return
	}

	user := &endUser{}

	err := user.parse(oidcToken)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to parse id_token: %s", err.Error())

		return
	}

	_, err = newPersistedData(o.store.users).getEndUser(user.Sub)
	if err != nil && !errors.Is(err, storage.ErrValueNotFound) {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to query user data: %s", err.Error())

		return
	}

	if errors.Is(err, storage.ErrValueNotFound) {
		err = newPersistedData(o.store.users).put(user.Sub, user)
		if err != nil {
			common.WriteErrorResponsef(w, logger,
				http.StatusInternalServerError, "failed to persist user data: %s", err.Error())

			return
		}
	}

	err = newPersistedData(o.store.tokens).put(user.Sub, &endUserTokens{
		ID:      user.Sub,
		Access:  oauthToken.AccessToken,
		Refresh: oauthToken.RefreshToken,
	})
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to persist user tokens: %s", err.Error())

		return
	}

	http.Redirect(w, r, o.uiEndpoint, http.StatusFound)
	logger.Debugf("redirected user to: %s", o.uiEndpoint)
}

func (o *Operation) fetchTokens(
	w http.ResponseWriter, r *http.Request) (oauthToken *oauth2.Token, oidcToken oidc.IDToken, proceed bool) {
	session, err := o.store.cookies.Get(r, cookieStoreName)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to create or decode session cookie: %s", err.Error())

		return nil, nil, false
	}

	stateCookie, found := session.Get(stateCookieName)
	if !found {
		common.WriteErrorResponsef(w, logger, http.StatusBadRequest, "missing state session cookie")

		return nil, nil, false
	}

	state := r.URL.Query().Get("state")
	if state == "" {
		common.WriteErrorResponsef(w, logger, http.StatusBadRequest, "missing state parameter")

		return nil, nil, false
	}

	if state != stateCookie {
		common.WriteErrorResponsef(w, logger, http.StatusBadRequest, "invalid state parameter")

		return nil, nil, false
	}

	session.Delete(stateCookieName)

	code := r.URL.Query().Get("code")
	if code == "" {
		common.WriteErrorResponsef(w, logger, http.StatusBadRequest, "missing code parameter")

		return nil, nil, false
	}

	oauthToken, err = o.oidcClient.Exchange(
		context.WithValue(
			r.Context(),
			oauth2.HTTPClient,
			&http.Client{Transport: &http.Transport{TLSClientConfig: o.tlsConfig}},
		),
		code,
	)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadGateway, "unable to exchange code for token: %s", err.Error())

		return nil, nil, false
	}

	oidcToken, err = o.oidcClient.VerifyIDToken(r.Context(), oauthToken)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadGateway, "cannot verify id_token: %s", err.Error())

		return nil, nil, false
	}

	err = session.Save(r, w)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save session cookies: %s", err.Error())

		return nil, nil, false
	}

	return oauthToken, oidcToken, true
}
