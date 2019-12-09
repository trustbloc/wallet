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
)

type server interface {
	ListenAndServe(host, rootPath string) error
}

// HTTPServer represents an actual HTTP server implementation.
type HTTPServer struct{}

// ListenAndServe starts the server using the standard Go HTTP server implementation.
func (s *HTTPServer) ListenAndServe(host, rootPath string) error {
	h := gzipped.FileServer(http.Dir(rootPath))
	return http.ListenAndServe(host, h)
}

type httpServerParameters struct {
	srv               server
	hostURL, wasmPath string
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

			parameters := &httpServerParameters{
				srv:      srv,
				hostURL:  hostURL,
				wasmPath: wasmPath,
			}
			return startHTTPServer(parameters)
		},
	}
}

func createFlags(startCmd *cobra.Command) {
	startCmd.Flags().StringP(hostURLFlagName, hostURLFlagShorthand, "", hostURLFlagUsage)
	startCmd.Flags().StringP(wasmPathFlagName, wasmPathFlagShorthand, "", wasmPathFlagUsage)
}

func getUserSetVar(cmd *cobra.Command, flagName, envKey string, isOptional bool) (string, error) {
	if cmd.Flags().Changed(flagName) {
		value, err := cmd.Flags().GetString(flagName)
		if err != nil {
			return "", fmt.Errorf(flagName+" flag not found: %s", err)
		}

		return value, nil
	}

	value, isSet := os.LookupEnv(envKey)

	if isOptional || isSet {
		return value, nil
	}

	return "", errors.New("Neither " + flagName + " (command line flag) nor " + envKey +
		" (environment variable) have been set.")
}

func startHTTPServer(parameters *httpServerParameters) error {
	if parameters.hostURL == "" {
		return errors.New("host URL not provided")
	}

	if parameters.wasmPath == "" {
		parameters.wasmPath = "."
	}

	err := parameters.srv.ListenAndServe(parameters.hostURL, parameters.wasmPath)
	if err != nil {
		return fmt.Errorf("http server closed unexpectedly: %s", err)
	}

	return err
}
