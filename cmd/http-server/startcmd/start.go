/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package startcmd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	oidcp "github.com/coreos/go-oidc"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/mux"
	"github.com/lpar/gzipped"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	oidc2 "github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"
	"github.com/trustbloc/edge-agent/pkg/restapi/device"
	"github.com/trustbloc/edge-agent/pkg/restapi/oidc"
	"github.com/trustbloc/edge-core/pkg/log"
	"github.com/trustbloc/edge-core/pkg/storage/memstore"
	cmdutils "github.com/trustbloc/edge-core/pkg/utils/cmd"
	tlsutils "github.com/trustbloc/edge-core/pkg/utils/tls"
)

const (
	hostURLFlagName      = "host-url"
	hostURLFlagShorthand = "u"
	hostURLFlagUsage     = "Host Name:Port." +
		" Alternatively, this can be set with the following environment variable: " + hostURLEnvKey
	hostURLEnvKey = "HTTP_SERVER_HOST_URL"

	wasmPathFlagName      = "wasm-path"
	wasmPathFlagShorthand = "w"
	wasmPathFlagUsage     = "WASM path." +
		" Defaults to current directory." +
		" Alternatively, this can be set with the following environment variable: " + wasmPathEnvKey
	wasmPathEnvKey = "HTTP_SERVER_WASM_PATH"

	tlsCertFileFlagName      = "tls-cert-file"
	tlsCertFileFlagShorthand = "c"
	tlsCertFileFlagUsage     = "tls certificate file." +
		" Alternatively, this can be set with the following environment variable: " + tlsCertFileEnvKey
	tlsCertFileEnvKey = "TLS_CERT_FILE"

	tlsKeyFileFlagName      = "tls-key-file"
	tlsKeyFileFlagShorthand = "k"
	tlsKeyFileFlagUsage     = "tls key file." +
		" Alternatively, this can be set with the following environment variable: " + tlsKeyFileEnvKey
	tlsKeyFileEnvKey = "TLS_KEY_FILE"

	tlsCACertsFlagName  = "tls-cacerts"
	tlsCACertsFlagUsage = "Comma-Separated list of ca certs path." +
		" Alternatively, this can be set with the following environment variable: " + tlsCACertsEnvKey
	tlsCACertsEnvKey = "TLS_CACERTS"

	// auto accept flag.
	agentAutoAcceptFlagName  = "auto-accept"
	agentAutoAcceptEnvKey    = "ARIESD_AUTO_ACCEPT"
	agentAutoAcceptFlagUsage = "Auto accept requests." +
		" Possible values [true] [false]. Defaults to false if not set." +
		" Alternatively, this can be set with the following environment variable: " + agentAutoAcceptEnvKey

	// log level.
	agentLogLevelFlagName  = "log-level"
	agentLogLevelEnvKey    = "ARIESD_LOG_LEVEL"
	agentLogLevelFlagUsage = "Log level." +
		" Possible values [INFO] [DEBUG] [ERROR] [WARNING] [CRITICAL] . Defaults to INFO if not set." +
		" Alternatively, this can be set with the following environment variable: " + agentLogLevelEnvKey

	// default label flag.
	agentDefaultLabelFlagName      = "agent-default-label"
	agentDefaultLabelEnvKey        = "ARIESD_DEFAULT_LABEL"
	agentDefaultLabelFlagShorthand = "l"
	agentDefaultLabelFlagUsage     = "Default Label for this agent. Defaults to blank if not set." +
		" Alternatively, this can be set with the following environment variable: " + agentDefaultLabelEnvKey

	// db namespace flag.
	agentDBNSFlagName      = "db-namespace"
	agentDBNSEnvKey        = "ARIESD_DB_NAMESPACE"
	agentDBNSFlagShorthand = "d"
	agentDBNSFlagUsage     = "database namespace." +
		" Alternatively, this can be set with the following environment variable: " + agentDBNSEnvKey

	// http resolver url flag.
	agentHTTPResolverFlagName      = "http-resolver-url"
	agentHTTPResolverEnvKey        = "ARIESD_HTTP_RESOLVER"
	agentHTTPResolverFlagShorthand = "r"
	agentHTTPResolverFlagUsage     = "HTTP binding DID resolver method and url. Values should be in `method@url` format." +
		" This flag can be repeated, allowing multiple http resolvers. Defaults to peer DID resolver if not set." +
		" Alternatively, this can be set with the following environment variable (in CSV format): " +
		agentHTTPResolverEnvKey

	blocDomainFlagName      = "bloc-domain"
	blocDomainFlagShorthand = "b"
	blocDomainFlagUsage     = "Bloc domain"
	blocDomainEnvKey        = "BLOC_DOMAIN"

	// wallet mediator url flag.
	walletMediatorURLFlagName      = "wallet-mediator-url"
	walletMediatorURLEnvKey        = "WALLET_MEDIATOR_URL"
	walletMediatorURLFlagShorthand = "m"
	walletMediatorURLFlagUsage     = "Mediator URL for wallet for performing DID communication" +
		" Alternatively, this can be set with the following environment variable: " +
		walletMediatorURLEnvKey

	// credential mediator url flag.
	credentialMediatorURLFlagName      = "credential-mediator-url"
	credentialMediatorURLEnvKey        = "CREDENTIAL_MEDIATOR_URL"
	credentialMediatorURLFlagShorthand = "a"
	credentialMediatorURLFlagUsage     = "credential mediator URL which provides Credential Mediator polyfill " +
		"for the W3C CCG Credential Handler API specification" +
		credentialMediatorURLEnvKey

	// blinded routing flag.
	blindedRoutingFlagName     = "blinded-routing"
	blindedRoutingEnvKey       = "BLINDED_ROUTING"
	blindedRoutingURLFlagUsage = "Flag to enable blinded routing to maintain identity privacy of the issuers " +
		"and verifiers involved. Possible values [true] [false]. Defaults to false if not set." +
		"Alternatively, this can be set with the following environment variable: " + blindedRoutingEnvKey

	// TODO Derive the SDS URL from the hub-auth bootstrap data #271.
	sdsURLFlagName      = "sds-url"
	sdsURLFlagShorthand = "s"
	sdsURLFlagUsage     = "URL SDS instance is running on."
	sdsURLEnvKey        = "HTTP_SERVER_SDS_URL"

	dependencyMaxRetriesFlagName   = "dep-maxretries"
	dependencyMaxRetriesFlagEnvKey = "HTTP_SERVER_DEP_MAXRETRIES"
	dependencyMaxRetriesFlagUsage  = "Optional. Sets the maximum number of retries while establishing connections with" +
		" external dependencies on startup. Default is 120. Delay between retries is 1s." +
		" Alternatively, this can be set with the following environment variable: " + dependencyMaxRetriesFlagEnvKey
	dependencyMaxRetriesDefault = uint64(120) // nolint:gomnd // false positive ("magic number")

	indexHTLMPath    = "/index.html"
	uiBasePath       = "/wallet/"
	uiConfigBasePath = "/walletconfig/"
	oidcBasePath     = "/oidc/"
	healthCheckPath  = "/healthcheck"
	deviceBasePath   = "/device/"
)

