/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/cucumber/godog"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/google/uuid"

	"github.com/trustbloc/edge-agent/pkg/restapi/oidc"
	"github.com/trustbloc/edge-agent/test/bdd/mock/authenticator"
	"github.com/trustbloc/edge-agent/test/bdd/pkg/bddcontext"
)

// HTTP server.
const (
	host                         = "https://bddtest-wallet-web.trustbloc.local:8077"
	loginPath                    = host + "/oidc/login"
	userProfilePath              = host + "/oidc/userinfo"
	userLogoutPath               = host + "/oidc/logout"
	dashboardPath                = "https://localhost:8078/dashboard"
	hubAuthHost                  = "https://localhost:8045"
	hubAuthSelectOIDCProviderURL = hubAuthHost + "/oauth2/login"
	mockLoginURL                 = "https://localhost:8099/mock/login"
	mockOIDCProviderName         = "mockbank" // providers.yaml
)

// device login endpoints.
const (
	registerBeginPath     = host + "/device/register/begin"
	registerFinishPath    = host + "/device/register/finish"
	deviceloginBeginPath  = host + "/device/login/begin"
	deviceloginFinishPath = host + "/device/login/finish"
)

// Mock Login Consent App.
const (
	mockLoginEndpoint = "https://localhost:8099/mock/login"
	mockAuthNEndpoint = "https://localhost:8099/mock/authn"
	mockAuthZEndpoint = "https://localhost:8099/mock/authz"
)

// AuthConfigRequest can be sent to mockAuthNEndpoint to configure the user authentication behaviour.
type AuthConfigRequest struct {
	Sub  string `json:"sub"`
	Fail bool   `json:"fail,omitempty"`
}

// ConsentConfigRequest can be sent to mockAuthZEndpoint to configure the user authorization behaviour.
type ConsentConfigRequest struct {
	UserClaims *UserClaims `json:"user_claims,omitempty"`
	Fail       bool        `json:"fail,omitempty"`
}

// UserClaims can be used by BDD tests to configure the claims in the id_token.
type UserClaims struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
}

type httpResponse struct {
	statusCode int
	body       string
	url        string
}

// NewSteps returns login BDD test steps.
func NewSteps(ctx *bddcontext.BDDContext) *Steps {
	return &Steps{
		ctx:               ctx,
		mockAuthenticator: authenticator.New(),
	}
}

// Steps are the login BDD test steps.
type Steps struct {
	browser             *http.Client
	ctx                 *bddcontext.BDDContext
	loginRedirectResult *httpResponse
	expectedUserClaims  *UserClaims
	authNResult         *httpResponse
	authZResult         *httpResponse
	mockAuthenticator   *authenticator.MockAuthenticator
}

// Register the login BDD test steps in the godog suite.
func (s *Steps) Register(gs *godog.Suite) {
	gs.Step("the user clicks on the Login button", s.userClicksLoginButton)
	gs.Step("the user is redirected to the OIDC provider", s.userRedirectedToOIDCProvider)
	gs.Step("the user is authenticated", s.userIsAuthenticated)
	gs.Step("the user consents to sharing their identity data", s.userAuthorizesAccessToTheirData)
	gs.Step("the user is redirected to the wallet's dashboard", s.userRedirectedToWallet)
	gs.Step("the user can retrieve their profile", s.userRetrievesProfile)
	gs.Step("the user is logged in", s.userIsLoggedIn)
	gs.Step("the user logs out", s.userLogout)
	gs.Step("the user cannot access their profile", s.userCannotAccessProfile)
	gs.Step("^the user registers a device$", s.registerDevice)
	gs.Step("^the user logs in through their device$", s.loginDevice)
}

func (s *Steps) userClicksLoginButton() error {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("failed to initialize the user's cookie jar : %w", err)
	}

	s.browser = &http.Client{
		Transport: &http.Transport{TLSClientConfig: s.ctx.TLSConfig},
		Jar:       cookieJar,
	}

	resp, err := s.browser.Get(loginPath) // nolint:noctx // no need to set context
	if err != nil {
		return fmt.Errorf("failed to invoke http server login endpoint %s: %w", loginPath, err)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close HTTP response body: %w", err)
	}

	request := fmt.Sprintf("%s?provider=%s", hubAuthSelectOIDCProviderURL, mockOIDCProviderName)

	resp, err = s.browser.Get(request) // nolint:noctx // no need to set context
	if err != nil {
		return fmt.Errorf("user failed to select OIDC provider using request %s: %w", request, err)
	}

	if !strings.HasPrefix(resp.Request.URL.String(), mockLoginURL) {
		return fmt.Errorf(
			"user at wrong third party OIDC provider; expected %s got %s",
			mockLoginURL, resp.Request.URL.String(),
		)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close HTTP response body: %w", err)
	}

	s.loginRedirectResult = &httpResponse{
		statusCode: resp.StatusCode,
		body:       string(body),
		url:        resp.Request.URL.String(),
	}

	return nil
}

