/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"github.com/rs/cors"
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
	orbDomainEnvKey           = "ORB_DOMAIN"
	contextProviderEnvKey     = "CONTEXT_PROVIDER_URL"
	keyTypeEnvKey             = "KEY_TYPE"
	keyAgreementTypeEnvKey    = "KEY_AGREEMENT_TYPE"
)

func main() {
	// initiate aries framework go options
	agent, err := startAriesAgent()
	if err != nil {
		panic(fmt.Errorf("failed to start aries-agent : %w", err))
	}

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("templates"))
	router.PathPrefix("/css/").Handler(fs)
	router.Handle("/", fs)

	// host demo sample ui pages
	err = startAdapterApp(agent, router)
	if err != nil {
		panic(fmt.Errorf("failed to get verifier-app : %w", err))
	}

	// start the server
	port := os.Getenv(demoPortEnvKey)
	if port == "" {
		panic("port to be served not provided")
	}

	handler := cors.New(
		cors.Options{
			AllowedMethods: []string{http.MethodGet, http.MethodPost},
			AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
		},
	).Handler(router)

	fmt.Println(http.ListenAndServeTLS(":"+port, os.Getenv(tlsCertFileEnvKey),
		os.Getenv(tlsKeyFileEnvKey), handler))
}
