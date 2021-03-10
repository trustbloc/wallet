/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package device

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	ariesstorage "github.com/hyperledger/aries-framework-go/spi/storage"
	"github.com/trustbloc/edge-core/pkg/log"

	"github.com/trustbloc/edge-agent/pkg/restapi/common"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/cookie"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
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
	deviceCookieName  = "device_user"
)

var logger = log.New("edge-agent/device-registration")

// Config holds all configuration for an Operation.
type Config struct {
	Storage         *StorageConfig
	WalletDashboard string
	TLSConfig       *tls.Config
	Cookie          *cookie.Config
	Webauthn        *webauthn.WebAuthn
	HubAuthURL      string
}

// StorageConfig holds storage config.
type StorageConfig struct {
	Storage      ariesstorage.Provider
	SessionStore ariesstorage.Provider
}

type stores struct {
	users   *user.Store
	cookies cookie.Store
	storage ariesstorage.Store
	session *session.Store
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Operation implements OIDC operations.
type Operation struct {
	store           *stores
	walletDashboard string
	tlsConfig       *tls.Config
	httpClient      httpClient
	webauthn        *webauthn.WebAuthn
	hubAuthURL      string
}

// New returns a new Operation.
func New(config *Config) (*Operation, error) {
	op := &Operation{
		store: &stores{
			cookies: cookie.NewStore(config.Cookie),
		},
		tlsConfig:       config.TLSConfig,
		httpClient:      &http.Client{Transport: &http.Transport{TLSClientConfig: config.TLSConfig}},
		webauthn:        config.Webauthn,
		walletDashboard: config.WalletDashboard,
		hubAuthURL:      config.HubAuthURL,
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

	userData, canProceed := o.getUserData(w, r, userSubCookieName)
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
		credCreationOpts.Attestation = protocol.PreferDirectAttestation
	}

	// generate PublicKeyCredentialCreationOptions, session data
	credentialParams, sessionData, err := o.webauthn.BeginRegistration(
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

	jsonResponse(w, credentialParams, http.StatusOK)

	logger.Debugf("Registration begins")
}

func (o *Operation) finishRegistration(w http.ResponseWriter, r *http.Request) { // nolint:funlen // not clean to split
	logger.Debugf("handling finish device registration: %s", r.URL.String())

	userData, canProceed := o.getUserData(w, r, userSubCookieName)
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

	// unfold webauthn.FinishRegistration, to access parsedResponse
	parsedResponse, err := protocol.ParseCredentialCreationResponse(r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to finish registration: parsing ccr: %#v", err)

		return
	}

	credential, err := o.webauthn.CreateCredential(device, sessionData, parsedResponse)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to finish registration: cred: %#v", err)

		return
	}

	deviceCerts, ok := parsedResponse.Response.AttestationObject.AttStatement["x5c"].([]interface{})
	if !ok || len(deviceCerts) == 0 {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to finish registration: device certificate missing")

		return
	}

	err = o.requestDeviceValidation(r.Context(), userData.Sub, string(credential.Authenticator.AAGUID), deviceCerts)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to finish registration: %#v", err)

		return
	}

	device.AddCredential(*credential)

	err = o.saveDeviceInfo(device)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save device info: %s", err.Error())

		return
	}

	o.saveCookie(w, r, userData.Sub, deviceCookieName)

	jsonResponse(w, credential, http.StatusOK)

	logger.Debugf("Registration success")
}

func (o *Operation) requestDeviceValidation(ctx context.Context, userSub, aaguid string, certs []interface{}) error {
	if len(certs) == 0 {
		return fmt.Errorf("missing certs")
	}

	var certPemList []string

	for _, certInterface := range certs {
		cert, ok := certInterface.([]byte)
		if !ok {
			return fmt.Errorf("can't cast certificate data to []byte")
		}

		certPemList = append(certPemList, string(pem.EncodeToMemory(
			&pem.Block{Bytes: cert, Type: "CERTIFICATE"},
		)))
	}

	postData, err := json.Marshal(&struct {
		X5c    []string `json:"x5c"`
		Sub    string   `json:"sub"`
		Aaguid string   `json:"aaguid"`
	}{
		X5c:    certPemList,
		Sub:    userSub,
		Aaguid: aaguid,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal cert data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.hubAuthURL+"/device", bytes.NewBuffer(postData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	_, _, err = common.SendHTTPRequest(req, o.httpClient, http.StatusOK, nil)
	if err != nil {
		return fmt.Errorf("failed response from hub-auth/device endpoint: %w", err)
	}

	return nil
}

func (o *Operation) beginLogin(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("handling begin device login: %s", r.URL.String())

	// get username
	userData, canProceed := o.getUserData(w, r, deviceCookieName)
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

	logger.Debugf("Login begin success")
	jsonResponse(w, options, http.StatusOK)
}

func (o *Operation) finishLogin(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("handling finish device login: %s", r.URL.String())
	// get username
	userData, canProceed := o.getUserData(w, r, deviceCookieName)
	if !canProceed {
		return
	}

	deviceData, err := o.getDeviceInfo(userData.Sub)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadRequest, "failed to get device data: %s", err.Error())

		return
	}

	o.saveCookie(w, r, userData.Sub, userSubCookieName)

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
	http.Redirect(w, r, o.walletDashboard, http.StatusFound)
	logger.Debugf("Login finish success")
}

func (o *Operation) getUserData(w http.ResponseWriter, r *http.Request, cookieName string) (userData *user.User,
	proceed bool) {
	cookieSession, err := o.store.cookies.Open(r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to create or decode session cookie: %s", err.Error())

		return nil, false
	}

	userSub, found := cookieSession.Get(cookieName)
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

func (o *Operation) saveCookie(w http.ResponseWriter, r *http.Request, usr, cookieName string) {
	logger.Debugf("device cookie begin %s", usr)

	deviceSession, err := o.store.cookies.Open(r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to read user session cookie: %s", err.Error())

		return
	}

	deviceSession.Set(cookieName, usr)
	err = deviceSession.Save(r, w)

	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save device cookie: %s", err.Error())

		return
	}
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
