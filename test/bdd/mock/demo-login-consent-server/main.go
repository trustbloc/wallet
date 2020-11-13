/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	tlsutils "github.com/trustbloc/edge-core/pkg/utils/tls"
)

const (
	adminURLEnvKey          = "ADMIN_URL"
	servePortEnvKey         = "SERVE_PORT"
	tlsSystemCertPoolEnvKey = "TLS_SYSTEMCERTPOOL"
	tlsCACertsEnvKey        = "TLS_CACERTS"

	loginHTML           = "./templates/login.html"
	consentHTML         = "./templates/consent.html"
	bankloginHTML       = "./templates/banklogin.html"
	bankconsentHTML     = "./templates/bankconsent.html"
	dlUploadHTML        = "./templates/uploadCred.html"
	dlUploadConsentHTML = "./templates/uploadCredConsent.html"

	bankLogin = "banklogin"
	dlUpload  = "uploaddrivinglicense"

	bankFlow     = "bank"
	dlUploadFlow = "dlUpload"
	defaultFlow  = "default"

	loginTypeCookie = "loginType"

	timeout = 10 * time.Second
)

type htmlTemplate interface {
	Execute(wr io.Writer, data interface{}) error
}

func main() {
	c, err := buildConsentServer()
	if err != nil {
		panic(err)
	}

	port := os.Getenv(servePortEnvKey)
	if port == "" {
		panic("port to be served not provided")
	}

	// Hydra login and consent handlers
	http.HandleFunc("/login", c.login)
	http.HandleFunc("/consent", c.consent)

	http.Handle("/img/", http.FileServer(http.Dir("templates")))

	fmt.Println(http.ListenAndServe(":"+port, nil))
}

func buildConsentServer() (*consentServer, error) {
	adminURL := os.Getenv(adminURLEnvKey)
	if adminURL == "" {
		return nil, fmt.Errorf("admin URL is required")
	}

	var tlsSystemCertPool bool

	tlsSystemCertPoolVal := os.Getenv(tlsSystemCertPoolEnvKey)
	if tlsSystemCertPoolVal != "" {
		var err error

		tlsSystemCertPool, err = strconv.ParseBool(tlsSystemCertPoolVal)
		if err != nil {
			return nil, fmt.Errorf("invalid value (%s) suppiled for `%s`, switching to default false",
				tlsSystemCertPoolVal, tlsSystemCertPoolEnvKey)
		}
	}

	var tlsCACerts []string

	tlsCACertsVal := os.Getenv(tlsCACertsEnvKey)
	if tlsCACertsVal != "" {
		tlsCACerts = strings.Split(tlsCACertsVal, ",")
	}

	return newConsentServer(adminURL, tlsSystemCertPool, tlsCACerts)
}

// newConsentServer returns new login consent server instance
func newConsentServer(adminURL string, tlsSystemCertPool bool, tlsCACerts []string) (*consentServer, error) {
	u, err := url.Parse(adminURL)
	if err != nil {
		return nil, err
	}

	loginTemplate, err := template.ParseFiles(loginHTML)
	if err != nil {
		return nil, err
	}

	consentTemplate, err := template.ParseFiles(consentHTML)
	if err != nil {
		return nil, err
	}

	bankLoginTemplate, err := template.ParseFiles(bankloginHTML)
	if err != nil {
		return nil, err
	}

	bankConsentTemplate, err := template.ParseFiles(bankconsentHTML)
	if err != nil {
		return nil, err
	}

	dlUploadTemplate, err := template.ParseFiles(dlUploadHTML)
	if err != nil {
		return nil, err
	}

	dlUploadConsentTemplate, err := template.ParseFiles(dlUploadConsentHTML)
	if err != nil {
		return nil, err
	}

	rootCAs, err := tlsutils.GetCertPool(tlsSystemCertPool, tlsCACerts)
	if err != nil {
		return nil, err
	}

	return &consentServer{
		hydraClient: client.NewHTTPClientWithConfig(nil,
			&client.TransportConfig{Schemes: []string{u.Scheme}, Host: u.Host, BasePath: u.Path}),
		loginTemplate:           loginTemplate,
		consentTemplate:         consentTemplate,
		bankLoginTemplate:       bankLoginTemplate,
		bankConsentTemplate:     bankConsentTemplate,
		dlUploadTemplate:        dlUploadTemplate,
		dlUploadConsentTemplate: dlUploadConsentTemplate,
		httpClient: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: rootCAs, MinVersion: tls.VersionTLS12}}},
	}, nil
}

