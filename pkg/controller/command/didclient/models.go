/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package didclient

// CreateDIDRequest model
//
// This is used for creating DID
//
type CreateDIDRequest struct {
	// Optional public DID to be used in invitation
	PublicKey PublicKey `json:"publicKey,omitempty"`
}

// PublicKey public key
type PublicKey struct {
	ID    string `json:"id,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}
