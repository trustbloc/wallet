/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/hyperledger/aries-framework-go/pkg/doc/presexch"
	"github.com/hyperledger/aries-framework-go/pkg/store/ld"
	"github.com/hyperledger/aries-framework-go/spi/storage"
)

// StoreProvider provides stores for JSON-LD contexts and remote providers.
type StoreProvider struct {
	ContextStore        ld.ContextStore
	RemoteProviderStore ld.RemoteProviderStore
}

// NewLDStoreProvider returns a new instance of StoreProvider.
func NewLDStoreProvider(storageProvider storage.Provider) (*StoreProvider, error) {
	contextStore, err := ld.NewContextStore(storageProvider)
	if err != nil {
		return nil, fmt.Errorf("create JSON-LD context store: %w", err)
	}

	remoteProviderStore, err := ld.NewRemoteProviderStore(storageProvider)
	if err != nil {
		return nil, fmt.Errorf("create remote provider store: %w", err)
	}

	return &StoreProvider{
		ContextStore:        contextStore,
		RemoteProviderStore: remoteProviderStore,
	}, nil
}

// JSONLDContextStore returns JSON-LD context store.
func (p *StoreProvider) JSONLDContextStore() ld.ContextStore {
	return p.ContextStore
}

// JSONLDRemoteProviderStore returns JSON-LD remote provider store.
func (p *StoreProvider) JSONLDRemoteProviderStore() ld.RemoteProviderStore {
	return p.RemoteProviderStore
}

func loadTemplate(w http.ResponseWriter, fileName string, data map[string]interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := t.Execute(w, data); err != nil {
		return err
	}

	return nil
}

func handleError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(ErrorResponse{
		Message: msg,
	})

	if err != nil {
		logger.Errorf("Unable to send error message, %s", err)
	}
}

// ErrorResponse to send error message in the response.
type ErrorResponse struct {
	Message string `json:"errMessage,omitempty"`
}

type OIDCAuthClaims struct {
	VPToken *VPToken `json:"vp_token"`
}

type VPToken struct {
	PresDef *presexch.PresentationDefinition `json:"presentation_definition"`
}

type OIDCTokenCliams struct {
	VPToken *VPTokenClaim `json:"_vp_token"`
}

type VPTokenClaim struct {
	PresSub *presexch.PresentationSubmission `json:"presentation_submission"`
}