// OIDC config.
const (
	oidcProviderURLFlagName  = "oidc-opurl"
	oidcProviderURLFlagUsage = "URL for the OIDC provider." +
		" Alternatively, this can be set with the following environment variable: " + oidcProviderURLEnvKey
	oidcProviderURLEnvKey = "HTTP_SERVER_OIDC_OPURL"

	oidcClientIDFlagName  = "oidc-clientid"
	oidcClientIDFlagUsage = "OAuth2 client_id for OIDC." +
		" Alternatively, this can be set with the following environment variable: " + oidcClientIDEnvKey
	oidcClientIDEnvKey = "HTTP_SERVER_OIDC_CLIENTID"

	oidcClientSecretFlagName  = "oidc-clientsecret" // nolint:gosec // false positive on 'secret'
	oidcClientSecretFlagUsage = "OAuth2 client secret for OIDC." +
		" Alternatively, this can be set with the following environment variable: " + oidcClientSecretEnvKey
	oidcClientSecretEnvKey = "HTTP_SERVER_OIDC_CLIENTSECRET" // nolint:gosec // false positive on 'SECRET'

	// assumed to be the same landing page for all callbacks from all OIDC providers.
	oidcCallbackURLFlagName  = "oidc-callback"
	oidcCallbackURLFlagUsage = "Base URL for the OIDC callback endpoint." +
		" Alternatively, this can be set with the following environment variable: " + oidcCallbackURLEnvKey
	oidcCallbackURLEnvKey = "HTTP_SERVER_OIDC_CALLBACK"
)