func (s *Steps) userRedirectedToOIDCProvider() error {
	if s.loginRedirectResult.statusCode != http.StatusOK {
		return fmt.Errorf(
			"expected StatusCode=%d on login redirection but got %d, msg=%s",
			http.StatusOK, s.loginRedirectResult.statusCode, s.loginRedirectResult.body)
	}

	if s.loginRedirectResult.body != "mock UI" {
		return fmt.Errorf("unexpected UI rendered: %s", s.loginRedirectResult.body)
	}

	if !strings.HasPrefix(s.loginRedirectResult.url, mockLoginEndpoint) {
		return fmt.Errorf("expected login URL %s but got %s", mockLoginEndpoint, s.loginRedirectResult.url)
	}

	return nil
}

func (s *Steps) userIsAuthenticated() error {
	s.expectedUserClaims = &UserClaims{
		Sub:        uuid.New().String(),
		Name:       "John Doe",
		GivenName:  "John",
		FamilyName: "Doe",
		Email:      "john.doe@test.com",
	}

	request, err := json.Marshal(&AuthConfigRequest{
		Sub: s.expectedUserClaims.Sub,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal auth config request: %w", err)
	}

	// nolint:noctx // no need to set the context
	resp, err := s.browser.Post(mockAuthNEndpoint, "application/json", bytes.NewReader(request))
	if err != nil {
		return fmt.Errorf("failed to post request to mock auth endpoint: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}

	s.authNResult = &httpResponse{
		statusCode: resp.StatusCode,
		body:       string(body),
		url:        resp.Request.URL.String(),
	}

	return nil
}

func (s *Steps) userAuthorizesAccessToTheirData() error {
	request, err := json.Marshal(&ConsentConfigRequest{
		UserClaims: s.expectedUserClaims,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal consent config request: %w", err)
	}

	// nolint:noctx // no need to set the context
	resp, err := s.browser.Post(mockAuthZEndpoint, "application/json", bytes.NewReader(request))
	if err != nil {
		return fmt.Errorf("failed to invoke mock authZ endpoint %s: %w", mockAuthZEndpoint, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close http response body: %w", err)
	}

	s.authZResult = &httpResponse{
		statusCode: resp.StatusCode,
		body:       string(body),
		url:        resp.Request.URL.String(),
	}

	return nil
}

func (s *Steps) userRedirectedToWallet() error {
	if s.authZResult.statusCode != http.StatusOK {
		fmt.Printf("url: %s\n", s.authZResult.url)
		fmt.Printf("body: %s\n", s.authZResult.body)

		return fmt.Errorf("expected status code %d but got %d", http.StatusFound, s.authZResult.statusCode)
	}

	if !strings.HasPrefix(s.authZResult.url, dashboardPath) {
		return fmt.Errorf("expected path %s but got %s", dashboardPath, s.authZResult.url)
	}

	return nil
}

func (s *Steps) userRetrievesProfile() error { // nolint:gocyclo,funlen // not much logic
	response, err := s.browser.Get(userProfilePath) // nolint:noctx // ignoring rule since these are bdd tests
	if err != nil {
		return fmt.Errorf("user failed to fetch their profile: %w", err)
	}

	profile := make(map[string]interface{})

	err = json.NewDecoder(response.Body).Decode(&profile)
	if err != nil {
		return fmt.Errorf("user failed to decode their profile: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}

	sub, found := profile["sub"]
	if !found {
		return fmt.Errorf("user did not receive their 'sub' value in their profile")
	}

	if sub != s.expectedUserClaims.Sub {
		return fmt.Errorf(
			"unexpected 'sub' value received in user profile; expected %s got %s",
			s.expectedUserClaims.Sub, sub)
	}

	b, err := json.Marshal(profile["bootstrap"])
	if err != nil {
		return fmt.Errorf("unmarshal bootstrap data: %w", err)
	}

	respData := oidc.BootstrapData{}

	err = json.Unmarshal(b, &respData)
	if err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}

	if respData.AuthzKeyStoreURL == "" {
		return errors.New("authz kms keystore url can't be empty")
	}

	if respData.OpsKeyStoreURL == "" {
		return errors.New("ops kms keystore url can't be empty")
	}

	if respData.OpsEDVVaultURL == "" {
		return errors.New("ops edv vault url can't be empty")
	}

	if respData.UserEDVVaultURL == "" {
		return errors.New("user edv vault url can't be empty")
	}

	if respData.EDVOpsKIDURL == "" {
		return errors.New("ops edv KID url can't be empty")
	}

	if respData.EDVHMACKIDURL == "" {
		return errors.New("hmac edv KID url can't be empty")
	}

	return nil
}

func (s *Steps) userIsLoggedIn() error {
	err := s.userClicksLoginButton()
	if err != nil {
		return err
	}

	err = s.userRedirectedToOIDCProvider()
	if err != nil {
		return err
	}

	err = s.userIsAuthenticated()
	if err != nil {
		return err
	}

	err = s.userAuthorizesAccessToTheirData()
	if err != nil {
		return err
	}

	err = s.userRedirectedToWallet()
	if err != nil {
		return err
	}

	return s.userRetrievesProfile()
}

func (s *Steps) userLogout() error {
	response, err := s.browser.Get(userLogoutPath) // nolint:noctx // ignoring rule since this is a bdd test
	if err != nil {
		return fmt.Errorf("user failed to logout: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"unexpected response code when logging out; expected %d got %d",
			http.StatusOK, response.StatusCode,
		)
	}

	err = response.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}

	return nil
}

func (s *Steps) userCannotAccessProfile() error {
	response, err := s.browser.Get(userProfilePath) // nolint:noctx // ignoring rule since this is a bdd test
	if err != nil {
		return fmt.Errorf("user failed to invoke user profile endpoint: %w", err)
	}

	if response.StatusCode != http.StatusForbidden {
		return fmt.Errorf(
			"unexpected response code when logging out; expected %d got %d",
			http.StatusForbidden, response.StatusCode,
		)
	}

	err = response.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close body: %w", err)
	}

	return nil
}

func (s *Steps) loginDevice() error {
	resp, err := s.browser.Get(deviceloginBeginPath) // nolint:noctx // no need to set context
	if err != nil {
		return fmt.Errorf("failed to invoke http server login endpoint %s: %w", deviceloginBeginPath, err)
	}

	defer func() {
		_ = resp.Body.Close() // nolint:errcheck // no need
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"unexpected status code at /device/login/begin; expected %d got %d, body: %s",
			http.StatusOK, resp.StatusCode, string(data),
		)
	}

	var credAssert protocol.CredentialAssertion

	err = json.Unmarshal(data, &credAssert)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cred request: %v", err)
	}

	authBytes, err := s.mockAuthenticator.Assert(host, &credAssert.Response)
	if err != nil {
		return fmt.Errorf("failed to generate credential assertion: %v", err)
	}

	resp2, err := s.browser.Post(deviceloginFinishPath, "application/json", bytes.NewReader(authBytes)) // nolint:noctx,lll // don't need in bdd test
	if err != nil {
		return fmt.Errorf("failed to invoke http server device finish login endpoint %s: %w", registerBeginPath, err)
	}

	defer func() {
		_ = resp2.Body.Close() // nolint:errcheck // no need
	}()

	data, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"unexpected status code at /device/login/finish; expected %d got %d, body: %s",
			http.StatusOK, resp.StatusCode, string(data),
		)
	}

	return nil
}

