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
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	oidcp "github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/lpar/gzipped"
	"github.com/spf13/cobra"
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

type ariesJSOpts struct {
	HTTPResolvedURLs  []string `json:"http-resolver-url,omitempty"`
	AgentDefaultLabel string   `json:"agent-default-label,omitempty"`
	AutoAccept        bool     `json:"auto-accept,omitempty"`
	LogLevel          string   `json:"log-level,omitempty"`
	DBNamespace       string   `json:"db-namespace,omitempty"`
}

type trustblocAgentJSOpts struct {
	BlocDomain            string `json:"blocDomain,omitempty"`
	WalletMediatorURL     string `json:"walletMediatorURL,omitempty"`
	CredentialMediatorURL string `json:"credentialMediatorURL,omitempty"`
	BlindedRouting        bool   `json:"blindedRouting,omitempty"`
	LogLevel              string `json:"log-level,omitempty"`
	SDSServerURL          string `json:"sdsServerURL,omitempty"`
}

type httpServerParameters struct {
	dependencyMaxRetries uint64
	srv                  server
	hostURL, wasmPath    string
	opts                 *ariesJSOpts
	trustblocAgentOpts   *trustblocAgentJSOpts
	tls                  *tlsParameters
	oidc                 *oidcParameters
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

			opt, optErr := fetchAriesWASMAgentOpts(cmd)
			if optErr != nil {
				return optErr
			}

			trustblocAgentOpts, err := fetchTrustBlocWASMAgentOpts(cmd)
			if err != nil {
				return err
			}

			oidcParams, err := getOIDCParams(cmd)
			if err != nil {
				return err
			}

			retries, err := getDependencyMaxRetries(cmd)
			if err != nil {
				return err
			}

			parameters := &httpServerParameters{
				dependencyMaxRetries: retries,
				srv:                  srv,
				wasmPath:             wasmPath,
				hostURL:              hostURL,
				opts:                 opt,
				trustblocAgentOpts:   trustblocAgentOpts,
				tls:                  tlsParams,
				oidc:                 oidcParams,
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
	// aries log level
	startCmd.Flags().StringP(agentLogLevelFlagName, "", "", agentLogLevelFlagUsage)
	// aries agent default label flag
	startCmd.Flags().StringP(agentDefaultLabelFlagName, agentDefaultLabelFlagShorthand, "",
		agentDefaultLabelFlagUsage)
	// aries auto accept flag
	startCmd.Flags().StringP(agentAutoAcceptFlagName, "", "", agentAutoAcceptFlagUsage)
	// aries http resolver url flag
	startCmd.Flags().StringArrayP(agentHTTPResolverFlagName, agentHTTPResolverFlagShorthand, []string{},
		agentHTTPResolverFlagUsage)
	// trustbloc agent bloc domain
	startCmd.Flags().StringP(blocDomainFlagName, blocDomainFlagShorthand, "",
		blocDomainFlagUsage)
	// trustbloc agent wallet mediator URL
	startCmd.Flags().StringP(walletMediatorURLFlagName, walletMediatorURLFlagShorthand, "",
		walletMediatorURLFlagUsage)
	// trustbloc agent credential mediator URL
	startCmd.Flags().StringP(credentialMediatorURLFlagName, credentialMediatorURLFlagShorthand, "",
		credentialMediatorURLFlagUsage)
	// blinded routing for wallet
	startCmd.Flags().StringP(blindedRoutingFlagName, "", "",
		blindedRoutingURLFlagUsage)
	startCmd.Flags().StringP(sdsURLFlagName, sdsURLFlagShorthand, "", sdsURLFlagUsage)
	startCmd.Flags().StringP(dependencyMaxRetriesFlagName, "", "", dependencyMaxRetriesFlagUsage)
	createOIDCFlags(startCmd)
	createTLSFlags(startCmd)
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

func fetchAriesWASMAgentOpts(cmd *cobra.Command) (*ariesJSOpts, error) {
	defaultLabel, err := cmdutils.GetUserSetVarFromString(
		cmd, agentDefaultLabelFlagName, agentDefaultLabelEnvKey, true)
	if err != nil {
		return nil, err
	}

	dbNS, err := cmdutils.GetUserSetVarFromString(cmd, agentDBNSFlagName, agentDBNSEnvKey, true)
	if err != nil {
		return nil, err
	}

	logLevel, err := cmdutils.GetUserSetVarFromString(cmd, agentLogLevelFlagName, agentLogLevelEnvKey, true)
	if err != nil {
		return nil, err
	}

	autoAcceptStr, err := cmdutils.GetUserSetVarFromString(cmd, agentAutoAcceptFlagName, agentAutoAcceptEnvKey, true)
	if err != nil {
		return nil, err
	}

	autoAccept, err := parseBoolFlagValue(autoAcceptStr)
	if err != nil {
		return nil, err
	}

	httpResolvers, err := cmdutils.GetUserSetVarFromArrayString(
		cmd, agentHTTPResolverFlagName, agentHTTPResolverEnvKey, true)
	if err != nil {
		return nil, err
	}

	return &ariesJSOpts{
		HTTPResolvedURLs:  httpResolvers,
		AgentDefaultLabel: defaultLabel,
		AutoAccept:        autoAccept,
		LogLevel:          logLevel,
		DBNamespace:       dbNS,
	}, nil
}

func fetchTrustBlocWASMAgentOpts(cmd *cobra.Command) (*trustblocAgentJSOpts, error) {
	blocDomain, err := cmdutils.GetUserSetVarFromString(cmd, blocDomainFlagName, blocDomainEnvKey, false)
	if err != nil {
		return nil, err
	}

	walletMediatorURL, err := cmdutils.GetUserSetVarFromString(cmd,
		walletMediatorURLFlagName, walletMediatorURLEnvKey, true)
	if err != nil {
		return nil, err
	}

	credentialMediatorURL, err := cmdutils.GetUserSetVarFromString(cmd,
		credentialMediatorURLFlagName, credentialMediatorURLEnvKey, true)
	if err != nil {
		return nil, err
	}

	blindedRoutingStr, err := cmdutils.GetUserSetVarFromString(cmd,
		blindedRoutingFlagName, blindedRoutingEnvKey, true)
	if err != nil {
		return nil, err
	}

	blindedRouting, err := parseBoolFlagValue(blindedRoutingStr)
	if err != nil {
		return nil, err
	}

	logLevel, err := cmdutils.GetUserSetVarFromString(cmd, agentLogLevelFlagName, agentLogLevelEnvKey, true)
	if err != nil {
		return nil, err
	}

	sdsServerURL, err := cmdutils.GetUserSetVarFromString(cmd, sdsURLFlagName, sdsURLEnvKey, true)
	if err != nil {
		return nil, err
	}

	return &trustblocAgentJSOpts{
		BlocDomain:            blocDomain,
		WalletMediatorURL:     walletMediatorURL,
		CredentialMediatorURL: credentialMediatorURL,
		BlindedRouting:        blindedRouting,
		LogLevel:              logLevel,
		SDSServerURL:          sdsServerURL,
	}, nil
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

	router, err := router(parameters)
	if err != nil {
		return fmt.Errorf("failed to configure router: %w", err)
	}

	logger.Infof("starting http-server on %s...", parameters.hostURL)

	err = parameters.srv.ListenAndServe(
		parameters.hostURL, parameters.tls.certFile, parameters.tls.keyFile,
		router)
	if err != nil {
		return fmt.Errorf("http server closed unexpectedly: %s", err)
	}

	return err
}

func router(config *httpServerParameters) (http.Handler, error) {
	root := mux.NewRouter()

	uiRouter := root.PathPrefix(uiBasePath).Subrouter()
	addUIHandler(uiRouter, config.wasmPath)

	uiConfigRouter := root.PathPrefix(uiConfigBasePath).Subrouter()
	addUIConfigHandlers(uiConfigRouter, config)

	oidcRouter := root.PathPrefix(oidcBasePath).Subrouter()

	err := addOIDCHandlers(oidcRouter, config)
	if err != nil {
		return nil, fmt.Errorf("failed to add OIDC handlers: %w", err)
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
	router.HandleFunc("/aries", func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("handling config request: %s", r.URL.String())
		bits, _ := json.Marshal(config.opts) // nolint:errcheck // marshalling *ariesJSOpts never fails

		_, err := w.Write(bits)
		if err != nil {
			logger.Errorf("failed to write config %s to response: %w", bits, err)
		}
	})
	router.HandleFunc("/trustbloc", func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("handling config request: %s", r.URL.String())
		bits, _ := json.Marshal(config.trustblocAgentOpts) // nolint:errcheck // marshalling *trustblocAgentJSOpts never fails

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
		OIDC: &oidc.OIDCConfig{
			Provider:     &oidc.OIDCProviderImpl{OP: provider},
			ClientID:     config.oidc.clientID,
			ClientSecret: config.oidc.clientSecret,
			Scopes:       []string{oidcp.ScopeOpenID, "profile", "email"},
			CallbackURL:  config.oidc.callbackURL,
		},
		Storage: &oidc.StorageConfig{
			Storage:          memstore.NewProvider(),
			TransientStorage: memstore.NewProvider(),
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
