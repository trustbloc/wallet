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

	"github.com/google/uuid"
	"github.com/trustbloc/edge-agent/pkg/restapi/common"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/cookie"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/tokens"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/edge-core/pkg/storage"
	"golang.org/x/oauth2"
)

// Endpoints.
const (
	oidcLoginPath    = "/login"
	oidcCallbackPath = "/callback"
)

// Stores.
const (
	transientStoreName = "edgeagent_oidc_trx"
	stateCookieName    = "oauth2_state"
	deviceCookieName   = "device_user"
)

var logger = log.New("hub-auth/oidc")

// Config holds all configuration for an Operation.
type Config struct {
	OIDCClient oidc.Client
	Storage    *StorageConfig
	UIEndpoint string
	TLSConfig  *tls.Config
	Keys       *KeyConfig
}

// KeyConfig holds configuration for cryptographic keys.
type KeyConfig struct {
	Auth []byte
	Enc  []byte
}

// StorageConfig holds storage config.
type StorageConfig struct {
	Storage          storage.Provider
	TransientStorage storage.Provider
}

type stores struct {
	users     *user.Store
	tokens    *tokens.Store
	transient storage.Store
	cookies   cookie.Store
}

// Operation implements OIDC operations.
type Operation struct {
	store      *stores
	oidcClient oidc.Client
	uiEndpoint string
	tlsConfig  *tls.Config
}

// New returns a new Operation.
func New(config *Config) (*Operation, error) {
	op := &Operation{
		oidcClient: config.OIDCClient,
		store: &stores{
			cookies: cookie.NewStore(config.Keys.Auth, config.Keys.Enc),
		},
		uiEndpoint: config.UIEndpoint,
		tlsConfig:  config.TLSConfig,
	}

	var err error

	op.store.transient, err = store.Open(config.Storage.TransientStorage, transientStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to open transient store: %w", err)
	}

	op.store.users, err = user.NewStore(config.Storage.Storage)
	if err != nil {
		return nil, fmt.Errorf("failed to open users store: %w", err)
	}

	op.store.tokens, err = tokens.NewStore(config.Storage.Storage)
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

	session, err := o.store.cookies.Open(r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to read user session cookie: %s", err.Error())

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

	usr, err := user.ParseIDToken(oidcToken)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to parse id_token: %s", err.Error())

		return
	}

	_, err = o.store.users.Get(usr.Sub)
	if err != nil && !errors.Is(err, storage.ErrValueNotFound) {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to query user data: %s", err.Error())

		return
	}

	if errors.Is(err, storage.ErrValueNotFound) {
		err = o.store.users.Save(usr)
		if err != nil {
			common.WriteErrorResponsef(w, logger,
				http.StatusInternalServerError, "failed to persist user data: %s", err.Error())

			return
		}
	}

	err = o.store.tokens.Save(&tokens.UserTokens{
		UserSub: usr.Sub,
		Access:  oauthToken.AccessToken,
		Refresh: oauthToken.RefreshToken,
	})
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to persist user tokens: %s", err.Error())

		return
	}

	session, err := o.store.cookies.Open(r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to create or decode user sub session cookie: %s", err.Error())

		return
	}

	session.Set(deviceCookieName, usr.Sub)

	http.Redirect(w, r, o.uiEndpoint, http.StatusFound)
	logger.Debugf("redirected user to: %s", o.uiEndpoint)
}

func (o *Operation) fetchTokens(
	w http.ResponseWriter, r *http.Request) (oauthToken *oauth2.Token, oidcToken oidc.IDToken, valid bool) {
	session, valid := o.getAndVerifyUserSession(w, r)
	if !valid {
		return
	}

	session.Delete(stateCookieName)

	code := r.URL.Query().Get("code")
	if code == "" {
		common.WriteErrorResponsef(w, logger, http.StatusBadRequest, "missing code parameter")

		return nil, nil, false
	}

	oauthToken, err := o.oidcClient.Exchange(
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

func (o *Operation) getAndVerifyUserSession(w http.ResponseWriter, r *http.Request) (cookie.Jar, bool) {
	session, err := o.store.cookies.Open(r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to create or decode session cookie: %s", err.Error())

		return nil, false
	}

	stateCookie, found := session.Get(stateCookieName)
	if !found {
		common.WriteErrorResponsef(w, logger, http.StatusBadRequest, "missing state session cookie")

		return nil, false
	}

	state := r.URL.Query().Get("state")
	if state == "" {
		common.WriteErrorResponsef(w, logger, http.StatusBadRequest, "missing state parameter")

		return nil, false
	}

	if state != stateCookie {
		common.WriteErrorResponsef(w, logger, http.StatusBadRequest, "invalid state parameter")

		return nil, false
	}

	return session, true
}
