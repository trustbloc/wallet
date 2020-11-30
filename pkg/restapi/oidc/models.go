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

// BootstrapData user bootsrap data.
type BootstrapData struct {
	UserEDVVaultURL   string `json:"edvVaultURL,omitempty"`
	OpsEDVVaultURL    string `json:"opsVaultURL,omitempty"`
	AuthzKeyStoreURL  string `json:"authzKeyStoreURL,omitempty"`
	OpsKeyStoreURL    string `json:"opsKeyStoreURL,omitempty"`
	EDVOpsKIDURL      string `json:"edvOpsKIDURL,omitempty"`
	EDVHMACKIDURL     string `json:"edvHMACKIDURL,omitempty"`
	UserEDVCapability string `json:"edvCapability,omitempty"`
}

type userBootstrapData struct {
	Data *BootstrapData `json:"data,omitempty"`
}

type secretRequest struct {
	Secret []byte `json:"secret,omitempty"`
}

type hubKMSHeader struct {
	secretShare string
	userSub     string
	accessToken string
}

type userConfig struct {
	Sub         string `json:"sub,omitempty"`
	SecretShare string `json:"walletSecretShare,omitempty"`
}
