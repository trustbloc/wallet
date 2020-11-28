/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package startcmd // nolint:testpackage // using private types in tests

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type mockServer struct {
	Err error
}

func (s *mockServer) ListenAndServe(host, certFile, keyFile string, handler http.Handler) error {
	return s.Err
}

func TestListenAndServe(t *testing.T) {
	router, err := router(&httpServerParameters{
		oidc: &oidcParameters{providerURL: mockOIDCProvider(t)},
		tls:  &tlsParameters{},
		keys: &keyParameters{},
		webAuth: &webauthParameters{
			rpDisplayName: "Foobar Corp.",
			rpID:          "localhost",
			rpOrigin:      "http://localhost",
		},
		keyServer: &keyServerParameters{
			authzKMSURL: "http://localhost",
		},
	})
	require.NoError(t, err)

	h := HTTPServer{}

	err = h.ListenAndServe("localhost:8080", "test.key", "test.cert", router)
	require.Error(t, err)
	require.Contains(t, err.Error(), "open test.key: no such file or directory")
}

func TestStartCmdContents(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	require.Equal(t, "start", startCmd.Use)
	require.Equal(t, "Start http server", startCmd.Short)
	require.Equal(t, "Start http server", startCmd.Long)

	checkFlagPropertiesCorrect(t, startCmd, hostURLFlagName, hostURLFlagShorthand, hostURLFlagUsage)
}

func TestStartCmdWithBlankArg(t *testing.T) {
	t.Run("test blank host arg", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "",
			"--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Equal(t, "host-url value is empty", err.Error())
	})

	t.Run("test blank tls cert arg", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080",
			"--" + agentUIURLFlagName, "ui",
			"--" + tlsCertFileFlagName, "",
			"--" + tlsKeyFileFlagName, "key",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Equal(t, "failed to configure tls cert file: tls-cert-file value is empty", err.Error())
	})

	t.Run("test blank tls cert arg", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080",
			"--" + agentUIURLFlagName, "ui",
			"--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Equal(t, "failed to configure tls key file: tls-key-file value is empty", err.Error())
	})
}

