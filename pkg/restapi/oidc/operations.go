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
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/edge-core/pkg/storage"
	"golang.org/x/oauth2"

	"github.com/trustbloc/edge-agent/pkg/restapi/common"
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
)

var logger = log.New("hub-auth/oidc")

// Config holds all configuration for an Operation.
type Config struct {
	OIDCClient OIDCClient
	Storage    *StorageConfig
	UIEndpoint string
	TLSConfig  *tls.Config
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
	op := &Operation{
		oidcClient: config.OIDCClient,
		store:      &store{},
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

	state := uuid.New().String()

	err := newTransientData(o.store.transient).Put(state, state)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save state to store: %s", err.Error())

		return
	}

	redirectURL := o.oidcClient.FormatRequest(state)

	http.Redirect(w, r, redirectURL, http.StatusFound)
	logger.Debugf("redirected to login url: %s", redirectURL)
}

// TODO setup session cookies: https://github.com/trustbloc/edge-agent/issues/379
// TODO encrypt data before storing: https://github.com/trustbloc/edge-agent/issues/380
func (o *Operation) oidcCallbackHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("handling oidc callback: %s", r.URL.String())

	state := r.URL.Query().Get("state")
	if state == "" {
		common.WriteErrorResponsef(w, logger, http.StatusBadRequest, "missing state parameter")

		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		common.WriteErrorResponsef(w, logger, http.StatusBadRequest, "missing code parameter")

		return
	}

	_, err := o.store.transient.Get(state)
	if errors.Is(err, storage.ErrValueNotFound) {
		common.WriteErrorResponsef(w, logger, http.StatusBadRequest, "invalid state parameter")

		return
	}

	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "unable to query transient store: %s", err.Error())

		return
	}

	err = o.store.transient.Delete(state)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to delete state from transient store: %s", err.Error())

		return
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

		return
	}

	oidcToken, err := o.oidcClient.VerifyIDToken(r.Context(), oauthToken)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadGateway, "cannot verify id_token: %s", err.Error())

		return
	}

	// TODO only save new user if one doesn't already exist in the store for the given `sub`:
	//  https://github.com/trustbloc/edge-agent/issues/381
	user := &endUser{ID: uuid.New().String()}

	err = user.parse(oidcToken)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to parse id_token: %s", err.Error())

		return
	}

	err = newPersistedData(o.store.users).put(user.ID, user)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to persist user data: %s", err.Error())

		return
	}

	err = newPersistedData(o.store.tokens).put(user.ID, &endUserTokens{
		ID:      user.ID,
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
