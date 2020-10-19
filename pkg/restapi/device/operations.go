/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package device

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/trustbloc/edge-agent/pkg/restapi/common"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/cookie"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/edge-core/pkg/storage"
)

// Endpoints.
const (
	registerBeginPath  = "register/begin"
	registerFinishPath = "register/finish"
)

// Stores.
const (
	deviceCookieName = "device_user"
)

var logger = log.New("edge-agent/device-registration")

// Config holds all configuration for an Operation.
type Config struct {
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
	Storage      storage.Provider
	SessionStore storage.Provider
}

type stores struct {
	users   *user.Store
	cookies cookie.Store
	session *session.Store
}

// Operation implements OIDC operations.
type Operation struct {
	store      *stores
	uiEndpoint string
	tlsConfig  *tls.Config
	webauthn   *webauthn.WebAuthn
}

// New returns a new Operation.
func New(config *Config) (*Operation, error) {
	op := &Operation{
		store: &stores{
			cookies: cookie.NewStore(config.Keys.Auth, config.Keys.Enc),
		},
		uiEndpoint: config.UIEndpoint,
		tlsConfig:  config.TLSConfig,
	}

	var err error

	op.store.session, err = session.NewStore()
	if err != nil {
		return nil, fmt.Errorf("failed to create web auth protocol session store: %w", err)
	}

	op.store.users, err = user.NewStore(config.Storage.Storage)
	if err != nil {
		return nil, fmt.Errorf("failed to open users store: %w", err)
	}

	return op, nil
}

// GetRESTHandlers get all controller API handler available for this service.
func (o *Operation) GetRESTHandlers() []common.Handler {
	return []common.Handler{
		common.NewHTTPHandler(registerBeginPath, http.MethodGet, o.beginRegistration),
		common.NewHTTPHandler(registerFinishPath, http.MethodPost, o.finishRegistration),
	}
}

func (o *Operation) beginRegistration(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("handling device registration: %s", r.URL.String())

	userData, canProceed := o.getUserData(w, r)
	if !canProceed {
		return
	}

	device := NewDevice(userData)

	webAuthnUser := protocol.UserEntity{
		ID:          device.WebAuthnID(),
		DisplayName: device.WebAuthnDisplayName(),
		CredentialEntity: protocol.CredentialEntity{
			Name: device.WebAuthnName(),
		},
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.User = webAuthnUser
		credCreationOpts.CredentialExcludeList = device.CredentialExcludeList()
	}

	// generate PublicKeyCredentialCreationOptions, session data
	protocolCredential, sessionData, err := o.webauthn.BeginRegistration(
		device,
		registerOptions,
	)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to begin registration %s", err.Error())

		return
	}

	protocolCredentialBytes, err := marshallResponse(protocolCredential)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadRequest, "failed to marshal protocol credential %s", err.Error())

		return
	}
	// store session data as marshaled JSON
	err = o.store.session.SaveWebauthnSession(userData.Sub, sessionData, r, w)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save web auth session %s", err.Error())

		return
	}

	r.Response = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(string(protocolCredentialBytes))),
	}

	logger.Debugf("Registration begins: %s", o.uiEndpoint)
}

func (o *Operation) finishRegistration(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("handling finish device registration: %s", r.URL.String())

	userData, canProceed := o.getUserData(w, r)
	if !canProceed {
		return
	}

	sessionData, err := o.store.session.GetWebauthnSession(userData.Sub, r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to get web auth session: %s", err.Error())

		return
	}

	device := NewDevice(userData)

	credential, err := o.webauthn.FinishRegistration(device, sessionData, r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to finish registration: %s", err.Error())

		return
	}

	credentialBytes, err := marshallResponse(credential)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to marshal finished credential %s", err.Error())

		return
	}
	// sending credential as part of response
	r.Response = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(string(credentialBytes))),
	}

	logger.Debugf("Registration success: %s", o.uiEndpoint)
}

func (o *Operation) getUserData(w http.ResponseWriter, r *http.Request) (userData *user.User, proceed bool) {
	cookieSession, err := o.store.cookies.Open(r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to create or decode session cookie: %s", err.Error())

		return nil, false
	}

	userSub, found := cookieSession.Get(deviceCookieName)
	if !found {
		common.WriteErrorResponsef(w, logger, http.StatusNotFound, "missing device user session cookie")

		return nil, false
	}

	userData, err = o.store.users.Get(fmt.Sprintf("%v", userSub))
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to read user session cookie: %s", err.Error())

		return nil, false
	}

	return userData, true
}

func marshallResponse(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}
