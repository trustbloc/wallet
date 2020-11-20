/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/trustbloc/edge-agent/pkg/restapi/common"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/cookie"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/tokens"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/edge-core/pkg/sss"
	"github.com/trustbloc/edge-core/pkg/sss/base"
	"github.com/trustbloc/edge-core/pkg/storage"
	"github.com/trustbloc/edv/pkg/client"
	"github.com/trustbloc/edv/pkg/restapi/models"
	"golang.org/x/oauth2"
)

// Endpoints.
const (
	oidcLoginPath    = "/login"
	oidcCallbackPath = "/callback"
	oidcUserInfoPath = "/userinfo"
	logoutPath       = "/logout"
)

// Stores.
const (
	transientStoreName      = "edgeagent_oidc_trx"
	stateCookieName         = "oauth2_state"
	userSubCookieName       = "user_sub"
	todoDeleteThisStoreName = "todo_delete"
)

var logger = log.New("hub-auth/oidc")

// Config holds all configuration for an Operation.
type Config struct {
	OIDCClient      oidc.Client
	Storage         *StorageConfig
	WalletDashboard string
	TLSConfig       *tls.Config
	Keys            *KeyConfig
	KeyServer       *KeyServerConfig
	UserEDVURL      string
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

// KeyServerConfig holds configuration for key management server.
type KeyServerConfig struct {
	AuthzKMSURL string
	OpsKMSURL   string
	KeyEDVURL   string
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type edvClient interface {
	CreateDataVault(config *models.DataVaultConfiguration, opts ...client.EDVOption) (string, error)
}

type stores struct {
	users             *user.Store
	tokens            *tokens.Store
	transient         storage.Store
	cookies           cookie.Store
	tempBootstrapData storage.Store
}

// Operation implements OIDC operations.
type Operation struct {
	store           *stores
	oidcClient      oidc.Client
	walletDashboard string
	tlsConfig       *tls.Config
	secretSplitter  sss.SecretSplitter
	httpClient      httpClient
	keyEDVClient    edvClient
	keyServer       *KeyServerConfig
	userEDVClient   edvClient
}

// New returns a new Operation.
func New(config *Config) (*Operation, error) {
	op := &Operation{
		oidcClient: config.OIDCClient,
		store: &stores{
			cookies: cookie.NewStore(config.Keys.Auth, config.Keys.Enc),
		},
		walletDashboard: config.WalletDashboard,
		tlsConfig:       config.TLSConfig,
		secretSplitter:  &base.Splitter{},
		httpClient:      &http.Client{Transport: &http.Transport{TLSClientConfig: config.TLSConfig}},
		keyEDVClient: client.New(
			config.KeyServer.KeyEDVURL,
			client.WithTLSConfig(config.TLSConfig),
		),
		keyServer: config.KeyServer,
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

	op.store.tempBootstrapData, err = store.Open(config.Storage.Storage, todoDeleteThisStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to open '%s': %w", todoDeleteThisStoreName, err)
	}

	if config.UserEDVURL != "" {
		op.userEDVClient = client.New(
			config.UserEDVURL,
			client.WithTLSConfig(config.TLSConfig),
		)
	}

	return op, nil
}

// GetRESTHandlers get all controller API handler available for this service.
func (o *Operation) GetRESTHandlers() []common.Handler {
	return []common.Handler{
		common.NewHTTPHandler(oidcLoginPath, http.MethodGet, o.oidcLoginHandler),
		common.NewHTTPHandler(oidcCallbackPath, http.MethodGet, o.oidcCallbackHandler),
		common.NewHTTPHandler(oidcUserInfoPath, http.MethodGet, o.userProfileHandler),
		common.NewHTTPHandler(logoutPath, http.MethodGet, o.userLogoutHandler),
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

	_, found := session.Get(userSubCookieName)
	if found {
		http.Redirect(w, r, o.walletDashboard, http.StatusMovedPermanently)

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
func (o *Operation) oidcCallbackHandler(w http.ResponseWriter, r *http.Request) { // nolint:funlen,gocyclo,lll // cannot reduce
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
		err = o.onboardUser(usr.Sub)
		if err != nil {
			common.WriteErrorResponsef(w, logger,
				http.StatusInternalServerError, "failed to onboard the user: %s", err.Error())

			return
		}

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

	session.Set(userSubCookieName, usr.Sub)

	err = session.Save(r, w)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save user sub cookie: %s", err.Error())

		return
	}

	http.Redirect(w, r, o.walletDashboard, http.StatusFound)
	logger.Debugf("redirected user to: %s", o.walletDashboard)
}

func (o *Operation) fetchTokens(
	w http.ResponseWriter, r *http.Request) (oauthToken *oauth2.Token, oidcToken oidc.Claimer, valid bool) {
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

func (o *Operation) userProfileHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("handling userprofile request")

	jar, err := o.store.cookies.Open(r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadRequest, "cannot open cookies: %s", err.Error())

		return
	}

	userSubCookie, found := jar.Get(userSubCookieName)
	if !found {
		common.WriteErrorResponsef(w, logger,
			http.StatusForbidden, "not logged in")

		return
	}

	userSub, ok := userSubCookie.(string)
	if !ok {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "invalid user sub cookie format")

		return
	}

	data, proceed := o.fetchUserData(w, r, userSub)
	if !proceed {
		return
	}

	common.WriteResponse(w, logger, data)
	logger.Debugf("finished handling userprofile request")
}

func (o *Operation) fetchUserData(w http.ResponseWriter, r *http.Request, sub string) (map[string]interface{}, bool) {
	tokns, err := o.store.tokens.Get(sub)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to fetch user tokens from store: %s", err.Error())

		return nil, false
	}

	userInfo, err := o.oidcClient.UserInfo(r.Context(), &oauth2.Token{
		AccessToken:  tokns.Access,
		TokenType:    "Bearer",
		RefreshToken: tokns.Refresh,
	})
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadGateway, "failed to fetch user info: %s", err.Error())

		return nil, false
	}

	data := make(map[string]interface{})

	err = userInfo.Claims(&data)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to extract claims from user info: %s", err.Error())

		return nil, false
	}

	data["bootstrap"], err = o.fetchBootstrapData(sub)
	if err != nil {
		common.WriteErrorResponsef(w, logger, http.StatusInternalServerError,
			"failed to fetch bootstrap data: %s", err.Error())

		return nil, false
	}

	return data, true
}

