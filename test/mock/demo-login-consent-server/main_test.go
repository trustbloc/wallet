/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConsent_New(t *testing.T) {
	tests := []struct {
		name     string
		adminURL string
		err      string
	}{
		{
			name:     "initialize with valid admin URL",
			adminURL: "sampleURL",
		},
		{
			name:     "initialize with valid admin URL",
			adminURL: " ?&=#+%!<>#\"{}|\\^[];",
			err:      "invalid URL escape",
		},
	}

	t.Parallel()

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			server, err := newConsentServer(tc.adminURL, false, []string{})
			if tc.err != "" {
				require.Contains(t, err.Error(), tc.err)
				return
			}

			require.NotNil(t, server)
			require.NotNil(t, server.hydraClient)
			require.NotNil(t, server.loginTemplate)
			require.NotNil(t, server.consentTemplate)
		})
	}
}

func TestConsent_buildConsentServer(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		err  string
	}{
		{
			name: "initialize without required ENV variables",
			env:  map[string]string{},
			err:  "admin URL is required",
		},
		{
			name: "initialize with only required ENV variables",
			env: map[string]string{
				adminURLEnvKey: "sampleURL",
			},
		},
		{
			name: "initialize with invalid tls system cert pool ENV variable",
			env: map[string]string{
				adminURLEnvKey:          "sampleURL",
				tlsSystemCertPoolEnvKey: "InVaLid",
			},
			err: "invalid value",
		},
		{
			name: "initialize with invalid tls system cert pool ENV variable",
			env: map[string]string{
				adminURLEnvKey:          "sampleURL",
				tlsSystemCertPoolEnvKey: "false",
				tlsCACertsEnvKey:        "certs",
			},
			err: "failed to read cert",
		},
		{
			name: "initialize with valid ENV variables",
			env: map[string]string{
				adminURLEnvKey:          "sampleURL",
				tlsSystemCertPoolEnvKey: "true",
				tlsCACertsEnvKey:        "",
			},
		},
		{
			name: "initialize with valid ENV variables",
			env: map[string]string{
				adminURLEnvKey:          "sampleURL",
				tlsSystemCertPoolEnvKey: "true",
				tlsCACertsEnvKey:        "",
			},
		},
	}

	t.Parallel()

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				require.NoError(t, os.Setenv(k, v))
			}

			server, err := buildConsentServer()
			if tc.err != "" {
				require.Contains(t, err.Error(), tc.err)
			} else {
				require.NotNil(t, server)
				require.NotNil(t, server.hydraClient)
				require.NotNil(t, server.loginTemplate)
				require.NotNil(t, server.consentTemplate)
			}
		})
	}
}

