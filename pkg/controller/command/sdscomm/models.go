/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package sdscomm

import "encoding/json"

type DIDDocData struct {
	Name     string          `json:"name,omitempty"`
	SignType string          `json:"signType,omitempty"`
	DID      json.RawMessage `json:"did,omitempty"`
}

type CredentialData struct {
	Name       string          `json:"name,omitempty"`
	Credential json.RawMessage `json:"credential,omitempty"`
}