func TestStartCmdWithMissingArg(t *testing.T) {
	t.Run("test missing host arg", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		err := startCmd.Execute()
		require.Error(t, err)
		require.Equal(t,
			"Neither host-url (command line flag) nor HTTP_SERVER_HOST_URL (environment variable) have been set.",
			err.Error())
	})

	t.Run("test invalid auto accept flag", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{Err: errors.New("error starting the server")})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + sessionCookieAuthKeyFlagName, key(t),
			"--" + sessionCookieEncKeyFlagName, key(t),
			"--" + webAuthRPDisplayFlagName, "Foobar Corp.",
			"--" + webAuthRPIDFlagName, "localhost",
			"--" + webAuthRPOriginFlagName, "http://localhost",
			"--" + authzKMSURLFlagName, "http://localhost",
			"--" + opsKMSURLFlagName, "http://localhost",
			"--" + keyEDVURLFlagName, "http://localhost",
			"--" + hubAuthURLFlagName, "http://localhost",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "error starting the server")
	})

	t.Run("test invalid tls-cacerts", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, "INVALID",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()

		require.EqualError(t, err,
			"failed to init tls cert pool: failed to read cert: open INVALID: no such file or directory")
	})

	t.Run("missing oidc provider URL", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.EqualError(t, err,
			"failed to configure OIDC provider URL: Neither oidc-opurl (command line flag) nor"+
				" HTTP_SERVER_OIDC_OPURL (environment variable) have been set.")
	})

	t.Run("invalid oidc provider URL", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + dependencyMaxRetriesFlagName, "1",
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, "INVALID",
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieAuthKeyFlagName, key(t),
			"--" + sessionCookieEncKeyFlagName, key(t),
			"--" + webAuthRPDisplayFlagName, "Foobar Corp.",
			"--" + webAuthRPIDFlagName, "localhost",
			"--" + webAuthRPOriginFlagName, "http://localhost",
			"--" + authzKMSURLFlagName, "http://localhost",
			"--" + opsKMSURLFlagName, "http://localhost",
			"--" + keyEDVURLFlagName, "http://localhost",
			"--" + hubAuthURLFlagName, "http://localhost",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to init OIDC provider")
	})

	t.Run("missing oidc client ID", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + webAuthRPDisplayFlagName, "Foobar Corp.",
			"--" + webAuthRPIDFlagName, "localhost",
			"--" + webAuthRPOriginFlagName, "http://localhost",
			"--" + authzKMSURLFlagName, "http://localhost",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.EqualError(t, err,
			"failed to configure OIDC clientID: Neither oidc-clientid (command line flag) nor"+
				" HTTP_SERVER_OIDC_CLIENTID (environment variable) have been set.")
	})

	t.Run("missing oidc client secret", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.EqualError(t, err,
			"failed to configure OIDC client secret: Neither oidc-clientsecret (command line flag) nor"+
				" HTTP_SERVER_OIDC_CLIENTSECRET (environment variable) have been set.")
	})

	t.Run("missing oidc callback", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + tlsCACertsFlagName, cert(t),
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.EqualError(t, err,
			"failed to configure OIDC callback URL: Neither oidc-callback (command line flag) nor"+
				" HTTP_SERVER_OIDC_CALLBACK (environment variable) have been set.")
	})

	t.Run("missing session cookie auth key", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + dependencyMaxRetriesFlagName, "1",
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, "INVALID",
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieEncKeyFlagName, key(t),
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(),
			"failed to configure session cookie auth key: Neither cookie-auth-key (command line flag) nor"+
				" HTTP_SERVER_COOKIE_AUTH_KEY (environment variable) have been set.")
	})

	t.Run("invalid session cookie auth key path", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + dependencyMaxRetriesFlagName, "1",
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, "INVALID",
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieAuthKeyFlagName, "INVALID",
			"--" + sessionCookieEncKeyFlagName, key(t),
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(),
			"failed to configure session cookie auth key: failed to read file INVALID: open INVALID:"+
				" no such file or directory")
	})

	t.Run("invalid session cookie auth key length", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + dependencyMaxRetriesFlagName, "1",
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, "INVALID",
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieAuthKeyFlagName, invalidKey(t),
			"--" + sessionCookieEncKeyFlagName, key(t),
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to configure session cookie auth key")
	})

	t.Run("missing session cookie enc key", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + dependencyMaxRetriesFlagName, "1",
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, "INVALID",
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieAuthKeyFlagName, key(t),
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(),
			"failed to configure session cookie enc key: Neither cookie-enc-key (command line flag) nor"+
				" HTTP_SERVER_COOKIE_ENC_KEY (environment variable) have been set.")
	})

	t.Run("invalid log level", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieAuthKeyFlagName, key(t),
			"--" + sessionCookieEncKeyFlagName, key(t),
			"--" + agentLogLevelFlagName, "INVALID",
			"--" + authzKMSURLFlagName, "http://localhost",
			"--" + opsKMSURLFlagName, "http://localhost",
			"--" + keyEDVURLFlagName, "http://localhost",
			"--" + hubAuthURLFlagName, "http://localhost",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()

		require.Error(t, err)
		require.Contains(t, err.Error(),
			"failed to set log level: failed to parse log level 'INVALID' : logger: invalid log level")
	})

	t.Run("missing authz key server url", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieAuthKeyFlagName, key(t),
			"--" + sessionCookieEncKeyFlagName, key(t),
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()

		require.Error(t, err)
		require.Contains(t, err.Error(),
			"Neither authz-kms-url (command line flag) nor HTTP_SERVER_AUTHZ_KMS_URL (environment variable) have been set.")
	})

	t.Run("missing ops edv server url", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieAuthKeyFlagName, key(t),
			"--" + sessionCookieEncKeyFlagName, key(t),
			"--" + authzKMSURLFlagName, "http://localhost",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()

		require.Error(t, err)
		require.Contains(t, err.Error(),
			"Neither key-edv-url (command line flag) nor HTTP_SERVER_KEY_EDV_URL (environment variable) have been set.")
	})

	t.Run("missing ops key server url", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieAuthKeyFlagName, key(t),
			"--" + sessionCookieEncKeyFlagName, key(t),
			"--" + authzKMSURLFlagName, "http://localhost",
			"--" + keyEDVURLFlagName, "http://localhost",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()

		require.Error(t, err)
		require.Contains(t, err.Error(),
			"Neither ops-kms-url (command line flag) nor HTTP_SERVER_OPS_KMS_URL (environment variable) have been set.")
	})

	t.Run("missing hub-auth server url", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieAuthKeyFlagName, key(t),
			"--" + sessionCookieEncKeyFlagName, key(t),
			"--" + authzKMSURLFlagName, "http://localhost",
			"--" + keyEDVURLFlagName, "http://localhost",
			"--" + opsKMSURLFlagName, "http://localhost",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()

		require.Error(t, err)
		require.Contains(t, err.Error(),
			"Neither hub-auth-url (command line flag) nor HTTP_SERVER_HUB_AUTH_URL (environment variable) have been set.")
	})

	t.Run("invalid use Remote KMS value - must be t/f", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieAuthKeyFlagName, key(t),
			"--" + sessionCookieEncKeyFlagName, key(t),
			"--" + authzKMSURLFlagName, "http://localhost",
			"--" + keyEDVURLFlagName, "http://localhost",
			"--" + opsKMSURLFlagName, "http://localhost",
			"--" + hubAuthURLFlagName, "http://localhost",
			"--" + useRemoteKMSFlagName, "xyz",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()

		require.EqualError(t, err, "use remote kms param to bool : strconv.ParseBool: parsing \"xyz\": "+
			"invalid syntax")
	})

	t.Run("empty use Remote KMS value - must be t/f", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentUIURLFlagName, "ui",
			"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
			"--" + oidcClientIDFlagName, uuid.New().String(),
			"--" + oidcClientSecretFlagName, uuid.New().String(),
			"--" + oidcCallbackURLFlagName, "http://test.com/callback",
			"--" + tlsCACertsFlagName, cert(t),
			"--" + sessionCookieAuthKeyFlagName, key(t),
			"--" + sessionCookieEncKeyFlagName, key(t),
			"--" + authzKMSURLFlagName, "http://localhost",
			"--" + keyEDVURLFlagName, "http://localhost",
			"--" + opsKMSURLFlagName, "http://localhost",
			"--" + hubAuthURLFlagName, "http://localhost",
			"--" + useRemoteKMSFlagName, "",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()

		require.EqualError(t, err, "get use remote kms param : use-remote-kms value is empty")
	})
}