// ConsentServer hydra login consent server
type consentServer struct {
	hydraClient             *client.OryHydra
	loginTemplate           htmlTemplate
	consentTemplate         htmlTemplate
	bankLoginTemplate       htmlTemplate
	bankConsentTemplate     htmlTemplate
	dlUploadTemplate        htmlTemplate
	dlUploadConsentTemplate htmlTemplate
	httpClient              *http.Client
}

func (c *consentServer) login(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		challenge := req.URL.Query().Get("login_challenge")
		fullData := map[string]interface{}{
			"login_challenge": challenge,
		}

		switch {
		case strings.Contains(req.Referer(), bankLogin):
			expire := time.Now().AddDate(0, 0, 1)
			cookie := http.Cookie{Name: loginTypeCookie, Value: bankFlow, Expires: expire}
			http.SetCookie(w, &cookie)

			err := c.bankLoginTemplate.Execute(w, fullData)
			if err != nil {
				fmt.Fprint(w, err.Error())
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		case strings.Contains(req.Referer(), dlUpload):
			expire := time.Now().AddDate(0, 0, 1)
			cookie := http.Cookie{Name: loginTypeCookie, Value: dlUploadFlow, Expires: expire}
			http.SetCookie(w, &cookie)

			err := c.dlUploadTemplate.Execute(w, fullData)
			if err != nil {
				fmt.Fprint(w, err.Error())
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		default:
			expire := time.Now().AddDate(0, 0, 1)
			cookie := http.Cookie{Name: loginTypeCookie, Value: defaultFlow, Expires: expire}
			http.SetCookie(w, &cookie)

			err := c.loginTemplate.Execute(w, fullData)
			if err != nil {
				fmt.Fprint(w, err.Error())
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
		}

	case http.MethodPost:
		c.acceptLoginRequest(w, req)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (c *consentServer) consent(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		c.showConsentPage(w, req)
		return
	case "POST":
		ok := parseRequestForm(w, req)
		if !ok {
			return
		}

		allowed, found := req.Form["submit"]
		if !found {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "consent value missing, Bad request!")

			return
		}

		switch allowed[0] {
		case "accept":
			c.acceptConsentRequest(w, req)
			return
		case "reject":
			c.rejectConsentRequest(w, req)
			return
		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "incorrect consent value, Bad request!")

			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (c *consentServer) acceptLoginRequest(w http.ResponseWriter, req *http.Request) {
	ok := parseRequestForm(w, req)
	if !ok {
		return
	}

	username, usernameSet := req.Form["email"]
	password, passwordSet := req.Form["password"]
	challenge, challengeSet := req.Form["challenge"]

	if !usernameSet || !passwordSet || !challengeSet || !c.authLogin(username[0], password[0]) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	loginRqstParams := admin.NewGetLoginRequestParamsWithHTTPClient(c.httpClient)
	loginRqstParams.SetTimeout(timeout)
	loginRqstParams.LoginChallenge = challenge[0]

	resp, err := c.hydraClient.Admin.GetLoginRequest(loginRqstParams)
	if err != nil {
		fmt.Fprint(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	loginOKRequest := admin.NewAcceptLoginRequestParamsWithHTTPClient(c.httpClient)

	b := &models.AcceptLoginRequest{
		Subject: &username[0],
	}

	loginOKRequest.SetBody(b)
	loginOKRequest.SetTimeout(timeout)
	loginOKRequest.LoginChallenge = resp.Payload.Challenge

	loginOKResponse, err := c.hydraClient.Admin.AcceptLoginRequest(loginOKRequest)
	if err != nil {
		fmt.Fprint(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.Redirect(w, req, loginOKResponse.Payload.RedirectTo, http.StatusFound)
}

func (c *consentServer) showConsentPage(w http.ResponseWriter, req *http.Request) { // nolint: gocyclo
	// get the consent request
	consentRqstParams := admin.NewGetConsentRequestParamsWithHTTPClient(c.httpClient)
	consentRqstParams.SetTimeout(timeout)
	consentRqstParams.ConsentChallenge = req.URL.Query().Get("consent_challenge")

	consentRequest, err := c.hydraClient.Admin.GetConsentRequest(consentRqstParams)
	if err != nil {
		fmt.Fprint(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	fullData := map[string]interface{}{
		"User":      consentRequest.Payload.Subject,
		"Challenge": consentRqstParams.ConsentChallenge,
		"Scope":     consentRequest.Payload.RequestedScope,
	}

	if consentRequest.Payload.Client != nil {
		fullData["ClientName"] = consentRequest.Payload.Client.ClientName
		fullData["ClientID"] = consentRequest.Payload.Client.ClientID
	}

	loginTypeCookie, err := req.Cookie(loginTypeCookie)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		fmt.Fprint(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	switch loginTypeCookie.Value {
	case bankFlow:
		err = c.bankConsentTemplate.Execute(w, fullData)
		if err != nil {
			fmt.Fprint(w, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	case dlUploadFlow:
		err = c.dlUploadConsentTemplate.Execute(w, fullData)
		if err != nil {
			fmt.Fprint(w, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		err = c.consentTemplate.Execute(w, fullData)
		if err != nil {
			fmt.Fprint(w, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (c *consentServer) acceptConsentRequest(w http.ResponseWriter, req *http.Request) {
	getConsentRequest := admin.NewGetConsentRequestParamsWithHTTPClient(c.httpClient)
	getConsentRequest.SetTimeout(timeout)
	getConsentRequest.ConsentChallenge = req.URL.Query().Get("consent_challenge")

	getConsentRequestResponse, err := c.hydraClient.Admin.GetConsentRequest(getConsentRequest)
	if err != nil {
		fmt.Fprint(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, remember := req.Form["remember"]
	b := &models.AcceptConsentRequest{
		GrantScope:               req.Form["grant_scope"],
		GrantAccessTokenAudience: getConsentRequestResponse.Payload.RequestedAccessTokenAudience,
		Remember:                 remember,
		HandledAt:                strfmt.DateTime(time.Now()),
	}

	consentOKRequest := admin.NewAcceptConsentRequestParamsWithHTTPClient(c.httpClient)
	consentOKRequest.SetBody(b)
	consentOKRequest.SetTimeout(timeout)
	consentOKRequest.ConsentChallenge = req.URL.Query().Get("consent_challenge")

	consentOKResponse, err := c.hydraClient.Admin.AcceptConsentRequest(consentOKRequest)
	if err != nil {
		fmt.Fprint(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.Redirect(w, req, consentOKResponse.Payload.RedirectTo, http.StatusFound)
}

func (c *consentServer) rejectConsentRequest(w http.ResponseWriter, req *http.Request) {
	consentDeniedRequest := admin.NewRejectConsentRequestParamsWithHTTPClient(c.httpClient)

	b := &models.RejectRequest{
		Error:            "access_denied",
		ErrorDescription: "The resource owner denied the request",
	}

	consentDeniedRequest.SetBody(b)
	consentDeniedRequest.SetTimeout(timeout)
	consentDeniedRequest.ConsentChallenge = req.URL.Query().Get("consent_challenge")

	consentDenyResponse, err := c.hydraClient.Admin.RejectConsentRequest(consentDeniedRequest)
	if err != nil {
		fmt.Fprint(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.Redirect(w, req, consentDenyResponse.Payload.RedirectTo, http.StatusFound)
}

// authLogin authenticates user login credentials,
// currently authenticating all users
func (c *consentServer) authLogin(usr, pwd string) bool {
	return true
}

// parseRequestForm parses request form.
// writes error to response and returns false when failed.
func parseRequestForm(w http.ResponseWriter, req *http.Request) bool {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprint(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)

		return false
	}

	return true
}