func (s *Steps) registerDevice() error {
	// first we call registerBeginPath
	resp, err := s.browser.Get(registerBeginPath) // nolint:noctx // no need to set context
	if err != nil {
		return fmt.Errorf("failed to invoke http server register device endpoint %s: %w", registerBeginPath, err)
	}

	defer func() {
		_ = resp.Body.Close() // nolint:errcheck // no need
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var credParams protocol.CredentialCreation

	err = json.Unmarshal(data, &credParams)
	if err != nil {
		return fmt.Errorf("failed to parse credential creation parameters: %w", err)
	}

	authData, err := s.mockAuthenticator.Authenticate(
		"fixtures/keys/device/ec-cacert.pem",
		"fixtures/keys/device/ec-cakey.pem",
		host, &credParams)
	if err != nil {
		return fmt.Errorf("failed to generate mock authenticator response: %w", err)
	}

	resp2, err := s.browser.Post(registerFinishPath, "application/json", bytes.NewReader(authData)) // nolint:noctx,lll // don't need in bdd test
	if err != nil {
		return fmt.Errorf("failed to invoke http server register device endpoint %s: %w", registerBeginPath, err)
	}

	defer func() {
		_ = resp2.Body.Close() // nolint:errcheck // no need
	}()

	data, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"unexpected status code at /device/register/finish; expected %d got %d, body: %s",
			http.StatusOK, resp.StatusCode, string(data),
		)
	}

	var cred webauthn.Credential

	err = json.Unmarshal(data, &cred)
	if err != nil || cred.ID == nil {
		return fmt.Errorf("failed to parse credential: %w", err)
	}

	return nil
}
