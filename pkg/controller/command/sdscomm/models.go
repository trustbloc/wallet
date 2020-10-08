/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package sdscomm

import "encoding/json"

// SaveDIDDocToSDSRequest holds parameters for saving a DID document to the SDS.
type SaveDIDDocToSDSRequest struct {
	Name     string          `json:"name,omitempty"`
	SignType string          `json:"signType,omitempty"`
	DID      json.RawMessage `json:"did,omitempty"`
	UserID   string          `json:"userID,omitempty"`
}

// SaveCredentialToSDSRequest holds parameters for saving a VC to the SDS.
type SaveCredentialToSDSRequest struct {
	Name       string          `json:"name,omitempty"`
	Credential json.RawMessage `json:"credential,omitempty"`
	UserID     string          `json:"userID,omitempty"`
}

// SavePresentationToSDSRequest holds parameters for saving a VP to the SDS.
type SavePresentationToSDSRequest struct {
	Name         string          `json:"name,omitempty"`
	Presentation json.RawMessage `json:"presentation,omitempty"`
	UserID       string          `json:"userID,omitempty"`
}
