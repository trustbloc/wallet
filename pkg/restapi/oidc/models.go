/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package oidc

import "encoding/json"

type createKeystoreReq struct {
	Controller string `json:"controller,omitempty"`
	VaultID    string `json:"vaultID,omitempty"`
}

type createKeyReq struct {
	KeyType string `json:"keyType,omitempty"`
}

type exportKeyResp struct {
	PublicKey string `json:"publicKey,omitempty"`
}

type updateCapabilityReq struct {
	EDVCapability json.RawMessage `json:"edvCapability,omitempty"`
}

type signReq struct {
	Message string `json:"message,omitempty"`
}

type signResp struct {
	Signature string `json:"signature,omitempty"`
}

// TODO delete during completion of https://github.com/trustbloc/edge-agent/issues/489.
type todoDeleteThisModel struct {
	UserEDVVaultURL   string          `json:"edvVaultURL"`
	OpsEDVVaultURL    string          `json:"opsVaultURL"`
	AuthzKeyStoreURL  string          `json:"authKeyStoreURL"`
	OpsKeyStoreURL    string          `json:"opsKeyStoreURL"`
	UserEDVCapability json.RawMessage `json:"uerEDVCapability,omitempty"`
}

type secretRequest struct {
	Secret []byte `json:"secret,omitempty"`
}
