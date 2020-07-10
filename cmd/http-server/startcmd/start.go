/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package startcmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/lpar/gzipped"
	"github.com/spf13/cobra"
	cmdutils "github.com/trustbloc/edge-core/pkg/utils/cmd"
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

	// auto accept flag
	agentAutoAcceptFlagName  = "auto-accept"
	agentAutoAcceptEnvKey    = "ARIESD_AUTO_ACCEPT"
	agentAutoAcceptFlagUsage = "Auto accept requests." +
		" Possible values [true] [false]. Defaults to false if not set." +
		" Alternatively, this can be set with the following environment variable: " + agentAutoAcceptEnvKey

	// log level
	agentLogLevelFlagName  = "log-level"
	agentLogLevelEnvKey    = "ARIESD_LOG_LEVEL"
	agentLogLevelFlagUsage = "Log level." +
		" Possible values [INFO] [DEBUG] [ERROR] [WARNING] [CRITICAL] . Defaults to INFO if not set." +
		" Alternatively, this can be set with the following environment variable: " + agentLogLevelEnvKey

	// default label flag
	agentDefaultLabelFlagName      = "agent-default-label"
	agentDefaultLabelEnvKey        = "ARIESD_DEFAULT_LABEL"
	agentDefaultLabelFlagShorthand = "l"
	agentDefaultLabelFlagUsage     = "Default Label for this agent. Defaults to blank if not set." +
		" Alternatively, this can be set with the following environment variable: " + agentDefaultLabelEnvKey

	// db namespace flag
	agentDBNSFlagName      = "db-namespace"
	agentDBNSEnvKey        = "ARIESD_DB_NAMESPACE"
	agentDBNSFlagShorthand = "d"
	agentDBNSFlagUsage     = "database namespace." +
		" Alternatively, this can be set with the following environment variable: " + agentDBNSEnvKey

	// http resolver url flag
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

	// wallet mediator url flag
	walletMediatorURLFlagName      = "wallet-mediator-url"
	walletMediatorURLEnvKey        = "WALLET_MEDIATOR_URL"
	walletMediatorURLFlagShorthand = "m"
	walletMediatorURLFlagUsage     = "Mediator URL for wallet for performing DID communication" +
		" Alternatively, this can be set with the following environment variable: " +
		agentHTTPResolverEnvKey

	// aries opts path
	ariesStartupOptsPath = "/aries/jsopts"
	indexHTLMPath        = "/index.html"
	basePath             = "/"

	// tustbloc agent opt path
	trustblocStartupOptsPath = "/trustbloc-agent/jsopts"
)

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
	BlocDomain        string `json:"blocDomain,omitempty"`
	WalletMediatorURL string `json:"walletMediatorURL,omitempty"`
	LogLevel          string `json:"log-level,omitempty"`
}

// VueHandler return a http.Handler that supports Vue Router app with history mode
func VueHandler(publicDir string, opts *ariesJSOpts, trustblocAgentOpts *trustblocAgentJSOpts) http.Handler {
	handler := gzipped.FileServer(http.Dir(publicDir))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		urlPath := req.URL.Path

		// aries js opts
		if urlPath == ariesStartupOptsPath {
			j, _ := json.Marshal(opts) // nolint errcheck

			w.Write(j) // nolint errcheck
			return
		}

		// trustbloc agent js opts
		if urlPath == trustblocStartupOptsPath {
			j, _ := json.Marshal(trustblocAgentOpts) // nolint errcheck

			w.Write(j) // nolint errcheck
			return
		}

		// static files
		if urlPath == basePath || strings.Contains(urlPath, ".") {
			handler.ServeHTTP(w, req)
			return
		}

		// the all 404 gonna be served as root
		http.ServeFile(w, req, path.Join(publicDir, indexHTLMPath))
	})
}

type httpServerParameters struct {
	srv                                        server
	hostURL, wasmPath, tlsCertFile, tlsKeyFile string
	opts                                       *ariesJSOpts
	trustblocAgentOpts                         *trustblocAgentJSOpts
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

			tlsCertFile, tlsCertFileErr := cmdutils.GetUserSetVarFromString(cmd, tlsCertFileFlagName, tlsCertFileEnvKey, true)
			if tlsCertFileErr != nil {
				return tlsCertFileErr
			}

			tlsKeyFile, tlsKeyFileErr := cmdutils.GetUserSetVarFromString(cmd, tlsKeyFileFlagName, tlsKeyFileEnvKey, true)
			if tlsKeyFileErr != nil {
				return tlsKeyFileErr
			}

			opt, optErr := fetchAriesWASMAgentOpts(cmd)
			if optErr != nil {
				return optErr
			}

			trustblocAgentOpts, err := fetchTrustBlocWASMAgentOpts(cmd)
			if err != nil {
				return err
			}

			parameters := &httpServerParameters{
				srv:                srv,
				wasmPath:           wasmPath,
				hostURL:            hostURL,
				tlsCertFile:        tlsCertFile,
				tlsKeyFile:         tlsKeyFile,
				opts:               opt,
				trustblocAgentOpts: trustblocAgentOpts,
			}
			return startHTTPServer(parameters)
		},
	}
}

func createFlags(startCmd *cobra.Command) {
	// wasm path flag
	startCmd.Flags().StringP(wasmPathFlagName, wasmPathFlagShorthand, "", wasmPathFlagUsage)
	// tls cert key flag
	startCmd.Flags().StringP(tlsKeyFileFlagName, tlsKeyFileFlagShorthand, "", tlsKeyFileFlagUsage)
	// host url flag
	startCmd.Flags().StringP(hostURLFlagName, hostURLFlagShorthand, "", hostURLFlagUsage)
	// tls cert file flag
	startCmd.Flags().StringP(tlsCertFileFlagName, tlsCertFileFlagShorthand, "", tlsCertFileFlagUsage)
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

	autoAccept, err := getAutoAcceptValue(cmd)
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

	logLevel, err := cmdutils.GetUserSetVarFromString(cmd, agentLogLevelFlagName, agentLogLevelEnvKey, true)
	if err != nil {
		return nil, err
	}

	return &trustblocAgentJSOpts{
		BlocDomain:        blocDomain,
		WalletMediatorURL: walletMediatorURL,
		LogLevel:          logLevel,
	}, nil
}

func getAutoAcceptValue(cmd *cobra.Command) (bool, error) {
	v, err := cmdutils.GetUserSetVarFromString(cmd, agentAutoAcceptFlagName, agentAutoAcceptEnvKey, true)
	if err != nil {
		return false, err
	}

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

	err := parameters.srv.ListenAndServe(
		parameters.hostURL, parameters.tlsCertFile, parameters.tlsKeyFile,
		VueHandler(parameters.wasmPath, parameters.opts, parameters.trustblocAgentOpts))
	if err != nil {
		return fmt.Errorf("http server closed unexpectedly: %s", err)
	}

	return err
}
