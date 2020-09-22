/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package startcmd // nolint:testpackage // using private types in tests

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	cmdutils "github.com/trustbloc/edge-core/pkg/utils/cmd"
)

type mockServer struct {
	Err error
}

func (s *mockServer) ListenAndServe(host, certFile, keyFile string, handler http.Handler) error {
	return s.Err
}

type mockHTTPResponseWriter struct{}

func (m *mockHTTPResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockHTTPResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m *mockHTTPResponseWriter) WriteHeader(_ int) {}

func TestVueHandler(t *testing.T) {
	h := VueHandler("", &ariesJSOpts{}, &trustblocAgentJSOpts{})
	require.NotNil(t, h)
	h.ServeHTTP(&mockHTTPResponseWriter{}, &http.Request{URL: &url.URL{}})
	h.ServeHTTP(&mockHTTPResponseWriter{}, &http.Request{URL: &url.URL{Path: "."}})
	h.ServeHTTP(&mockHTTPResponseWriter{}, &http.Request{URL: &url.URL{Path: "/aries/jsopts"}})
}

func TestListenAndServe(t *testing.T) {
	h := HTTPServer{}
	err := h.ListenAndServe("localhost:8080", "test.key", "test.cert",
		VueHandler("", &ariesJSOpts{}, &trustblocAgentJSOpts{}))
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
			"--" + tlsCertFileFlagName, "",
			"--" + tlsKeyFileFlagName, "key",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Equal(t, "tls-cert-file value is empty", err.Error())
	})

	t.Run("test blank tls cert arg", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080",
			"--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Equal(t, "tls-key-file value is empty", err.Error())
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

	t.Run("test missing bloc domain arg", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080",
			"--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Equal(t,
			"Neither bloc-domain (command line flag) nor BLOC_DOMAIN (environment variable) have been set.",
			err.Error())
	})

	t.Run("test invalid auto accept flag", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080",
			"--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + agentAutoAcceptFlagName, "invalid",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Equal(t,
			"invalid option - set true or false as the value",
			err.Error())
	})

	t.Run("test invalid auto accept flag", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{Err: errors.New("error starting the server")})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + blocDomainFlagName, "domain",
			"--" + agentAutoAcceptFlagName, "false",
			"--" + agentHTTPResolverFlagName, "sidetree@http://localhost:8901",
			"--" + sdsURLFlagName, "someURL",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "error starting the server")
	})

	t.Run("test invalid blinded routing flag", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{Err: errors.New("error starting the server")})

		args := []string{
			"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
			"--" + tlsKeyFileFlagName, "key",
			"--" + blocDomainFlagName, "domain",
			"--" + blindedRoutingFlagName, "invalid",
			"--" + agentHTTPResolverFlagName, "sidetree@http://localhost:8901",
			"--" + sdsURLFlagName, "someURL",
		}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid option - set true or false as the value")
	})
}

func TestStartCmdValidArgs(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	args := []string{
		"--" + hostURLFlagName, "localhost:8080", "--" + tlsCertFileFlagName, "cert",
		"--" + tlsKeyFileFlagName, "key",
		"--" + blocDomainFlagName, "domain",
		"--" + walletMediatorURLFlagName, "http://localhost:8999",
		"--" + credentialMediatorURLFlagName, "http://auth.sample/mediator",
		"--" + blindedRoutingFlagName, "true",
		"--" + agentAutoAcceptFlagName, "false",
		"--" + agentHTTPResolverFlagName, "sidetree@http://localhost:8901",
		"--" + sdsURLFlagName, "someURL",
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

	err = os.Setenv(blocDomainEnvKey, "domain")
	require.NoError(t, err)

	err = os.Setenv(sdsURLEnvKey, "someURL")
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

func TestGetUserSetVar(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	t.Run("missing mandatory value", func(t *testing.T) {
		vals, err := cmdutils.GetUserSetVarFromArrayString(startCmd, agentLogLevelFlagName, agentLogLevelEnvKey, false)
		require.Error(t, err)
		require.Equal(t,
			"Neither log-level (command line flag) nor ARIESD_LOG_LEVEL (environment variable) have been set.",
			err.Error())
		require.Empty(t, vals)
	})

	t.Run("valid env value", func(t *testing.T) {
		err := os.Setenv(agentLogLevelEnvKey, "sidetree@localhost:8080,uni@localhost:8900")
		require.NoError(t, err)

		vals, err := cmdutils.GetUserSetVarFromArrayString(startCmd, agentLogLevelFlagName, agentLogLevelEnvKey, true)
		require.NoError(t, err)
		require.Equal(t, 2, len(vals))
	})
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
