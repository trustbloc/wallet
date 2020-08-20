/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sdscomm

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/trustbloc/edge-core/pkg/log"
	sdsclient "github.com/trustbloc/edv/pkg/client"
	"github.com/trustbloc/edv/pkg/restapi/messages"
	"github.com/trustbloc/edv/pkg/restapi/models"
)

const sdsVaultIDLen = 16

var logger = log.New("edge-agent-sdscomm")

var errBlankSDSURL = errors.New("SDS server URL cannot be blank")
var errBlankAgentUsername = errors.New("agent username cannot be blank")

type SDSComm struct {
	sdsServerURL  string
	agentUsername string
	sdsClient     *sdsclient.Client
}

type DIDDocData struct {
	DID      json.RawMessage `json:"did,omitempty"`
	SignType string          `json:"signType,omitempty"`
	Name     string          `json:"name,omitempty"`
}

func New(sdsServerURL, agentUsername string) (*SDSComm, error) {
	if sdsServerURL == "" {
		return nil, errBlankSDSURL
	}

	if agentUsername == "" {
		return nil, errBlankAgentUsername
	}

	return &SDSComm{
		sdsServerURL:  sdsServerURL,
		agentUsername: agentUsername,
		sdsClient:     sdsclient.New(sdsServerURL),
	}, nil
}

// CreateDIDVault creates the user's DID vault if it doesn't exist.
func (e *SDSComm) CreateDIDVault() error {
	didVaultID := e.GetDIDVaultID()

	_, err := e.sdsClient.CreateDataVault(&models.DataVaultConfiguration{ReferenceID: didVaultID})
	if err != nil {
		if !strings.Contains(err.Error(), messages.ErrDuplicateVault.Error()) {
			return err
		}

		logger.Debugf("%s vault already exists. Skipping vault creation.", didVaultID)
	}

	return nil
}

func (e *SDSComm) StoreDIDDocument(didData *DIDDocData) error {
	structuredDoc := models.StructuredDocument{
		ID: didData.Name,
	}

	structuredDoc.Content = make(map[string]interface{})

	structuredDoc.Content["didDoc"] = didData.DID
	structuredDoc.Content["signType"] = didData.SignType

	encryptedDocID, err := generateSDSCompatibleID()
	if err != nil {
		return fmt.Errorf("failed to generate encrypted doc ID: %w", err)
	}

	indexedAttribute := models.IndexedAttribute{
		Name:   "FriendlyName",
		Value:  didData.Name, //TODO Don't leak friendly name to SDS #265
		Unique: true,
	}

	indexedAttributeCollection := models.IndexedAttributeCollection{
		IndexedAttributes: []models.IndexedAttribute{indexedAttribute},
	}

	indexedAttributeCollections := []models.IndexedAttributeCollection{indexedAttributeCollection}

	structuredDocBytes, err := json.Marshal(structuredDoc)

	// TODO encrypt data before storing in SDS #267
	encryptedDoc := models.EncryptedDocument{
		ID:                          encryptedDocID,
		Sequence:                    0,
		IndexedAttributeCollections: indexedAttributeCollections,
		JWE:                         structuredDocBytes,
	}

	_, err = e.sdsClient.CreateDocument(e.GetDIDVaultID(), &encryptedDoc)
	if err != nil {
		return fmt.Errorf("failed to store DID document: %w", err)
	}

	return nil
}

// TODO don't leak username to SDS: #265
func (e *SDSComm) GetDIDVaultID() string {
	return e.agentUsername + "_dids"
}

func generateSDSCompatibleID() (string, error) {
	randomBytes := make([]byte, sdsVaultIDLen)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	base58EncodedUUID := base58.Encode(randomBytes)

	return base58EncodedUUID, nil
}
