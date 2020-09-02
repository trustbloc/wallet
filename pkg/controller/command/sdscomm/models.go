/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package sdscomm

import "encoding/json"

type SaveDIDDocToSDSRequest struct {
	Name     string          `json:"name,omitempty"`
	SignType string          `json:"signType,omitempty"`
	DID      json.RawMessage `json:"did,omitempty"`
	UserID   string          `json:"userID,omitempty"`
}

type SaveCredentialToSDSRequest struct {
	Name       string          `json:"name,omitempty"`
	Credential json.RawMessage `json:"credential,omitempty"`
	UserID     string          `json:"userID,omitempty"`
}

type SavePresentationToSDSRequest struct {
	Name         string          `json:"name,omitempty"`
	Presentation json.RawMessage `json:"presentation,omitempty"`
	UserID       string          `json:"userID,omitempty"`
}
