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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/fingerprint"
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
	"github.com/trustbloc/edge-core/pkg/zcapld"
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
	transientStoreName = "edgeagent_oidc_trx"
	stateCookieName    = "oauth2_state"
	userSubCookieName  = "user_sub"
)

// external url paths.
const (
	hubAuthSecretPath        = "/secret"
	hubAuthBootstrapDataPath = "/bootstrap"
	hubKMSCreateKeyStorePath = "/kms/keystores"
	keysEndpoint             = "/kms/keystores/%s/keys"
	exportKeyEndpoint        = "/kms/keystores/%s/keys/%s/export"
	capabilityEndpoint       = "/kms/keystores/%s/capability"
	signEndpoint             = "/kms/keystores/%s/keys/%s/sign"
)

const (
	edvResource = "urn:edv:vault"
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
	HubAuthURL      string
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
	CreateDataVault(config *models.DataVaultConfiguration, opts ...client.ReqOption) (string, []byte, error)
}

type stores struct {
	users     *user.Store
	tokens    *tokens.Store
	transient storage.Store
	cookies   cookie.Store
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
	hubAuthURL      string
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
		keyServer:  config.KeyServer,
		hubAuthURL: config.HubAuthURL,
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
		walletSecretShare, onboardErr := o.onboardUser(usr.Sub, oauthToken.AccessToken)
		if onboardErr != nil {
			common.WriteErrorResponsef(w, logger,
				http.StatusInternalServerError, "failed to onboard the user: %s", onboardErr.Error())

			return
		}

		usr.SecretShare = walletSecretShare

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

	walletUserData, err := o.store.users.Get(sub)
	if err != nil {
		common.WriteErrorResponsef(w, logger, http.StatusInternalServerError,
			"failed to fetch bootstrap data: %s", err.Error())

		return nil, false
	}

	userBootStrapData, err := o.fetchBootstrapData(tokns.Access)
	if err != nil {
		common.WriteErrorResponsef(w, logger, http.StatusInternalServerError,
			"failed to fetch bootstrap data: %s", err.Error())

		return nil, false
	}

	data["bootstrap"] = userBootStrapData.Data
	data["userConfig"] = &userConfig{
		Sub:         sub,
		SecretShare: walletUserData.SecretShare,
	}

	return data, true
}

func (o *Operation) fetchBootstrapData(accessToken string) (*userBootstrapData, error) {
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, o.hubAuthURL+hubAuthBootstrapDataPath, nil)
	if err != nil {
		return nil, err
	}

	addAccessToken(req, accessToken)

	data, _, err := sendHTTPRequest(req, o.httpClient, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("get bootstrap data : %w", err)
	}

	bootstrapData := &userBootstrapData{}

	return bootstrapData, json.Unmarshal(data, bootstrapData)
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