func (o *Operation) fetchBootstrapData(sub string) (*todoDeleteThisModel, error) {
	rawBootstrapData, err := o.store.tempBootstrapData.Get(sub)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bootstrap data: %w", err)
	}

	bootstrapData := &todoDeleteThisModel{}

	return bootstrapData, json.Unmarshal(rawBootstrapData, bootstrapData)
}

func (o *Operation) userLogoutHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("handling logout request")

	jar, err := o.store.cookies.Open(r)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusBadRequest, "cannot open cookies: %s", err.Error())

		return
	}

	_, found := jar.Get(userSubCookieName)
	if !found {
		logger.Infof("missing user cookie - this is a no-op")

		return
	}

	jar.Delete(userSubCookieName)

	err = jar.Save(r, w)
	if err != nil {
		common.WriteErrorResponsef(w, logger, http.StatusInternalServerError,
			"failed to delete user sub cookie: %s", err.Error())
	}

	logger.Debugf("finished handling logout request")
}

func (o *Operation) onboardUser(sub string) error {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return fmt.Errorf("create user secret key : %w", err)
	}

	secrets, err := o.secretSplitter.Split(b, 2, 2)
	if err != nil {
		return fmt.Errorf("split user secret key : %w", err)
	}

	// TODO https://github.com/trustbloc/edge-agent/issues/488 send half secret key to hub-auth and remove logger
	logger.Infof(string(secrets[0]))
	logger.Infof(string(secrets[1]))

	authzKeyStoreURL, err := createKeyStore(o.keyServer.AuthzKMSURL, sub, "", o.httpClient)
	if err != nil {
		return fmt.Errorf("create authz keystore : %w", err)
	}

	// TODO https://github.com/trustbloc/edge-agent/issues/493 create controller
	controller := uuid.New().URN()

	opsEDVVaultURL, err := createEDVDataVault(o.keyEDVClient, controller)
	if err != nil {
		return fmt.Errorf("create key edv vault : %w", err)
	}

	opsKeyStoreURL, err := createKeyStore(o.keyServer.OpsKMSURL, controller, opsEDVVaultURL, o.httpClient)
	if err != nil {
		return fmt.Errorf("create operational keystore : %w", err)
	}

	var userEDVVaultURL string

	if o.userEDVClient != nil {
		userEDVVaultURL, err = createEDVDataVault(o.userEDVClient, controller)
		if err != nil {
			return fmt.Errorf("create user edv vault : %w", err)
		}
	}

	// TODO https://github.com/trustbloc/edge-agent/issues/489 send keystore/vault ids to hub-auth instead of saving
	bits, err := json.Marshal(&todoDeleteThisModel{
		UserEDVVaultURL:  userEDVVaultURL,
		OpsEDVVaultURL:   opsEDVVaultURL,
		AuthzKeyStoreURL: authzKeyStoreURL,
		OpsKeyStoreURL:   opsKeyStoreURL,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal temporary bootstrap data model: %w", err)
	}

	err = o.store.tempBootstrapData.Put(sub, bits)
	if err != nil {
		return fmt.Errorf("failed to store temporary bootstrap data model: %w", err)
	}

	return nil
}

func createKeyStore(baseURL, controller, vaultID string, httpClient httpClient) (string, error) {
	reqBytes, err := json.Marshal(createKeystoreReq{
		Controller:         controller,
		OperationalVaultID: vaultID,
	})
	if err != nil {
		return "", fmt.Errorf("marshal create keystore req : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+"/kms/keystores", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}

	// TODO https://github.com/trustbloc/edge-agent/issues/488 pass half secret and oauth token to hub-kms
	_, headers, err := sendHTTPRequest(req, httpClient, http.StatusCreated)
	if err != nil {
		return "", fmt.Errorf("create authz keystore : %w", err)
	}

	keystoreURL := headers.Get("Location")

	return keystoreURL, nil
}

func createEDVDataVault(edvClient edvClient, controller string) (string, error) {
	config := models.DataVaultConfiguration{
		Sequence:    0,
		Controller:  controller,
		ReferenceID: uuid.New().String(),
		KEK:         models.IDTypePair{ID: uuid.New().URN(), Type: "AesKeyWrappingKey2019"},
		HMAC:        models.IDTypePair{ID: uuid.New().URN(), Type: "Sha256HmacKey2019"},
	}

	vaultURL, err := edvClient.CreateDataVault(&config)
	if err != nil {
		return "", fmt.Errorf("create data vault : %w", err)
	}

	return vaultURL, nil
}

func sendHTTPRequest(req *http.Request, httpClient httpClient, status int) ([]byte, http.Header, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("http request : %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			logger.Errorf("failed to close response body")
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("http request: failed to read resp body %d : %w", resp.StatusCode, err)
	}

	if resp.StatusCode != status {
		return nil, nil, fmt.Errorf("http request: expected=%d actual=%d body=%s", status, resp.StatusCode, string(body))
	}

	return body, resp.Header, nil
}
