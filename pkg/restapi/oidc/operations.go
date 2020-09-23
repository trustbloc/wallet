/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/edge-core/pkg/storage"

	"github.com/trustbloc/edge-agent/pkg/restapi/common"
)

// Endpoints.
const (
	oidcLoginPath    = "/oidc/login"
	oidcCallbackPath = "/oidc/callback"
)

// Misc.
const (
	transientStoreName = "hubauth_trx"
)

var logger = log.New("hub-auth/oidc")

// Config holds all configuration for an Operation.
type Config struct {
	OIDC    *OIDCConfig
	Storage *StorageConfig
}

// OIDCConfig holds OIDC config.
type OIDCConfig struct {
	Provider     OIDCProvider
	ClientID     string
	ClientSecret string
	Scopes       []string
	CallbackURL  string
}

// StorageConfig holds storage config.
type StorageConfig struct {
	Storage          storage.Provider
	TransientStorage storage.Provider
}

// Operation implements OIDC operations.
type Operation struct {
	transientStore storage.Store
	oidcConfig     *OIDCConfig
}

// New returns a new Operation.
func New(config *Config) (*Operation, error) {
	op := &Operation{
		oidcConfig: config.OIDC,
	}

	var err error

	op.transientStore, err = openStore(config.Storage.TransientStorage, transientStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to open transient store: %w", err)
	}

	return op, nil
}

// GetRESTHandlers get all controller API handler available for this service.
func (o *Operation) GetRESTHandlers() []common.Handler {
	return []common.Handler{
		common.NewHTTPHandler(oidcLoginPath, http.MethodGet, o.oidcLoginHandler),
	}
}

func (o *Operation) oidcLoginHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("handling login request: %s", r.URL.String())

	oauth2Config := oauth2.Config{
		ClientID:     o.oidcConfig.ClientID,
		ClientSecret: o.oidcConfig.ClientSecret,
		Endpoint:     o.oidcConfig.Provider.Endpoint(),
		RedirectURL:  o.oidcConfig.CallbackURL,
		Scopes:       o.oidcConfig.Scopes,
	}

	state := uuid.New().String()

	err := newTransientData(o.transientStore).Put(state, state)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save state to store: %s", err.Error())

		return
	}

	redirectURL := oauth2Config.AuthCodeURL(state)

	http.Redirect(w, r, redirectURL, http.StatusFound)
	logger.Debugf("redirected to login url: %s", redirectURL)
}