func TestConsentServer_Login(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/oauth2/auth/requests/login/accept" {
			fmt.Fprint(res, `{"redirect_to":"sampleURL"}`)
		}

		res.WriteHeader(http.StatusOK)
	}))

	defer func() { testServer.Close() }()

	tests := []struct {
		name             string
		adminURL         string
		method           string
		url              string
		form             map[string][]string
		responseHTML     []string
		responseStatus   int
		referer          string
		loginTemplate    htmlTemplate
		bankTemplate     htmlTemplate
		dlUploadTemplate htmlTemplate
		err              string
	}{
		{
			name:           "/login Method not allowed",
			adminURL:       testServer.URL,
			method:         http.MethodPatch,
			url:            "?login_challenge=12345",
			responseStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "/login GET SUCCESS",
			adminURL:       testServer.URL,
			method:         http.MethodGet,
			url:            "?login_challenge=12345",
			responseHTML:   []string{"<title>Login Page</title>", `name="challenge" value="12345"`},
			responseStatus: http.StatusOK,
		},
		{
			name:           "/login GET FAILURE (template error)",
			adminURL:       testServer.URL,
			method:         http.MethodGet,
			url:            "?login_challenge=12345",
			responseStatus: http.StatusOK,
			err:            "template error",
			loginTemplate:  &mockTemplate{executeErr: fmt.Errorf("template error")},
		},
		{
			name:           "/bank login GET SUCCESS",
			adminURL:       testServer.URL,
			method:         http.MethodGet,
			url:            "?login_challenge=12345",
			responseHTML:   []string{"<title>Bank Login Page</title>", `name="challenge" value="12345"`},
			responseStatus: http.StatusOK,
			referer:        bankLogin,
		},
		{
			name:           "/bank login GET FAILURE (template error)",
			adminURL:       testServer.URL,
			method:         http.MethodGet,
			url:            "?login_challenge=12345",
			responseStatus: http.StatusOK,
			err:            "template error",
			bankTemplate:   &mockTemplate{executeErr: fmt.Errorf("template error")},
			referer:        bankLogin,
		},
		{
			name:           "/dlUpload login GET SUCCESS",
			adminURL:       testServer.URL,
			method:         http.MethodGet,
			url:            "?login_challenge=12345",
			responseHTML:   []string{"<title>Upload Credential</title>", `name="challenge" value="12345"`},
			responseStatus: http.StatusOK,
			referer:        dlUpload,
		},
		{
			name:             "/dlUpload login GET FAILURE (template error)",
			adminURL:         testServer.URL,
			method:           http.MethodGet,
			url:              "?login_challenge=12345",
			responseStatus:   http.StatusOK,
			err:              "template error",
			dlUploadTemplate: &mockTemplate{executeErr: fmt.Errorf("template error")},
			referer:          dlUpload,
		},
		{
			name:           "/login POST FAILURE (missing form body)",
			adminURL:       testServer.URL,
			method:         http.MethodPost,
			url:            "?login_challenge=12345",
			err:            "missing form body",
			responseStatus: http.StatusOK,
		},
		{
			name:     "/login POST FAILURE (missing login credentials)",
			adminURL: testServer.URL,
			method:   http.MethodPost,
			form: map[string][]string{
				"email":     {"uname"},
				"challenge": {"12345"},
			},
			responseStatus: http.StatusForbidden,
		},
		{
			name:     "/login POST FAILURE (missing challenge)",
			adminURL: testServer.URL,
			method:   http.MethodPost,
			form: map[string][]string{
				"email":    {"uname"},
				"password": {"pwd"},
			},
			responseStatus: http.StatusForbidden,
		},
		{
			name:     "/login POST SUCCESS",
			adminURL: testServer.URL,
			method:   http.MethodPost,
			form: map[string][]string{
				"email":     {"uname"},
				"password":  {"pwd"},
				"challenge": {"12345"},
			},
			responseStatus: http.StatusOK,
		},
	}

	t.Parallel()

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			server, err := newConsentServer(tc.adminURL, false, []string{})
			require.NotNil(t, server)
			require.NoError(t, err)

			if tc.loginTemplate != nil {
				server.loginTemplate = tc.loginTemplate
			}

			if tc.bankTemplate != nil {
				server.bankLoginTemplate = tc.bankTemplate
			}

			if tc.dlUploadTemplate != nil {
				server.dlUploadTemplate = tc.dlUploadTemplate
			}

			req, err := http.NewRequest(tc.method, tc.url, nil)
			require.NoError(t, err)

			if tc.form != nil {
				req.PostForm = url.Values(tc.form)
			}

			req.Header.Set("Referer", tc.referer)

			res := httptest.NewRecorder()

			server.login(res, req)

			if len(tc.responseHTML) > 0 {
				for _, html := range tc.responseHTML {
					require.Contains(t, res.Body.String(), html)
				}
			}

			if tc.err != "" {
				require.Contains(t, res.Body.String(), tc.err)
			}
			require.Equal(t, tc.responseStatus, res.Code, res.Body.String())
		})
	}
}

