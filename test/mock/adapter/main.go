/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	demoPortEnvKey            = "DEMO_PORT"
	demoExternalURLEnvKey     = "EXTRERAL_URL"
	didCommInternalHostEnvKey = "INTERNAL_DIDCOMM_HOST"
	didCommExternalHostEnvKey = "EXTERNAL_DIDCOMM_HOST"
	tlsKeyFileEnvKey          = "TLS_KEY_FILE"
	tlsCertFileEnvKey         = "TLS_CERT_FILE"
	tlsCACertsEnvKey          = "TLS_CACERTS"
	orbDomainEnvKey           = "testnet.orb.local"
)

func main() {
	// initiate aries framework go options
	agent, err := startAriesAgent()
	if err != nil {
		panic(fmt.Errorf("failed to start aries-agent : %w", err))
	}

	router := mux.NewRouter()

	// host demo sample ui pages
	err = startVerifierApp(agent, router)
	if err != nil {
		panic(fmt.Errorf("failed to get verifier-app : %w", err))
	}

	// start the server
	port := os.Getenv(demoPortEnvKey)
	if port == "" {
		panic("port to be served not provided")
	}

	fmt.Println(http.ListenAndServeTLS(":"+port, os.Getenv(tlsCertFileEnvKey),
		os.Getenv(tlsKeyFileEnvKey), router))
}