// Keys.
const (
	sessionCookieAuthKeyFlagName  = "cookie-auth-key"
	sessionCookieAuthKeyFlagUsage = "Path to the pem-encoded 32-byte key to use to authenticate session cookies." +
		" Alternatively, this can be set with the following environment variable: " + sessionCookieAuthKeyEnvKey
	sessionCookieAuthKeyEnvKey = "HTTP_SERVER_COOKIE_AUTH_KEY"

	sessionCookieEncKeyFlagName  = "cookie-enc-key"
	sessionCookieEncKeyFlagUsage = "Path to the pem-encoded 32-byte key to use to encrypt session cookies." +
		" Alternatively, this can be set with the following environment variable: " + sessionCookieEncKeyEnvKey
	sessionCookieEncKeyEnvKey = "HTTP_SERVER_COOKIE_ENC_KEY"
)

// WebAuth Config.
const (
	webAuthRPDisplayFlagName  = "webauth-rp-displayname"
	webAuthRPDisplayEnvKey    = "HTTP_SERVER_RP_DISPLAY_NAME"
	webAuthRPDisplayFlagUsage = "WebAuth rp display name use to display name." +
		" Alternatively, this can be set with the following environment variable: " + webAuthRPDisplayEnvKey

	webAuthRPOriginFlagName  = "webauth-rp-origin"
	webAuthRPOriginEnvKey    = "HTTP_SERVER_RP_ORIGIN_NAME"
	webAuthRPOriginFlagUsage = "WebAuth rp origin url use for webauth requests ." +
		" Alternatively, this can be set with the following environment variable: " + webAuthRPOriginEnvKey

	webAuthRPIDFlagName  = "webauth-rp-id"
	webAuthRPIDEnvKey    = "HTTP_SERVER_RP_ID"
	webAuthRPIDFlagUsage = "WebAuth rp ID is generally the domain name for your site ." +
		" Alternatively, this can be set with the following environment variable: " + webAuthRPIDEnvKey
)

var logger = log.New("edge-agent/http-server")

type server interface {
	ListenAndServe(host, certFile, keyFile string, handler http.Handler) error
}

// HTTPServer represents an actual HTTP server implementation.
type HTTPServer struct{}

// ListenAndServe starts the server using the standard Go HTTP server implementation.
func (s *HTTPServer) ListenAndServe(host, certFile, keyFile string, handler http.Handler) error {
	if certFile != "" && keyFile != "" {
		return http.ListenAndServeTLS(host, certFile, keyFile, handler)
	}

	return http.ListenAndServe(host, handler)
}

type agentJSOpts struct {
	HTTPResolvedURLs  []string `json:"http-resolver-url,omitempty"`
	AgentDefaultLabel string   `json:"agent-default-label,omitempty"`
	AutoAccept        bool     `json:"auto-accept,omitempty"`
	LogLevel          string   `json:"log-level,omitempty"`
	DBNamespace       string   `json:"indexedDB-namespace,omitempty"`

	BlocDomain            string `json:"blocDomain,omitempty"`
	WalletMediatorURL     string `json:"walletMediatorURL,omitempty"`
	CredentialMediatorURL string `json:"credentialMediatorURL,omitempty"`
	BlindedRouting        bool   `json:"blindedRouting,omitempty"`
	SDSServerURL          string `json:"sdsServerURL,omitempty"`
}

type httpServerParameters struct {
	dependencyMaxRetries uint64
	srv                  server
	hostURL, wasmPath    string
	opts                 *agentJSOpts
	tls                  *tlsParameters
	oidc                 *oidcParameters
	keys                 *keyParameters
	webAuth              *webauthParameters
}

type tlsParameters struct {
	certFile string
	keyFile  string
	config   *tls.Config
}

type oidcParameters struct {
	providerURL  string
	clientID     string
	clientSecret string
	callbackURL  string
}

type webauthParameters struct {
	rpDisplayName string
	rpID          string
	rpOrigin      string
}

type keyParameters struct {
	sessionCookieAuthKey []byte
	sessionCookieEncKey  []byte
}