func TestConsentServer_Consent(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.RequestURI, "/oauth2/auth/requests/consent/") {
			fmt.Fprint(res, `{"redirect_to":"sampleURL"}`)
		}
		res.WriteHeader(http.StatusOK)
	}))

	defer func() { testServer.Close() }()

	tests := []struct {
		name                    string
		adminURL                string
		method                  string
		url                     string
		form                    map[string][]string
		responseHTML            []string
		responseStatus          int
		cookie                  *http.Cookie
		consentTemplate         htmlTemplate
		bankConsentTemplate     htmlTemplate
		dlUploadConsentTemplate htmlTemplate
		err                     string
	}{
		{
			name:           "/consent Method not allowed",
			adminURL:       testServer.URL,
			method:         http.MethodPatch,
			url:            "?consent_challenge=12345",
			responseStatus: http.StatusMethodNotAllowed,
			cookie:         &http.Cookie{Name: loginTypeCookie, Value: defaultFlow},
		},
		{
			name:           "/consent GET SUCCESS",
			adminURL:       testServer.URL,
			method:         http.MethodGet,
			url:            "?consent_challenge=12345",
			responseHTML:   []string{"<title>Consent Page</title>"},
			responseStatus: http.StatusOK,
			cookie:         &http.Cookie{Name: loginTypeCookie, Value: defaultFlow},
		},
		{
			name:            "/consent GET FAILURE (template error)",
			adminURL:        testServer.URL,
			method:          http.MethodGet,
			url:             "?consent_challenge=12345",
			responseStatus:  http.StatusOK,
			err:             "template error",
			consentTemplate: &mockTemplate{executeErr: fmt.Errorf("template error")},
			cookie:          &http.Cookie{Name: loginTypeCookie, Value: defaultFlow},
		},
		{
			name:           "/bank consent GET SUCCESS",
			adminURL:       testServer.URL,
			method:         http.MethodGet,
			url:            "?consent_challenge=12345",
			responseHTML:   []string{"<title>Bank Consent Page</title>"},
			responseStatus: http.StatusOK,
			cookie:         &http.Cookie{Name: loginTypeCookie, Value: bankFlow},
		},
		{
			name:                "/bank consent GET FAILURE (template error)",
			adminURL:            testServer.URL,
			method:              http.MethodGet,
			url:                 "?consent_challenge=12345",
			responseStatus:      http.StatusOK,
			err:                 "template error",
			bankConsentTemplate: &mockTemplate{executeErr: fmt.Errorf("template error")},
			cookie:              &http.Cookie{Name: loginTypeCookie, Value: bankFlow},
		},
		{
			name:           "/dlUpload consent GET SUCCESS",
			adminURL:       testServer.URL,
			method:         http.MethodGet,
			url:            "?consent_challenge=12345",
			responseHTML:   []string{"<title>Consent Page</title>"},
			responseStatus: http.StatusOK,
			cookie:         &http.Cookie{Name: loginTypeCookie, Value: dlUploadFlow},
		},
		{
			name:                    "/dlUpload consent GET FAILURE (template error)",
			adminURL:                testServer.URL,
			method:                  http.MethodGet,
			url:                     "?consent_challenge=12345",
			responseStatus:          http.StatusOK,
			err:                     "template error",
			dlUploadConsentTemplate: &mockTemplate{executeErr: fmt.Errorf("template error")},
			cookie:                  &http.Cookie{Name: loginTypeCookie, Value: dlUploadFlow},
		},
		{
			name:           "/consent POST FAILURE (missing form body)",
			adminURL:       testServer.URL,
			method:         http.MethodPost,
			err:            "missing form body",
			responseStatus: http.StatusOK,
			cookie:         &http.Cookie{Name: loginTypeCookie, Value: bankFlow},
		},
		{
			name:           "/consent POST FAILURE (missing submit)",
			adminURL:       testServer.URL,
			method:         http.MethodPost,
			form:           map[string][]string{},
			responseStatus: http.StatusBadRequest,
			err:            "consent value missing",
			cookie:         &http.Cookie{Name: loginTypeCookie, Value: defaultFlow},
		},
		{
			name:     "/consent POST FAILURE (invalid submit value)",
			adminURL: testServer.URL,
			method:   http.MethodPost,
			form: map[string][]string{
				"submit": {"xyz"},
			},
			responseStatus: http.StatusBadRequest,
			err:            "incorrect consent value",
			cookie:         &http.Cookie{Name: loginTypeCookie, Value: defaultFlow},
		},
		{
			name:     "/consent POST accept consent value",
			adminURL: testServer.URL,
			method:   http.MethodPost,
			form: map[string][]string{
				"submit": {"accept"},
			},
			responseStatus: http.StatusOK,
			cookie:         &http.Cookie{Name: loginTypeCookie, Value: defaultFlow},
		},
		{
			name:     "/consent POST accept consent value",
			adminURL: testServer.URL,
			method:   http.MethodPost,
			form: map[string][]string{
				"submit": {"reject"},
			},
			responseStatus: http.StatusOK,
			cookie:         &http.Cookie{Name: loginTypeCookie, Value: defaultFlow},
		},
	}

	t.Parallel()

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			server, err := newConsentServer(tc.adminURL, false, []string{})
			require.NotNil(t, server)
			require.NoError(t, err)

			req, err := http.NewRequest(tc.method, tc.url, nil)
			require.NoError(t, err)

			if tc.form != nil {
				req.PostForm = url.Values(tc.form)
			}

			if tc.consentTemplate != nil {
				server.consentTemplate = tc.consentTemplate
			}

			if tc.bankConsentTemplate != nil {
				server.bankConsentTemplate = tc.bankConsentTemplate
			}

			if tc.dlUploadConsentTemplate != nil {
				server.dlUploadConsentTemplate = tc.dlUploadConsentTemplate
			}

			req.AddCookie(tc.cookie)

			res := httptest.NewRecorder()

			server.consent(res, req)

			if len(tc.responseHTML) > 0 {
				for _, html := range tc.responseHTML {
					require.Contains(t, res.Body.String(), html)
				}
			}

			if tc.err != "" {
				require.Contains(t, res.Body.String(), tc.err)
			}
			require.Equal(t, tc.responseStatus, res.Code, res.Body.String())
		})
	}
}

type mockTemplate struct {
	executeErr error
}

func (m *mockTemplate) Execute(wr io.Writer, data interface{}) error {
	return m.executeErr
}