func (o *Operation) onboardUser(sub, accessToken string) (string, error) { // nolint:funlen,gocyclo // not much logic
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("create user secret key : %w", err)
	}

	secrets, err := o.secretSplitter.Split(b, 2, 2)
	if err != nil {
		return "", fmt.Errorf("split user secret key : %w", err)
	}

	walletSecretShare := base64.StdEncoding.EncodeToString(secrets[0])
	hubAuthSecretShare := secrets[1]

	err = postSecret(o.hubAuthURL, accessToken, hubAuthSecretShare, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("post half secret to hub-auth : %w", err)
	}

	h := &hubKMSHeader{
		userSub:     sub,
		accessToken: accessToken,
		secretShare: walletSecretShare,
	}

	authzKeyStoreURL, _, err := createKeyStore(o.keyServer.AuthzKMSURL, sub, "", h, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("create authz keystore : %w", err)
	}

	authzKeyStoreID := getKeystoreID(authzKeyStoreURL)

	keyID, err := createKey(o.keyServer.AuthzKMSURL, authzKeyStoreID, kms.ED25519, h, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("failed create authz key : %w", err)
	}

	pkBytes, err := exportPublicKey(o.keyServer.AuthzKMSURL, authzKeyStoreID, keyID, h, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("failed export public key: %w", err)
	}

	_, controller := fingerprint.CreateDIDKey(pkBytes)

	opsEDVVaultURL, opsEDVCapability, err := createEDVDataVault(o.keyEDVClient, controller, accessToken)
	if err != nil {
		return "", fmt.Errorf("create edv vault : %w", err)
	}

	opsEDVVaultID := getVaultID(opsEDVVaultURL)

	opsKeyStoreURL, opsKeyStoreEDVDIDKey, err := createKeyStore(o.keyServer.OpsKMSURL, controller,
		opsEDVVaultID, &hubKMSHeader{accessToken: accessToken}, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("create operational keystore : %w", err)
	}

	if len(opsEDVCapability) != 0 {
		if errUpdate := updateEDVCapabilityInKeyStore(o.keyServer.OpsKMSURL, getKeystoreID(opsKeyStoreURL), controller,
			opsEDVVaultID, opsEDVCapability, opsKeyStoreEDVDIDKey, newKMSSigner(o.keyServer.AuthzKMSURL,
				authzKeyStoreID, keyID, h, o.httpClient), o.httpClient); errUpdate != nil {
			return "", errUpdate
		}
	}

	var userEDVVaultURL string

	var userEDVCapability []byte

	if o.userEDVClient != nil {
		userEDVVaultURL, userEDVCapability, err = createEDVDataVault(o.userEDVClient, controller, accessToken)
		if err != nil {
			return "", fmt.Errorf("create user edv vault : %w", err)
		}
	}

	edvOpsKID, err := createKey(o.keyServer.OpsKMSURL, getKeystoreID(opsKeyStoreURL), kms.ECDH256KWAES256GCM, h,
		o.httpClient)
	if err != nil {
		return "", fmt.Errorf("create edv operational key : %w", err)
	}

	edvOpsKIDURL := fmt.Sprintf("%s/keys/%s", opsKeyStoreURL, edvOpsKID)

	hmacEDVKID, err := createKey(o.keyServer.OpsKMSURL, getKeystoreID(opsKeyStoreURL), kms.HMACSHA256Tag256, h,
		o.httpClient)
	if err != nil {
		return "", fmt.Errorf("create edv hmac key : %w", err)
	}

	hmacEDVKIDURL := fmt.Sprintf("%s/keys/%s", opsKeyStoreURL, hmacEDVKID)

	data := &BootstrapData{
		UserEDVVaultURL:   userEDVVaultURL,
		OpsEDVVaultURL:    opsEDVVaultURL,
		AuthzKeyStoreURL:  authzKeyStoreURL,
		OpsKeyStoreURL:    opsKeyStoreURL,
		EDVOpsKIDURL:      edvOpsKIDURL,
		EDVHMACKIDURL:     hmacEDVKIDURL,
		UserEDVCapability: string(userEDVCapability),
	}

	err = postUserBootstrapData(o.hubAuthURL, accessToken, data, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("update user bootstrap data : %w", err)
	}

	return walletSecretShare, nil
}