func TestStartCmdValidArgs(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	args := []string{
		"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
		"--" + tlsKeyFileFlagName, "key",
		"--" + agentUIURLFlagName, "ui",
		"--" + oidcProviderURLFlagName, mockOIDCProvider(t),
		"--" + oidcClientIDFlagName, uuid.New().String(),
		"--" + oidcClientSecretFlagName, uuid.New().String(),
		"--" + oidcCallbackURLFlagName, "http://test.com/callback",
		"--" + tlsCACertsFlagName, cert(t),
		"--" + sessionCookieAuthKeyFlagName, key(t),
		"--" + sessionCookieEncKeyFlagName, key(t),
		"--" + webAuthRPDisplayFlagName, "Foobar Corp.",
		"--" + webAuthRPIDFlagName, "localhost",
		"--" + webAuthRPOriginFlagName, "http://localhost",
		"--" + authzKMSURLFlagName, "http://localhost",
		"--" + opsKMSURLFlagName, "http://localhost",
		"--" + keyEDVURLFlagName, "http://localhost",
		"--" + hubAuthURLFlagName, "http://localhost",
	}
	startCmd.SetArgs(args)

	err := startCmd.Execute()

	require.NoError(t, err)
}

func TestStartCmdValidArgsEnvVar(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	err := os.Setenv(hostURLEnvKey, "localhost:8080")
	require.NoError(t, err)

	err = os.Setenv(tlsCertFileEnvKey, "cert")
	require.NoError(t, err)

	err = os.Setenv(tlsKeyFileEnvKey, "key")
	require.NoError(t, err)

	err = os.Setenv(agentUIURLEnvKey, "ui")
	require.NoError(t, err)

	err = os.Setenv(oidcProviderURLEnvKey, mockOIDCProvider(t))
	require.NoError(t, err)

	err = os.Setenv(oidcClientIDEnvKey, uuid.New().String())
	require.NoError(t, err)

	err = os.Setenv(oidcClientSecretEnvKey, uuid.New().String())
	require.NoError(t, err)

	err = os.Setenv(oidcCallbackURLEnvKey, "http://test.com/callback")
	require.NoError(t, err)

	err = os.Setenv(sessionCookieEncKeyEnvKey, key(t))
	require.NoError(t, err)

	err = os.Setenv(sessionCookieAuthKeyEnvKey, key(t))
	require.NoError(t, err)

	err = os.Setenv(webAuthRPDisplayEnvKey, "Foobar Corp.")
	require.NoError(t, err)

	err = os.Setenv(webAuthRPIDEnvKey, "http://localhost")
	require.NoError(t, err)

	err = os.Setenv(webAuthRPOriginEnvKey, "localhost")
	require.NoError(t, err)

	err = os.Setenv(authzKMSURLEnvKey, "localhost")
	require.NoError(t, err)

	err = os.Setenv(opsKMSURLEnvKey, "localhost")
	require.NoError(t, err)

	err = os.Setenv(keyEDVURLEnvKey, "localhost")
	require.NoError(t, err)

	err = os.Setenv(hubAuthURLEnvKey, "localhost")
	require.NoError(t, err)

	err = startCmd.Execute()

	require.NoError(t, err)
}