// GetStartCmd returns the Cobra start command.
func GetStartCmd(srv server) *cobra.Command {
	startCmd := createStartCmd(srv)

	createFlags(startCmd)

	return startCmd
}

func createStartCmd(srv server) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start http server",
		Long:  "Start http server",
		RunE: func(cmd *cobra.Command, args []string) error {
			hostURL, hostURLErr := cmdutils.GetUserSetVarFromString(cmd, hostURLFlagName, hostURLEnvKey, false)
			if hostURLErr != nil {
				return hostURLErr
			}

			wasmPath, err := cmdutils.GetUserSetVarFromString(cmd, wasmPathFlagName, wasmPathEnvKey, true)
			if err != nil {
				return err
			}

			tlsParams, err := getTLSParams(cmd)
			if err != nil {
				return err
			}

			opt, optErr := fetchAgentWASMOpts(cmd)
			if optErr != nil {
				return optErr
			}

			oidcParams, err := getOIDCParams(cmd)
			if err != nil {
				return err
			}

			retries, err := getDependencyMaxRetries(cmd)
			if err != nil {
				return err
			}

			keys, err := getKeyParams(cmd)
			if err != nil {
				return err
			}

			webAuthParams, err := getWebAuthParams(cmd)
			if err != nil {
				return err
			}

			parameters := &httpServerParameters{
				dependencyMaxRetries: retries,
				srv:                  srv,
				wasmPath:             wasmPath,
				hostURL:              hostURL,
				opts:                 opt,
				tls:                  tlsParams,
				oidc:                 oidcParams,
				webAuth:              webAuthParams,
				keys:                 keys,
			}

			return startHTTPServer(parameters)
		},
	}
}

func createFlags(startCmd *cobra.Command) {
	// wasm path flag
	startCmd.Flags().StringP(wasmPathFlagName, wasmPathFlagShorthand, "", wasmPathFlagUsage)
	// host url flag
	startCmd.Flags().StringP(hostURLFlagName, hostURLFlagShorthand, "", hostURLFlagUsage)
	// aries db path flag
	startCmd.Flags().StringP(agentDBNSFlagName, agentDBNSFlagShorthand, "", agentDBNSFlagUsage)
	// agent log level
	startCmd.Flags().StringP(agentLogLevelFlagName, "", "", agentLogLevelFlagUsage)
	// aries default label flag
	startCmd.Flags().StringP(agentDefaultLabelFlagName, agentDefaultLabelFlagShorthand, "",
		agentDefaultLabelFlagUsage)
	// aries auto accept flag
	startCmd.Flags().StringP(agentAutoAcceptFlagName, "", "", agentAutoAcceptFlagUsage)
	// aries http resolver url flag
	startCmd.Flags().StringArrayP(agentHTTPResolverFlagName, agentHTTPResolverFlagShorthand, []string{},
		agentHTTPResolverFlagUsage)
	// agent bloc domain
	startCmd.Flags().StringP(blocDomainFlagName, blocDomainFlagShorthand, "",
		blocDomainFlagUsage)
	// agent wallet mediator URL
	startCmd.Flags().StringP(walletMediatorURLFlagName, walletMediatorURLFlagShorthand, "",
		walletMediatorURLFlagUsage)
	// agent credential mediator URL
	startCmd.Flags().StringP(credentialMediatorURLFlagName, credentialMediatorURLFlagShorthand, "",
		credentialMediatorURLFlagUsage)
	// blinded routing for wallet
	startCmd.Flags().StringP(blindedRoutingFlagName, "", "",
		blindedRoutingURLFlagUsage)
	startCmd.Flags().StringP(sdsURLFlagName, sdsURLFlagShorthand, "", sdsURLFlagUsage)
	startCmd.Flags().StringP(dependencyMaxRetriesFlagName, "", "", dependencyMaxRetriesFlagUsage)
	createOIDCFlags(startCmd)
	createTLSFlags(startCmd)
	createKeyFlags(startCmd)
	createWebAuthFlags(startCmd)
}

func createTLSFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(tlsKeyFileFlagName, tlsKeyFileFlagShorthand, "", tlsKeyFileFlagUsage)
	cmd.Flags().StringP(tlsCertFileFlagName, tlsCertFileFlagShorthand, "", tlsCertFileFlagUsage)
	cmd.Flags().StringArrayP(tlsCACertsFlagName, "", []string{}, tlsCACertsFlagUsage)
}

func createOIDCFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(oidcProviderURLFlagName, "", "", oidcProviderURLFlagUsage)
	cmd.Flags().StringP(oidcClientIDFlagName, "", "", oidcClientIDFlagUsage)
	cmd.Flags().StringP(oidcClientSecretFlagName, "", "", oidcClientSecretFlagUsage)
	cmd.Flags().StringP(oidcCallbackURLFlagName, "", "", oidcCallbackURLFlagUsage)
}

func createKeyFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(sessionCookieAuthKeyFlagName, "", "", sessionCookieAuthKeyFlagUsage)
	cmd.Flags().StringP(sessionCookieEncKeyFlagName, "", "", sessionCookieEncKeyFlagUsage)
}

func createWebAuthFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(webAuthRPDisplayFlagName, "", "", webAuthRPDisplayFlagUsage)
	cmd.Flags().StringP(webAuthRPOriginFlagName, "", "", webAuthRPOriginFlagUsage)
	cmd.Flags().StringP(webAuthRPIDFlagName, "", "", webAuthRPIDFlagUsage)
}

func getDependencyMaxRetries(cmd *cobra.Command) (uint64, error) {
	retriesConfig, err := cmdutils.GetUserSetVarFromString(cmd,
		dependencyMaxRetriesFlagName, dependencyMaxRetriesFlagEnvKey, true)
	if err != nil {
		return 0, fmt.Errorf("failed to configure dependencyMaxRetries: %w", err)
	}

	maxRetries := dependencyMaxRetriesDefault

	if retriesConfig != "" {
		retries, err := strconv.ParseUint(retriesConfig, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse dependencyMaxRetries value '%s': %w", retriesConfig, err)
		}

		if retries > 0 {
			maxRetries = retries
		}
	}

	return maxRetries, nil
}

func getTLSParams(cmd *cobra.Command) (*tlsParameters, error) {
	params := &tlsParameters{}

	var err error

	params.certFile, err = cmdutils.GetUserSetVarFromString(cmd, tlsCertFileFlagName, tlsCertFileEnvKey, true)
	if err != nil {
		return nil, fmt.Errorf("failed to configure tls cert file: %w", err)
	}

	params.keyFile, err = cmdutils.GetUserSetVarFromString(cmd, tlsKeyFileFlagName, tlsKeyFileEnvKey, true)
	if err != nil {
		return nil, fmt.Errorf("failed to configure tls key file: %w", err)
	}

	rootCAs, err := cmdutils.GetUserSetVarFromArrayString(cmd, tlsCACertsFlagName, tlsCACertsEnvKey, true)
	if err != nil {
		return nil, fmt.Errorf("failed to configure root CAs: %w", err)
	}

	if len(rootCAs) > 0 {
		certPool, err := tlsutils.GetCertPool(false, rootCAs)
		if err != nil {
			return nil, fmt.Errorf("failed to init tls cert pool: %w", err)
		}

		params.config = &tls.Config{
			RootCAs:    certPool,
			MinVersion: tls.VersionTLS13,
		}
	}

	return params, nil
}

func getOIDCParams(cmd *cobra.Command) (*oidcParameters, error) {
	params := &oidcParameters{}

	var err error

	params.clientID, err = cmdutils.GetUserSetVarFromString(cmd, oidcClientIDFlagName, oidcClientIDEnvKey, false)
	if err != nil {
		return nil, fmt.Errorf("failed to configure OIDC clientID: %w", err)
	}

	params.clientSecret, err = cmdutils.GetUserSetVarFromString(
		cmd, oidcClientSecretFlagName, oidcClientSecretEnvKey, false)
	if err != nil {
		return nil, fmt.Errorf("failed to configure OIDC client secret: %w", err)
	}

	params.callbackURL, err = cmdutils.GetUserSetVarFromString(
		cmd, oidcCallbackURLFlagName, oidcCallbackURLEnvKey, false)
	if err != nil {
		return nil, fmt.Errorf("failed to configure OIDC callback URL: %w", err)
	}

	params.providerURL, err = cmdutils.GetUserSetVarFromString(
		cmd, oidcProviderURLFlagName, oidcProviderURLEnvKey, false)
	if err != nil {
		return nil, fmt.Errorf("failed to configure OIDC provider URL: %w", err)
	}

	return params, nil
}

