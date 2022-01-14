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
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/fingerprint"
	ariesstorage "github.com/hyperledger/aries-framework-go/spi/storage"
	"github.com/igor-pavlenko/httpsignatures-go"
	"github.com/lafriks/go-shamir"
	"github.com/piprate/json-gold/ld"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/edge-core/pkg/zcapld"
	"github.com/trustbloc/edv/pkg/client"
	"github.com/trustbloc/edv/pkg/restapi/models"
	"golang.org/x/oauth2"

	"github.com/trustbloc/edge-agent/pkg/restapi/common"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/cookie"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/tokens"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
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
	authSecretPath        = "/secret"
	authBootstrapDataPath = "/bootstrap"
	createKeyStorePath    = "/v1/keystores"
	createDIDPath         = "/v1/keystores/did"
	keysPath              = "/v1/keystores/%s/keys"
	signPath              = "/v1/keystores/%s/keys/%s/sign"
)

const (
	edvResource           = "urn:edv:vault"
	providerQueryParam    = "provider"
	walletTokenExpiryMins = "20"
	actionCreateKey       = "createKey"
)

var logger = log.New("hub-auth/oidc")

// Config holds all configuration for an Operation.
type Config struct {
	OIDCClient      oidc.Client
	Storage         *StorageConfig
	WalletDashboard string
	TLSConfig       *tls.Config
	KeyServer       *KeyServerConfig
	UserEDVURL      string
	HubAuthURL      string
	JSONLDLoader    ld.DocumentLoader
	Cookie          *cookie.Config
}

// StorageConfig holds storage config.
type StorageConfig struct {
	Storage          ariesstorage.Provider
	TransientStorage ariesstorage.Provider
}

// KeyServerConfig holds configuration for key management server.
type KeyServerConfig struct {
	AuthzKMSURL string
	OpsKMSURL   string
	KeyEDVURL   string
}

type edvClient interface {
	CreateDataVault(config *models.DataVaultConfiguration, opts ...client.ReqOption) (string, []byte, error)
}

type stores struct {
	users     *user.Store
	tokens    *tokens.Store
	transient ariesstorage.Store
	cookies   cookie.Store
}

// Operation implements OIDC operations.
type Operation struct {
	store           *stores
	oidcClient      oidc.Client
	walletDashboard string
	tlsConfig       *tls.Config
	httpClient      common.HTTPClient
	keyEDVClient    edvClient
	keyServer       *KeyServerConfig
	userEDVClient   edvClient
	hubAuthURL      string
	jsonLDLoader    ld.DocumentLoader
	userEDVURL      string
}

