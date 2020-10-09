/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sdscomm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/trustbloc/edge-core/pkg/log"
	sdsclient "github.com/trustbloc/edv/pkg/client"
	"github.com/trustbloc/edv/pkg/edvutils"
	"github.com/trustbloc/edv/pkg/restapi/messages"
	"github.com/trustbloc/edv/pkg/restapi/models"
)

var logger = log.New("edge-agent-sdscomm")

// SDSComm implements SDS commands.
type SDSComm struct {
	SDSServerURL string
	sdsClient    *sdsclient.Client
}

// New returns a new SDSComm.
func New(sdsServerURL string) *SDSComm {
	return &SDSComm{
		SDSServerURL: sdsServerURL,
		sdsClient:    sdsclient.New(sdsServerURL),
	}
}

// StoreDIDDocument in the SDS.
func (e *SDSComm) StoreDIDDocument(saveDIDDocToSDSRequest *SaveDIDDocToSDSRequest) error {
	err := e.ensureVaultExists(e.getDIDVaultID(saveDIDDocToSDSRequest.UserID))
	if err != nil {
		return fmt.Errorf(failureEnsuringDIDVaultExistsErrMsg, err)
	}

	structuredDoc := models.StructuredDocument{
		ID: saveDIDDocToSDSRequest.Name,
	}

	structuredDoc.Content = make(map[string]interface{})

	structuredDoc.Content["didDoc"] = saveDIDDocToSDSRequest.DID
	structuredDoc.Content["signType"] = saveDIDDocToSDSRequest.SignType

	err = e.storeDocument(e.getDIDVaultID(saveDIDDocToSDSRequest.UserID), saveDIDDocToSDSRequest.Name, &structuredDoc)
	if err != nil {
		return fmt.Errorf(failureStoringDIDDocErrMsg, err)
	}

	return nil
}

// StoreCredential in the SDS.
// nolint:dupl // don't know how to de-duplicate this with StorePresentation without significant changes
func (e *SDSComm) StoreCredential(saveCredentialToSDSRequest *SaveCredentialToSDSRequest) error {
	err := e.ensureVaultExists(e.getCredentialVaultID(saveCredentialToSDSRequest.UserID))
	if err != nil {
		return fmt.Errorf(failureEnsuringCredVaultExistsErrMsg, err)
	}

	structuredDoc := models.StructuredDocument{
		ID: saveCredentialToSDSRequest.Name,
	}

	structuredDoc.Content = make(map[string]interface{})

	structuredDoc.Content["credential"] = saveCredentialToSDSRequest.Credential

	err = e.storeDocument(
		e.getCredentialVaultID(saveCredentialToSDSRequest.UserID),
		saveCredentialToSDSRequest.Name,
		&structuredDoc,
	)
	if err != nil {
		return fmt.Errorf(failureStoringCredErrMsg, err)
	}

	return nil
}

// StorePresentation in the SDS.
// nolint:dupl // don't know how to de-duplicate this with StoreCredential without significant changes
func (e *SDSComm) StorePresentation(savePresentationToSDSRequest *SavePresentationToSDSRequest) error {
	err := e.ensureVaultExists(e.getPresentationVaultID(savePresentationToSDSRequest.UserID))
	if err != nil {
		return fmt.Errorf(failureEnsuringPresVaultExistsErrMsg, err)
	}

	structuredDoc := models.StructuredDocument{
		ID: savePresentationToSDSRequest.Name,
	}

	structuredDoc.Content = make(map[string]interface{})

	structuredDoc.Content["presentation"] = savePresentationToSDSRequest.Presentation

	err = e.storeDocument(e.getPresentationVaultID(savePresentationToSDSRequest.UserID),
		savePresentationToSDSRequest.Name, &structuredDoc)
	if err != nil {
		return fmt.Errorf(failureStoringPresErrMsg, err)
	}

	return nil
}

func (e *SDSComm) ensureVaultExists(vaultID string) error {
	if e.SDSServerURL == "" {
		return errSDSServerURLBlank
	}

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

// TODO don't leak username to SDS: #265.
func (e *SDSComm) getDIDVaultID(userID string) string {
	return strings.ToLower(userID) + "_dids"
}

// TODO don't leak username to SDS: #265.
func (e *SDSComm) getCredentialVaultID(userID string) string {
	return strings.ToLower(userID) + "_credentials"
}

// TODO don't leak username to SDS: #265.
func (e *SDSComm) getPresentationVaultID(userID string) string {
	return strings.ToLower(userID) + "_presentations"
}

func (e *SDSComm) storeDocument(vaultID, friendlyName string, structuredDoc *models.StructuredDocument) error {
	encryptedDocID, err := edvutils.GenerateEDVCompatibleID()
	if err != nil {
		return fmt.Errorf(failureGeneratingEncryptedDocIDErrMsg, err)
	}

	indexedAttribute := models.IndexedAttribute{
		Name:   "FriendlyName",
		Value:  friendlyName, // TODO Don't leak friendly name to SDS #265
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