func getWebAuthParams(cmd *cobra.Command) (*webauthParameters, error) {
	params := &webauthParameters{}

	var err error

	params.rpDisplayName, err = cmdutils.GetUserSetVarFromString(
		cmd, webAuthRPDisplayFlagName, webAuthRPDisplayEnvKey, true)
	if err != nil {
		return nil, fmt.Errorf("failed to configure webauth rp display name: %w", err)
	}

	params.rpID, err = cmdutils.GetUserSetVarFromString(
		cmd, webAuthRPIDFlagName, webAuthRPIDEnvKey, true)
	if err != nil {
		return nil, fmt.Errorf("failed to configure webauth rpID: %w", err)
	}

	params.rpOrigin, err = cmdutils.GetUserSetVarFromString(
		cmd, webAuthRPOriginFlagName, webAuthRPOriginEnvKey, true)
	if err != nil {
		return nil, fmt.Errorf("failed to configure webauth rp origin: %w", err)
	}

	return params, nil
}

func getKeyParams(cmd *cobra.Command) (*keyParameters, error) {
	params := &keyParameters{}

	sessionCookieAuthKeyPath, err := cmdutils.GetUserSetVarFromString(cmd,
		sessionCookieAuthKeyFlagName, sessionCookieAuthKeyEnvKey, false)
	if err != nil {
		return nil, fmt.Errorf("failed to configure session cookie auth key: %w", err)
	}

	params.sessionCookieAuthKey, err = parseKey(sessionCookieAuthKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to configure session cookie auth key: %w", err)
	}

	sessionCookieEncKeyPath, err := cmdutils.GetUserSetVarFromString(cmd,
		sessionCookieEncKeyFlagName, sessionCookieEncKeyEnvKey, false)
	if err != nil {
		return nil, fmt.Errorf("failed to configure session cookie enc key: %w", err)
	}

	params.sessionCookieEncKey, err = parseKey(sessionCookieEncKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to configure session cooie enc key: %w", err)
	}

	return params, nil
}

func parseKey(file string) ([]byte, error) {
	const (
		keyLen = 32
		bitNum = 8
	)

	bits, err := ioutil.ReadFile(filepath.Clean(file))
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", file, err)
	}

	if len(bits) != keyLen {
		return nil, fmt.Errorf("%s: need key of %d bits but got %d", file, keyLen*bitNum, len(bits)*bitNum)
	}

	return bits, nil
}

func initOIDCProvider(providerURL string, retries uint64, tlsConfig *tls.Config) (*oidcp.Provider, error) {
	var provider *oidcp.Provider

	err := backoff.RetryNotify(
		func() error {
			var provErr error
			provider, provErr = oidcp.NewProvider(
				oidcp.ClientContext(
					context.Background(),
					&http.Client{Transport: &http.Transport{
						TLSClientConfig: tlsConfig,
					}},
				),
				providerURL,
			)

			return provErr
		},
		backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), retries),
		func(retryErr error, d time.Duration) {
			fmt.Printf(
				"failed to init OIDC provider - will sleep for %s before trying again: %s\n", d, retryErr.Error())
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init OIDC provider: %w", err)
	}

	return provider, nil
}

func fetchAgentWASMOpts(cmd *cobra.Command) (*agentJSOpts, error) {
	opts := &agentJSOpts{}

	if err := setAriesOptions(opts, cmd); err != nil {
		return nil, err
	}

	if err := setTrustblocOptions(opts, cmd); err != nil {
		return nil, err
	}

	return opts, nil
}

