/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package device

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/trustbloc/edge-agent/pkg/restapi/common"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/cookie"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/edge-core/pkg/storage"
)

// Endpoints.
const (
	registerBeginPath  = "/register/begin"
	registerFinishPath = "/register/finish"
	loginBeginPath     = "/login/begin"
	loginFinishPath    = "/login/finish"
)

// Stores.
const (
	deviceStoreName   = "edgeagent_device_trx"
	userSubCookieName = "user_sub"
)

var logger = log.New("edge-agent/device-registration")

// Config holds all configuration for an Operation.
type Config struct {
	Storage    *StorageConfig
	UIEndpoint string
	TLSConfig  *tls.Config
	Keys       *KeyConfig
	Webauthn   *webauthn.WebAuthn
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
	storage storage.Store
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
		webauthn:   config.Webauthn,
	}

	var err error

	op.store.storage, err = store.Open(config.Storage.Storage, deviceStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to open store: %w", err)
	}

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
		common.NewHTTPHandler(loginBeginPath, http.MethodGet, o.beginLogin),
		common.NewHTTPHandler(loginFinishPath, http.MethodPost, o.finishLogin),
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
	// store session data as marshaled JSON
	err = o.store.session.SaveWebauthnSession(userData.Sub, sessionData, r, w)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save web auth session %s", err.Error())

		return
	}

	jsonResponse(w, protocolCredential, http.StatusOK)

	logger.Infof("Registration begins")
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
			http.StatusInternalServerError, "failed to finish registration: %+v", err.Error())

		return
	}

	device.AddCredential(*credential)

	err = o.saveDeviceInfo(device)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save device info: %s", err.Error())

		return
	}

	jsonResponse(w, credential, http.StatusOK)

	logger.Infof("Registration success")
}

func (o *Operation) beginLogin(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("handling begin device login: %s", r.URL.String())

	// get username
	userData, canProceed := o.getUserData(w, r)
	if !canProceed {
		return
	}

	deviceData, err := o.getDeviceInfo(userData.Sub)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadRequest, "failed to get device data: %s", err.Error())

		return
	}

	// generate PublicKeyCredentialRequestOptions, session data
	options, sessionData, err := o.webauthn.BeginLogin(deviceData)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to begin login: %s", err.Error())

		return
	}

	// store session data as marshaled JSON
	err = o.store.session.SaveWebauthnSession(deviceData.ID, sessionData, r, w)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save web auth login session: %s", err.Error())

		return
	}

	jsonResponse(w, options, http.StatusOK)
}

func (o *Operation) finishLogin(w http.ResponseWriter, r *http.Request) {
	// get username
	userData, canProceed := o.getUserData(w, r)
	if !canProceed {
		return
	}

	deviceData, err := o.getDeviceInfo(userData.Sub)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadRequest, "failed to get device data: %s", err.Error())

		return
	}
	// load the session data
	sessionData, err := o.store.session.GetWebauthnSession(deviceData.ID, r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadRequest, "failed to get web auth login session: %s", err.Error())

		return
	}

	_, err = o.webauthn.FinishLogin(deviceData, sessionData, r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadRequest, "failed to finish login: %s", err.Error())

		return
	}

	// handle successful login
	jsonResponse(w, "Login Success", http.StatusOK)
}

func (o *Operation) getUserData(w http.ResponseWriter, r *http.Request) (userData *user.User, proceed bool) {
	cookieSession, err := o.store.cookies.Open(r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to create or decode session cookie: %s", err.Error())

		return nil, false
	}

	userSub, found := cookieSession.Get(userSubCookieName)
	if !found {
		common.WriteErrorResponsef(w, logger, http.StatusNotFound, "missing device user session cookie")

		return nil, false
	}

	userData, err = o.store.users.Get(fmt.Sprintf("%v", userSub))
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to get user data %s:", err.Error())

		return nil, false
	}

	displayName := strings.Split(fmt.Sprintf("%v", userSub), "@")[0]

	userData.FamilyName = displayName

	return userData, true
}

func (o *Operation) saveDeviceInfo(device *Device) error {
	deviceBytes, err := json.Marshal(device)
	if err != nil {
		return fmt.Errorf("failed to marshall the device data: %s", err)
	}

	err = o.store.storage.Put(device.ID, deviceBytes)

	if err != nil {
		return fmt.Errorf("failed to save the device data: %s", err)
	}

	return nil
}

func (o *Operation) getDeviceInfo(username string) (*Device, error) {
	// fetch user and check if user doesn't exist
	userData, err := o.store.users.Get(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %s", err.Error())
	}

	deviceDataBytes, err := o.store.storage.Get(userData.Sub)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch device data: %s", err.Error())
	}

	deviceData := &Device{}

	err = json.Unmarshal(deviceDataBytes, deviceData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall device data: %s", err.Error())
	}

	return deviceData, nil
}

func jsonResponse(w http.ResponseWriter, resp interface{}, c int) {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", respBytes)
}