func TestStartCmdWithBlankEnvVar(t *testing.T) {
	t.Run("test blank host env var", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		err := os.Setenv(hostURLEnvKey, "")
		require.NoError(t, err)

		err = os.Setenv(tlsCertFileEnvKey, "cert")
		require.NoError(t, err)

		err = os.Setenv(tlsKeyFileEnvKey, "key")
		require.NoError(t, err)

		err = startCmd.Execute()
		require.Error(t, err)
		require.Equal(t, "HTTP_SERVER_HOST_URL value is empty", err.Error())
	})
}

func TestHealthCheckHandler(t *testing.T) {
	result := httptest.NewRecorder()
	healthCheckHandler(result, nil)
	require.Equal(t, http.StatusOK, result.Code)
}

func checkFlagPropertiesCorrect(t *testing.T, cmd *cobra.Command, flagName, flagShorthand, flagUsage string) {
	flag := cmd.Flag(flagName)

	require.NotNil(t, flag)
	require.Equal(t, flagName, flag.Name)
	require.Equal(t, flagShorthand, flag.Shorthand)
	require.Equal(t, flagUsage, flag.Usage)
	require.Equal(t, "", flag.Value.String())

	flagAnnotations := flag.Annotations
	require.Nil(t, flagAnnotations)
}

func mockOIDCProvider(t *testing.T) string {
	h := &testOIDCProvider{}
	srv := httptest.NewServer(h)
	h.baseURL = srv.URL

	t.Cleanup(srv.Close)

	return srv.URL
}

type oidcConfigJSON struct {
	Issuer      string   `json:"issuer"`
	AuthURL     string   `json:"authorization_endpoint"`
	TokenURL    string   `json:"token_endpoint"`
	JWKSURL     string   `json:"jwks_uri"`
	UserInfoURL string   `json:"userinfo_endpoint"`
	Algorithms  []string `json:"id_token_signing_alg_values_supported"`
}

type testOIDCProvider struct {
	baseURL string
}

func (t *testOIDCProvider) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	response, err := json.Marshal(&oidcConfigJSON{
		Issuer:      t.baseURL,
		AuthURL:     fmt.Sprintf("%s/oauth2/auth", t.baseURL),
		TokenURL:    fmt.Sprintf("%s/oauth2/token", t.baseURL),
		JWKSURL:     fmt.Sprintf("%s/oauth2/certs", t.baseURL),
		UserInfoURL: fmt.Sprintf("%s/oauth2/userinfo", t.baseURL),
		Algorithms:  []string{"RS256"},
	})
	if err != nil {
		panic(err)
	}

	_, err = w.Write(response)
	if err != nil {
		panic(err)
	}
}

func cert(t *testing.T) string {
	file, err := ioutil.TempFile("", "*.pem")
	require.NoError(t, err)

	t.Cleanup(func() {
		fileErr := file.Close()
		require.NoError(t, fileErr)
		fileErr = os.Remove(file.Name())
		require.NoError(t, fileErr)
	})

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	secret, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &secret.PublicKey, secret)
	require.NoError(t, err)

	err = pem.Encode(file, &pem.Block{Type: "CERTIFICATE", Bytes: der})
	require.NoError(t, err)

	return file.Name()
}

func key(t *testing.T) string {
	t.Helper()

	key := make([]byte, 32)

	n, err := rand.Reader.Read(key)
	require.NoError(t, err)
	require.Equal(t, 32, n)

	file, err := ioutil.TempFile("", "test_*.key")
	require.NoError(t, err)

	t.Cleanup(func() {
		delErr := os.Remove(file.Name())
		require.NoError(t, delErr)
	})

	err = ioutil.WriteFile(file.Name(), key, os.ModeAppend)
	require.NoError(t, err)

	return file.Name()
}

func invalidKey(t *testing.T) string {
	t.Helper()

	key := make([]byte, 18)

	n, err := rand.Reader.Read(key)
	require.NoError(t, err)
	require.Equal(t, 18, n)

	file, err := ioutil.TempFile("", "test_*.key")
	require.NoError(t, err)

	t.Cleanup(func() {
		delErr := os.Remove(file.Name())
		require.NoError(t, delErr)
	})

	err = ioutil.WriteFile(file.Name(), key, os.ModeAppend)
	require.NoError(t, err)

	return file.Name()
}