//nolint:dupl // code differs in parameters that are processed
func setAriesOptions(opts *agentJSOpts, cmd *cobra.Command) error {
	defaultLabel, err := cmdutils.GetUserSetVarFromString(cmd,
		agentDefaultLabelFlagName, agentDefaultLabelEnvKey, true)
	if err != nil {
		return err
	}

	dbNS, err := cmdutils.GetUserSetVarFromString(cmd, agentDBNSFlagName, agentDBNSEnvKey, true)
	if err != nil {
		return err
	}

	logLevel, err := cmdutils.GetUserSetVarFromString(cmd, agentLogLevelFlagName, agentLogLevelEnvKey, true)
	if err != nil {
		return err
	}

	autoAcceptStr, err := cmdutils.GetUserSetVarFromString(cmd, agentAutoAcceptFlagName, agentAutoAcceptEnvKey, true)
	if err != nil {
		return err
	}

	autoAccept, err := parseBoolFlagValue(autoAcceptStr)
	if err != nil {
		return err
	}

	httpResolvers, err := cmdutils.GetUserSetVarFromArrayString(cmd,
		agentHTTPResolverFlagName, agentHTTPResolverEnvKey, true)
	if err != nil {
		return err
	}

	opts.AgentDefaultLabel = defaultLabel
	opts.DBNamespace = dbNS
	opts.LogLevel = logLevel
	opts.AutoAccept = autoAccept
	opts.HTTPResolvedURLs = httpResolvers

	return nil
}

//nolint:dupl // code differs in parameters that are processed
func setTrustblocOptions(opts *agentJSOpts, cmd *cobra.Command) error {
	blocDomain, err := cmdutils.GetUserSetVarFromString(cmd, blocDomainFlagName, blocDomainEnvKey, false)
	if err != nil {
		return err
	}

	walletMediatorURL, err := cmdutils.GetUserSetVarFromString(cmd,
		walletMediatorURLFlagName, walletMediatorURLEnvKey, true)
	if err != nil {
		return err
	}

	credentialMediatorURL, err := cmdutils.GetUserSetVarFromString(cmd,
		credentialMediatorURLFlagName, credentialMediatorURLEnvKey, true)
	if err != nil {
		return err
	}

	blindedRoutingStr, err := cmdutils.GetUserSetVarFromString(cmd, blindedRoutingFlagName, blindedRoutingEnvKey, true)
	if err != nil {
		return err
	}

	blindedRouting, err := parseBoolFlagValue(blindedRoutingStr)
	if err != nil {
		return err
	}

	sdsServerURL, err := cmdutils.GetUserSetVarFromString(cmd, sdsURLFlagName, sdsURLEnvKey, true)
	if err != nil {
		return err
	}

	opts.BlocDomain = blocDomain
	opts.WalletMediatorURL = walletMediatorURL
	opts.CredentialMediatorURL = credentialMediatorURL
	opts.BlindedRouting = blindedRouting
	opts.SDSServerURL = sdsServerURL

	return nil
}

func parseBoolFlagValue(v string) (bool, error) {
	if v == "" {
		return false, nil
	}

	val, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("invalid option - set true or false as the value")
	}

	return val, nil
}

func startHTTPServer(parameters *httpServerParameters) error {
	if parameters.wasmPath == "" {
		parameters.wasmPath = "."
	}

	err := setLogLevel(parameters.opts.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to set log level: %w", err)
	}

	router, err := router(parameters)
	if err != nil {
		return fmt.Errorf("failed to configure router: %w", err)
	}

	handler := cors.New(
		cors.Options{
			AllowedMethods: []string{http.MethodGet, http.MethodPost},
			AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
		},
	).Handler(router)

	logger.Infof("starting http-server on %s...", parameters.hostURL)

	err = parameters.srv.ListenAndServe(
		parameters.hostURL, parameters.tls.certFile, parameters.tls.keyFile,
		handler)
	if err != nil {
		return fmt.Errorf("http server closed unexpectedly: %s", err)
	}

	return err
}

func router(config *httpServerParameters) (http.Handler, error) {
	root := mux.NewRouter()

	root.HandleFunc(healthCheckPath, healthCheckHandler).Methods(http.MethodGet)

	uiRouter := root.PathPrefix(uiBasePath).Subrouter()
	addUIHandler(uiRouter, config.wasmPath)

	uiConfigRouter := root.PathPrefix(uiConfigBasePath).Subrouter()
	addUIConfigHandlers(uiConfigRouter, config)

	oidcRouter := root.PathPrefix(oidcBasePath).Subrouter()

	err := addOIDCHandlers(oidcRouter, config)
	if err != nil {
		return nil, fmt.Errorf("failed to add OIDC handlers: %w", err)
	}

	deviceRouter := root.PathPrefix(deviceBasePath).Subrouter()

	err = addDeviceHandlers(deviceRouter, config)
	if err != nil {
		return nil, fmt.Errorf("failed to add device handlers: %w", err)
	}

	return root, nil
}

