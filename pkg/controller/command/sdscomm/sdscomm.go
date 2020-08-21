/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sdscomm

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/trustbloc/edv/pkg/edvutils"

	"github.com/trustbloc/edge-core/pkg/log"
	sdsclient "github.com/trustbloc/edv/pkg/client"
	"github.com/trustbloc/edv/pkg/restapi/messages"
	"github.com/trustbloc/edv/pkg/restapi/models"
)

var logger = log.New("edge-agent-sdscomm")

var errBlankSDSURL = errors.New("SDS server URL cannot be blank")
var errBlankAgentUsername = errors.New("agent username cannot be blank")

type SDSComm struct {
	sdsServerURL  string
	agentUsername string
	sdsClient     *sdsclient.Client
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

func (e *SDSComm) StoreDIDDocument(didData *DIDDocData) error {
	err := e.ensureVaultExists(e.getDIDVaultID())
	if err != nil {
		return fmt.Errorf(failureEnsuringDIDVaultExistsErrMsg, err)
	}

	structuredDoc := models.StructuredDocument{
		ID: didData.Name,
	}

	structuredDoc.Content = make(map[string]interface{})

	structuredDoc.Content["didDoc"] = didData.DID
	structuredDoc.Content["signType"] = didData.SignType

	err = e.storeDocument(e.getDIDVaultID(), didData.Name, &structuredDoc)
	if err != nil {
		return fmt.Errorf(failureStoringDIDDocErrMsg, err)
	}

	return nil
}

func (e *SDSComm) StoreCredential(credentialData *CredentialData) error {
	err := e.ensureVaultExists(e.getCredentialVaultID())
	if err != nil {
		return fmt.Errorf(failureEnsuringCredVaultExistsErrMsg, err)
	}

	structuredDoc := models.StructuredDocument{
		ID: credentialData.Name,
	}

	structuredDoc.Content = make(map[string]interface{})

	structuredDoc.Content["credential"] = credentialData.Credential

	err = e.storeDocument(e.getCredentialVaultID(), credentialData.Name, &structuredDoc)
	if err != nil {
		return fmt.Errorf(failureStoringCredErrMsg, err)
	}

	return nil
}

func (e *SDSComm) StorePresentation(presentationData *PresentationData) error {
	err := e.ensureVaultExists(e.getPresentationVaultID())
	if err != nil {
		return fmt.Errorf(failureEnsuringPresVaultExistsErrMsg, err)
	}

	structuredDoc := models.StructuredDocument{
		ID: presentationData.Name,
	}

	structuredDoc.Content = make(map[string]interface{})

	structuredDoc.Content["presentation"] = presentationData.Presentation

	err = e.storeDocument(e.getPresentationVaultID(), presentationData.Name, &structuredDoc)
	if err != nil {
		return fmt.Errorf(failureStoringPresErrMsg, err)
	}

	return nil
}

func (e *SDSComm) ensureVaultExists(vaultID string) error {
	_, err := e.sdsClient.CreateDataVault(&models.DataVaultConfiguration{ReferenceID: vaultID})
	if err != nil {
		if !strings.Contains(err.Error(), messages.ErrDuplicateVault.Error()) {
			return fmt.Errorf(unexpectedErrorOnCreateVaultCall, err)
		}

		logger.Debugf(vaultAlreadyExistsLogMsg, vaultID)
	} else {
		logger.Debugf(newVaultCreatedLogMsg, vaultID)
	}

	return nil
}

// TODO don't leak username to SDS: #265
func (e *SDSComm) getDIDVaultID() string {
	return e.agentUsername + "_dids"
}

// TODO don't leak username to SDS: #265
func (e *SDSComm) getCredentialVaultID() string {
	return e.agentUsername + "_credentials"
}

// TODO don't leak username to SDS: #265
func (e *SDSComm) getPresentationVaultID() string {
	return e.agentUsername + "_presentations"
}

func (e *SDSComm) storeDocument(vaultID, friendlyName string, structuredDoc *models.StructuredDocument) error {
	encryptedDocID, err := edvutils.GenerateEDVCompatibleID()
	if err != nil {
		return fmt.Errorf(failureGeneratingEncryptedDocIDErrMsg, err)
	}

	indexedAttribute := models.IndexedAttribute{
		Name:   "FriendlyName",
		Value:  friendlyName, //TODO Don't leak friendly name to SDS #265
		Unique: true,
	}

	indexedAttributeCollection := models.IndexedAttributeCollection{
		IndexedAttributes: []models.IndexedAttribute{indexedAttribute},
	}

	indexedAttributeCollections := []models.IndexedAttributeCollection{indexedAttributeCollection}

	structuredDocBytes, err := json.Marshal(structuredDoc)
	if err != nil {
		return fmt.Errorf(failureMarshallingStructuredDocErrMsg, err)
	}

	// TODO encrypt data before storing in SDS #267
	encryptedDoc := models.EncryptedDocument{
		ID:                          encryptedDocID,
		Sequence:                    0,
		IndexedAttributeCollections: indexedAttributeCollections,
		JWE:                         structuredDocBytes,
	}

	_, err = e.sdsClient.CreateDocument(vaultID, &encryptedDoc)
	if err != nil {
		return fmt.Errorf(failureStoringDocErrMsg, err)
	}

	return nil
}