// New returns a new Operation.
func New(config *Config) (*Operation, error) {
	op := &Operation{
		oidcClient: config.OIDCClient,
		store: &stores{
			cookies: cookie.NewStore(config.Cookie),
		},
		walletDashboard: config.WalletDashboard,
		tlsConfig:       config.TLSConfig,
		httpClient:      &http.Client{Transport: &http.Transport{TLSClientConfig: config.TLSConfig}},
		keyEDVClient: client.New(
			config.KeyServer.KeyEDVURL,
			client.WithTLSConfig(config.TLSConfig),
		),
		keyServer:    config.KeyServer,
		hubAuthURL:   config.HubAuthURL,
		jsonLDLoader: config.JSONLDLoader,
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
		op.userEDVURL = config.UserEDVURL
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

	// It is not mandatory parameter as we have to support login and signup flow for now Issue-785.
	providerID := r.URL.Query().Get(providerQueryParam)

	session, err := o.store.cookies.Open(r)
	if err != nil {
		// log the error and continue
		logger.Warnf("failed to read user session cookie: %s", err.Error())
	} else {
		_, found := session.Get(userSubCookieName)
		if found {
			http.Redirect(w, r, o.walletDashboard, http.StatusMovedPermanently)

			return
		}
	}

	state := uuid.New().String()
	session.Set(stateCookieName, state)

	err = session.Save(r, w)
	if err != nil {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to save session cookie: %s", err.Error())

		return
	}

	authOption := oauth2.SetAuthURLParam(providerQueryParam, providerID)

	redirectURL := o.oidcClient.FormatRequest(state, authOption)

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
	if err != nil && !errors.Is(err, ariesstorage.ErrDataNotFound) {
		common.WriteErrorResponsef(w, logger,
			http.StatusInternalServerError, "failed to query user data: %s", err.Error())

		return
	}

	if errors.Is(err, ariesstorage.ErrDataNotFound) {
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
		AccessToken: tokns.Access,
		SecretShare: walletUserData.SecretShare,
	}

	return data, true
}

func (o *Operation) fetchBootstrapData(accessToken string) (*userBootstrapData, error) {
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, o.hubAuthURL+authBootstrapDataPath, nil)
	if err != nil {
		return nil, err
	}

	addAccessToken(req, accessToken)

	data, _, err := common.SendHTTPRequest(req, o.httpClient, http.StatusOK, logger)
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

	secrets, err := shamir.Split(b, 2, 2)
	if err != nil {
		return "", fmt.Errorf("split user secret: %w", err)
	}

	walletSecretShare := secrets[0]
	authSecretShare := secrets[1]

	err = postSecret(o.hubAuthURL, accessToken, authSecretShare, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("post secret share to auth server: %w", err)
	}

	h := &kmsHeader{
		userSub:     sub,
		accessToken: accessToken,
		secretShare: walletSecretShare,
	}

	authzKeyStoreURL, err := createAuthzKeyStore(o.keyServer.AuthzKMSURL, sub, h, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("create authz keystore: %w", err)
	}

	authzKeyStoreID := getKeystoreID(authzKeyStoreURL)

	keyID, pubKey, err := createKey(authzKeyStoreURL, kms.ED25519, h, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("create authz key: %w", err)
	}

	_, controller := fingerprint.CreateDIDKey(pubKey)

	// EDV vault for storing user's keys
	kmsVaultURL, kmsEDVCapability, err := createEDVDataVault(o.keyEDVClient, controller, accessToken)
	if err != nil {
		return "", fmt.Errorf("create edv vault for kms: %w", err)
	}

	// create EDV controller on operational KMS
	edvController, err := createEDVController(o.keyServer.OpsKMSURL, accessToken, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("create edv controller: %w", err)
	}

	// create chain capabilities for KMS to use EDV storage
	edvZCAPs, err := createChainCapability(controller, getVaultID(kmsVaultURL), kmsEDVCapability, edvController,
		newKMSSigner(o.keyServer.AuthzKMSURL, authzKeyStoreID, keyID, h, o.httpClient), o.jsonLDLoader)
	if err != nil {
		return "", fmt.Errorf("create chain capability: %w", err)
	}

	opKeyStoreURL, opKeyStoreCapability, err := createOpKeyStore(o.keyServer.OpsKMSURL, controller, kmsVaultURL,
		edvZCAPs, accessToken, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("create operational key store: %w", err)
	}

	compressedOPSKMSCapability := base64.URLEncoding.EncodeToString(opKeyStoreCapability)

	var (
		userEDVVaultURL   string
		userEDVCapability []byte
	)

	if o.userEDVClient != nil {
		userEDVVaultURL, userEDVCapability, err = createEDVDataVault(o.userEDVClient, controller, accessToken)
		if err != nil {
			return "", fmt.Errorf("create user edv vault : %w", err)
		}
	}

	edvOpsKID, err := createOpKey(
		o.keyServer.OpsKMSURL,
		getKeystoreID(opKeyStoreURL),
		kms.NISTP256ECDHKW, // TODO make default key type configurable.
		controller,
		compressedOPSKMSCapability,
		newKMSSigner(o.keyServer.AuthzKMSURL, authzKeyStoreID, keyID, h, o.httpClient),
		h,
		o.httpClient)
	if err != nil {
		return "", fmt.Errorf("create edv operational key: %w", err)
	}

	edvOpsKIDURL := fmt.Sprintf("%s/keys/%s", opKeyStoreURL, edvOpsKID)

	hmacEDVKID, err := createOpKey(
		o.keyServer.OpsKMSURL,
		getKeystoreID(opKeyStoreURL),
		kms.HMACSHA256Tag256,
		controller,
		compressedOPSKMSCapability,
		newKMSSigner(o.keyServer.AuthzKMSURL, authzKeyStoreID, keyID, h, o.httpClient),
		h,
		o.httpClient,
	)
	if err != nil {
		return "", fmt.Errorf("create edv hmac key: %w", err)
	}

	hmacEDVKIDURL := fmt.Sprintf("%s/keys/%s", opKeyStoreURL, hmacEDVKID)

	// TODO remove OPSKMSCapability: https://github.com/trustbloc/edge-agent/issues/583.
	data := &BootstrapData{
		User:              uuid.NewString(),
		UserEDVVaultURL:   userEDVVaultURL, // TODO to be removed after universal wallet migration
		OpsEDVVaultURL:    kmsVaultURL,     // TODO to be removed after universal wallet migration
		AuthzKeyStoreURL:  authzKeyStoreURL,
		OpsKeyStoreURL:    opKeyStoreURL,
		EDVOpsKIDURL:      edvOpsKIDURL,
		EDVHMACKIDURL:     hmacEDVKIDURL,
		UserEDVCapability: string(userEDVCapability),
		OPSKMSCapability:  compressedOPSKMSCapability,
		UserEDVVaultID:    getVaultID(userEDVVaultURL),
		UserEDVServer:     o.userEDVURL,
		UserEDVEncKID:     edvOpsKID,
		UserEDVMACKID:     hmacEDVKID,
		TokenExpiry:       walletTokenExpiryMins,
	}

	err = postUserBootstrapData(o.hubAuthURL, accessToken, data, o.httpClient)
	if err != nil {
		return "", fmt.Errorf("update user bootstrap data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(walletSecretShare), nil
}

func postSecret(baseURL, accessToken string, secret []byte, httpClient common.HTTPClient) error {
	reqBytes, err := json.Marshal(secretRequest{
		Secret: secret,
	})
	if err != nil {
		return fmt.Errorf("marshal secret req : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+authSecretPath, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	addAccessToken(req, accessToken)

	_, _, err = common.SendHTTPRequest(req, httpClient, http.StatusOK, logger)
	if err != nil {
		return err
	}

	return nil
}

func postUserBootstrapData(baseURL, accessToken string, data *BootstrapData, httpClient common.HTTPClient) error {
	reqBytes, err := json.Marshal(userBootstrapData{
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("marshal boostrap data : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+authBootstrapDataPath, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	addAccessToken(req, accessToken)

	_, _, err = common.SendHTTPRequest(req, httpClient, http.StatusOK, logger)
	if err != nil {
		return err
	}

	return nil
}

func createAuthzKeyStore(baseURL, controller string, h *kmsHeader, httpClient common.HTTPClient) (string, error) {
	reqBytes, err := json.Marshal(createKeyStoreReq{
		Controller: controller,
	})
	if err != nil {
		return "", fmt.Errorf("marshal create keystore req : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+createKeyStorePath, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.accessToken))
	req.Header.Set("Secret-Share", base64.StdEncoding.EncodeToString(h.secretShare))

	respBody, _, err := common.SendHTTPRequest(req, httpClient, http.StatusOK, logger)
	if err != nil {
		return "", fmt.Errorf("create authz key store: %w", err)
	}

	var resp createKeyStoreResp

	if err = json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("unmarshal create key store resp: %w", err)
	}

	return resp.KeyStoreURL, nil
}

func createEDVController(baseURL, accessToken string, httpClient common.HTTPClient) (string, error) {
	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+createDIDPath, nil)
	if err != nil {
		return "", fmt.Errorf("new create did req: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	respBody, _, err := common.SendHTTPRequest(req, httpClient, http.StatusOK, logger)
	if err != nil {
		return "", fmt.Errorf("create edv controller on kms: %w", err)
	}

	var resp createDIDResp

	if err = json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("unmarshal create did resp: %w", err)
	}

	return resp.DID, nil
}

func createChainCapability(controller, vaultID string, edvCapability []byte, kmsDIDKey string, s signer,
	documentLoader ld.DocumentLoader) ([]byte, error) {
	capability, err := zcapld.ParseCapability(edvCapability)
	if err != nil {
		return nil, fmt.Errorf("parse edv capability: %w", err)
	}

	chainCapability, err := zcapld.NewCapability(&zcapld.Signer{
		SignatureSuite:     ed25519signature2018.New(suite.WithSigner(s)),
		SuiteType:          ed25519signature2018.SignatureType,
		VerificationMethod: controller,
		ProcessorOpts:      []jsonld.ProcessorOpts{jsonld.WithDocumentLoader(documentLoader)},
	}, zcapld.WithParent(capability.ID), zcapld.WithInvoker(kmsDIDKey),
		zcapld.WithAllowedActions("read", "write"),
		zcapld.WithInvocationTarget(vaultID, edvResource),
		zcapld.WithCapabilityChain(capability.Parent, capability.ID))
	if err != nil {
		return nil, fmt.Errorf("new capability: %w", err)
	}

	b, err := json.Marshal(chainCapability)
	if err != nil {
		return nil, fmt.Errorf("marshal capability: %w", err)
	}

	return b, nil
}

func createOpKeyStore(baseURL, controller, vaultURL string, edvZCAPs []byte, accessToken string,
	httpClient common.HTTPClient) (string, []byte, error) {
	reqBytes, err := json.Marshal(createKeyStoreReq{
		Controller: controller,
		EDV: &edvOptions{
			VaultURL:   "https://" + vaultURL,
			Capability: edvZCAPs,
		},
	})
	if err != nil {
		return "", nil, fmt.Errorf("marshal create keystore req: %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+createKeyStorePath, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	respBody, _, err := common.SendHTTPRequest(req, httpClient, http.StatusOK, logger)
	if err != nil {
		return "", nil, fmt.Errorf("create ops key store: %w", err)
	}

	var resp createKeyStoreResp

	if err = json.Unmarshal(respBody, &resp); err != nil {
		return "", nil, fmt.Errorf("unmarshal create key store resp: %w", err)
	}

	return resp.KeyStoreURL, resp.Capability, nil
}

func sign(r *http.Request, controller, invocationAction, compressedKMSCapability string, s signer) error {
	r.Header.Set(
		zcapld.CapabilityInvocationHTTPHeader,
		fmt.Sprintf(`zcap capability="%s",action="%s"`, compressedKMSCapability, invocationAction),
	)

	hs := httpsignatures.NewHTTPSignatures(&zcapld.AriesDIDKeySecrets{})
	hs.SetSignatureHashAlgorithm(&zcapld.AriesDIDKeySignatureHashAlgorithm{
		KMS:    &zcapRemoteKMS{},
		Crypto: &zcapRemoteCrypto{signer: s},
	})

	if err := hs.Sign(controller, r); err != nil {
		return fmt.Errorf("failed to sign http request: %w", err)
	}

	return nil
}

func createKey(keyStoreURL, keyType string, h *kmsHeader, httpClient common.HTTPClient) (string, []byte, error) {
	b, err := json.Marshal(createKeyReq{
		KeyType: keyType,
	})
	if err != nil {
		return "", nil, fmt.Errorf("marshal create key req : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, keyStoreURL+"/keys", bytes.NewBuffer(b))
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.accessToken))
	req.Header.Set("Secret-Share", base64.StdEncoding.EncodeToString(h.secretShare))

	respBody, _, err := common.SendHTTPRequest(req, httpClient, http.StatusOK, logger)
	if err != nil {
		return "", nil, fmt.Errorf("create authz key: %w", err)
	}

	var resp createKeyResp

	if err = json.Unmarshal(respBody, &resp); err != nil {
		return "", nil, fmt.Errorf("unmarshal create key resp: %w", err)
	}

	return getKeyID(resp.KeyURL), resp.PublicKey, nil
}

func createOpKey(baseURL, keystoreID, keyType, controller, compressedKMSCapability string,
	s signer, h *kmsHeader, httpClient common.HTTPClient) (string, error) {
	reqBytes, err := json.Marshal(createKeyReq{
		KeyType: keyType,
	})
	if err != nil {
		return "", fmt.Errorf("marshal create key req : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, baseURL+fmt.Sprintf(keysPath, keystoreID), bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.accessToken))
	// TODO this secret should be signed over below
	req.Header.Add("Secret-Share", base64.StdEncoding.EncodeToString(h.secretShare))

	err = sign(req, controller, actionCreateKey, compressedKMSCapability, s)
	if err != nil {
		return "", fmt.Errorf("sign req: %w", err)
	}

	respBody, _, err := common.SendHTTPRequest(req, httpClient, http.StatusOK, logger)
	if err != nil {
		return "", fmt.Errorf("create key: %w", err)
	}

	var resp createKeyResp

	if err = json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("unmarshal create key resp: %w", err)
	}

	return getKeyID(resp.KeyURL), nil
}

func getKeystoreID(location string) string {
	const (
		keystoreIDPos = 5
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
		keyIDPos = 7
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

func addAccessToken(r *http.Request, token string) {
	r.Header.Set(
		"authorization",
		fmt.Sprintf("Bearer %s", base64.StdEncoding.EncodeToString([]byte(token))),
	)
}