func addUIHandler(router *mux.Router, publicDir string) {
	handler := gzipped.FileServer(http.Dir(publicDir))
	router.Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		logger.Infof("handling ui request: %s", urlPath)

		if strings.Contains(urlPath, ".") {
			r.URL.Path = r.URL.Path[len(uiBasePath):]
			handler.ServeHTTP(w, r)

			return
		}

		http.ServeFile(w, r, path.Join(publicDir, indexHTLMPath))
	})
}

func addUIConfigHandlers(router *mux.Router, config *httpServerParameters) {
	router.HandleFunc("/agent", func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("handling config request: %s", r.URL.String())
		bits, _ := json.Marshal(config.opts) // nolint:errcheck // marshalling *agentJSOpts never fails

		_, err := w.Write(bits)
		if err != nil {
			logger.Errorf("failed to write config %s to response: %w", bits, err)
		}
	})
}

func addOIDCHandlers(router *mux.Router, config *httpServerParameters) error {
	provider, err := initOIDCProvider(config.oidc.providerURL, config.dependencyMaxRetries, config.tls.config)
	if err != nil {
		return fmt.Errorf("failed to init OIDC provider: %w", err)
	}

	oidcOps, err := oidc.New(&oidc.Config{
		WalletDashboard: uiBasePath + "dashboard",
		TLSConfig:       config.tls.config,
		OIDCClient: oidc2.NewClient(&oidc2.Config{
			TLSConfig:    config.tls.config,
			Provider:     &oidc2.ProviderAdapter{OP: provider, TLSConfig: config.tls.config},
			CallbackURL:  config.oidc.callbackURL,
			ClientID:     config.oidc.clientID,
			ClientSecret: config.oidc.clientSecret,
			Scopes:       []string{oidcp.ScopeOpenID, "profile", "email"},
		}),
		Storage: &oidc.StorageConfig{
			Storage:          memstore.NewProvider(),
			TransientStorage: memstore.NewProvider(),
		},
		Keys: &oidc.KeyConfig{
			Auth: config.keys.sessionCookieAuthKey,
			Enc:  config.keys.sessionCookieEncKey,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to init oidc ops: %w", err)
	}

	for _, handler := range oidcOps.GetRESTHandlers() {
		router.HandleFunc(handler.Path(), handler.Handle()).Methods(handler.Method())
	}

	return nil
}

func addDeviceHandlers(router *mux.Router, config *httpServerParameters) error {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: config.webAuth.rpDisplayName, // Display Name for your site
		RPID:          config.webAuth.rpID,          // Generally the domain name for your site
		RPOrigin:      config.webAuth.rpOrigin,      // The origin URL for WebAuthn requests
	})
	if err != nil {
		return fmt.Errorf("failed to init webauthn: %w", err)
	}

	deviceOps, err := device.New(&device.Config{
		UIEndpoint: uiBasePath,
		TLSConfig:  config.tls.config,
		Storage: &device.StorageConfig{
			Storage:      memstore.NewProvider(),
			SessionStore: memstore.NewProvider(),
		},
		Keys: &device.KeyConfig{
			Auth: config.keys.sessionCookieAuthKey,
			Enc:  config.keys.sessionCookieEncKey,
		},
		Webauthn: webAuthn,
	})
	if err != nil {
		return fmt.Errorf("failed to init device ops: %w", err)
	}

	for _, handler := range deviceOps.GetRESTHandlers() {
		router.HandleFunc(handler.Path(), handler.Handle()).Methods(handler.Method())
	}

	return nil
}

type healthCheckResp struct {
	Status      string    `json:"status"`
	CurrentTime time.Time `json:"currentTime"`
}

func healthCheckHandler(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)

	err := json.NewEncoder(rw).Encode(&healthCheckResp{
		Status:      "success",
		CurrentTime: time.Now(),
	})
	if err != nil {
		logger.Errorf("healthcheck response failure, %s", err)
	}
}

func setLogLevel(logLevel string) error {
	if logLevel == "" {
		logLevel = "INFO"
	}

	return setEdgeCoreLogLevel(logLevel)
}

func setEdgeCoreLogLevel(logLevel string) error {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level '%s' : %w", logLevel, err)
	}

	log.SetLevel("", level)

	return nil
}
