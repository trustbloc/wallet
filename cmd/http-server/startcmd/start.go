/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package startcmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/lpar/gzipped"
	"github.com/spf13/cobra"
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
)

type server interface {
	ListenAndServe(host, certFile, keyFile, rootPath string) error
}

// HTTPServer represents an actual HTTP server implementation.
type HTTPServer struct{}

// ListenAndServe starts the server using the standard Go HTTP server implementation.
func (s *HTTPServer) ListenAndServe(host, certFile, keyFile, rootPath string) error {
	return http.ListenAndServeTLS(host, certFile, keyFile, VueHandler(rootPath))
}

// VueHandler return a http.Handler that supports Vue Router app with history mode
func VueHandler(publicDir string) http.Handler {
	handler := gzipped.FileServer(http.Dir(publicDir))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		urlPath := req.URL.Path

		// static files
		if strings.Contains(urlPath, ".") || urlPath == "/" {
			handler.ServeHTTP(w, req)
			return
		}

		// the all 404 gonna be served as root
		http.ServeFile(w, req, path.Join(publicDir, "/index.html"))
	})
}

type httpServerParameters struct {
	srv                                        server
	hostURL, wasmPath, tlsCertFile, tlsKeyFile string
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
			hostURL, err := getUserSetVar(cmd, hostURLFlagName, hostURLEnvKey, false)
			if err != nil {
				return err
			}

			wasmPath, err := getUserSetVar(cmd, wasmPathFlagName, wasmPathEnvKey, true)
			if err != nil {
				return err
			}

			tlsCertFile, err := getUserSetVar(cmd, tlsCertFileFlagName, tlsCertFileEnvKey, false)
			if err != nil {
				return err
			}

			tlsKeyFile, err := getUserSetVar(cmd, tlsKeyFileFlagName, tlsKeyFileEnvKey, false)
			if err != nil {
				return err
			}

			parameters := &httpServerParameters{
				srv:         srv,
				hostURL:     hostURL,
				wasmPath:    wasmPath,
				tlsCertFile: tlsCertFile,
				tlsKeyFile:  tlsKeyFile,
			}
			return startHTTPServer(parameters)
		},
	}
}

func createFlags(startCmd *cobra.Command) {
	startCmd.Flags().StringP(hostURLFlagName, hostURLFlagShorthand, "", hostURLFlagUsage)
	startCmd.Flags().StringP(wasmPathFlagName, wasmPathFlagShorthand, "", wasmPathFlagUsage)
	startCmd.Flags().StringP(tlsCertFileFlagName, tlsCertFileFlagShorthand, "", tlsCertFileFlagUsage)
	startCmd.Flags().StringP(tlsKeyFileFlagName, tlsKeyFileFlagShorthand, "", tlsKeyFileFlagUsage)
}

func getUserSetVar(cmd *cobra.Command, flagName, envKey string, isOptional bool) (string, error) {
	if cmd.Flags().Changed(flagName) {
		value, err := cmd.Flags().GetString(flagName)
		if err != nil {
			return "", fmt.Errorf(flagName+" flag not found: %s", err)
		}

		if value == "" {
			return "", fmt.Errorf("%s value is empty", flagName)
		}

		return value, nil
	}

	value, isSet := os.LookupEnv(envKey)

	if isOptional || isSet {
		if !isOptional && value == "" {
			return "", fmt.Errorf("%s value is empty", envKey)
		}

		return value, nil
	}

	return "", errors.New("Neither " + flagName + " (command line flag) nor " + envKey +
		" (environment variable) have been set.")
}

func startHTTPServer(parameters *httpServerParameters) error {
	if parameters.wasmPath == "" {
		parameters.wasmPath = "."
	}

	err := parameters.srv.ListenAndServe(
		parameters.hostURL, parameters.tlsCertFile, parameters.tlsKeyFile, parameters.wasmPath)
	if err != nil {
		return fmt.Errorf("http server closed unexpectedly: %s", err)
	}

	return err
}