func postSecret(baseURL, accessToken string, secret []byte, httpClient httpClient) error {
	reqBytes, err := json.Marshal(secretRequest{
		Secret: secret,
	})
	if err != nil {
		return fmt.Errorf("marshal secret req : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+hubAuthSecretPath, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	addAccessToken(req, accessToken)

	_, _, err = sendHTTPRequest(req, httpClient, http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func postUserBootstrapData(baseURL, accessToken string, data *BootstrapData, httpClient httpClient) error {
	reqBytes, err := json.Marshal(userBootstrapData{
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("marshal boostrap data : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+hubAuthBootstrapDataPath, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	addAccessToken(req, accessToken)

	_, _, err = sendHTTPRequest(req, httpClient, http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func createKeyStore(baseURL, controller, vaultID string, h *hubKMSHeader,
	httpClient httpClient) (string, string, error) {
	reqBytes, err := json.Marshal(createKeystoreReq{
		Controller: controller,
		VaultID:    vaultID,
	})
	if err != nil {
		return "", "", fmt.Errorf("marshal create keystore req : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+hubKMSCreateKeyStorePath, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", "", err
	}

	addAuthZKMSHeaders(req, h)

	_, headers, err := sendHTTPRequest(req, httpClient, http.StatusCreated)
	if err != nil {
		return "", "", fmt.Errorf("create authz keystore : %w", err)
	}

	keystoreURL := headers.Get("Location")
	edvDIDKey := headers.Get("Edvdidkey")

	return keystoreURL, edvDIDKey, nil
}

func updateEDVCapabilityInKeyStore(baseURL, keystoreID, controller, vaultID string, edvCapability []byte,
	kmsDIDKey string, s signer, httpClient httpClient) error {
	capability, err := zcapld.ParseCapability(edvCapability)
	if err != nil {
		return err
	}

	chainCapability, err := zcapld.NewCapability(&zcapld.Signer{
		SignatureSuite:     ed25519signature2018.New(suite.WithSigner(s)),
		SuiteType:          ed25519signature2018.SignatureType,
		VerificationMethod: controller,
	}, zcapld.WithParent(capability.ID), zcapld.WithInvoker(kmsDIDKey),
		zcapld.WithAllowedActions("read", "write"),
		zcapld.WithInvocationTarget(vaultID, edvResource),
		zcapld.WithCapabilityChain(capability.Parent, capability.ID))
	if err != nil {
		return err
	}

	chainCapabilityBytes, err := json.Marshal(chainCapability)
	if err != nil {
		return err
	}

	reqBytes, err := json.Marshal(updateCapabilityReq{
		EDVCapability: chainCapabilityBytes,
	})
	if err != nil {
		return fmt.Errorf("marshal create update capability req : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+fmt.Sprintf(capabilityEndpoint, keystoreID), bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	_, _, err = sendHTTPRequest(req, httpClient, http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to update edv capability keystore : %w", err)
	}

	return nil
}

func createKey(baseURL, keystoreID, keyType string, h *hubKMSHeader, httpClient httpClient) (string, error) {
	reqBytes, err := json.Marshal(createKeyReq{
		KeyType: keyType,
	})
	if err != nil {
		return "", fmt.Errorf("marshal create key req : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+fmt.Sprintf(keysEndpoint, keystoreID), bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}

	addAuthZKMSHeaders(req, h)

	_, headers, err := sendHTTPRequest(req, httpClient, http.StatusCreated)
	if err != nil {
		return "", fmt.Errorf("create authz key : %w", err)
	}

	return getKeyID(headers.Get("Location")), nil
}

func exportPublicKey(baseURL, keystoreID, keyID string, h *hubKMSHeader, httpClient httpClient) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodGet, baseURL+fmt.Sprintf(exportKeyEndpoint, keystoreID, keyID), nil)
	if err != nil {
		return nil, err
	}

	addAuthZKMSHeaders(req, h)

	resp, _, err := sendHTTPRequest(req, httpClient, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("failed to export authz key : %w", err)
	}

	var exportKey exportKeyResp

	if errUnmarshal := json.Unmarshal(resp, &exportKey); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	pkBytes, err := base64.URLEncoding.DecodeString(exportKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return pkBytes, nil
}

func getKeystoreID(location string) string {
	const (
		keystoreIDPos = 3
	)

	s := strings.Split(location, "/")

	keystoreID := ""
	if len(s) > keystoreIDPos {
		keystoreID = s[keystoreIDPos]
	}

	return keystoreID
}

func getKeyID(location string) string {
	const (
		keyIDPos = 5
	)

	s := strings.Split(location, "/")

	keyID := ""
	if len(s) > keyIDPos {
		keyID = s[keyIDPos]
	}

	return keyID
}

func getVaultID(vaultURL string) string {
	parts := strings.Split(vaultURL, "/")

	return parts[len(parts)-1]
}

func createEDVDataVault(edvClient edvClient, controller, accessToken string) (string, []byte, error) {
	config := models.DataVaultConfiguration{
		Sequence:    0,
		Controller:  controller,
		ReferenceID: uuid.New().String(),
		KEK:         models.IDTypePair{ID: uuid.New().URN(), Type: "AesKeyWrappingKey2019"},
		HMAC:        models.IDTypePair{ID: uuid.New().URN(), Type: "Sha256HmacKey2019"},
	}

	vaultURL, capability, err := edvClient.CreateDataVault(&config,
		client.WithRequestHeader(func(req *http.Request) (*http.Header, error) {
			req.Header.Set("Authorization", "Bearer "+accessToken)

			return &req.Header, nil
		}))
	if err != nil {
		return "", nil, fmt.Errorf("create data vault : %w", err)
	}

	return vaultURL, capability, nil
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

func addAuthZKMSHeaders(r *http.Request, h *hubKMSHeader) {
	r.Header.Add("Hub-Kms-Secret", h.secretShare)
	r.Header.Add("Hub-Kms-User", h.userSub)
	logger.Errorf(h.secretShare)
	logger.Errorf(h.userSub)

	addAccessToken(r, h.accessToken)
}

func addAccessToken(r *http.Request, token string) {
	r.Header.Set(
		"authorization",
		fmt.Sprintf("Bearer %s", base64.StdEncoding.EncodeToString([]byte(token))),
	)
}
