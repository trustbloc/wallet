/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package oidc

type createKeyStoreReq struct {
	Controller string      `json:"controller"`
	EDV        *edvOptions `json:"edv"`
}

type edvOptions struct {
	VaultURL   string `json:"vault_url"`
	Capability []byte `json:"capability"`
}

type createKeyStoreResp struct {
	KeyStoreURL string `json:"key_store_url"`
	Capability  []byte `json:"capability"`
}

type createDIDResp struct {
	DID string `json:"did"`
}

type createKeyReq struct {
	KeyType string `json:"key_type"`
}

type createKeyResp struct {
	KeyURL    string `json:"key_url"`
	PublicKey []byte `json:"public_key"`
}

type signReq struct {
	Message []byte `json:"message"`
}

type signResp struct {
	Signature []byte `json:"signature"`
}

// BootstrapData user bootsrap data.
// TODO to be refactored for universal wallet migration.
type BootstrapData struct {
	User              string `json:"user,omitempty"`
	UserEDVVaultURL   string `json:"edvVaultURL,omitempty"` // TODO remove this
	OpsEDVVaultURL    string `json:"opsVaultURL,omitempty"` // TODO remove this
	AuthzKeyStoreURL  string `json:"authzKeyStoreURL,omitempty"`
	OpsKeyStoreURL    string `json:"opsKeyStoreURL,omitempty"`
	EDVOpsKIDURL      string `json:"edvOpsKIDURL,omitempty"`
	EDVHMACKIDURL     string `json:"edvHMACKIDURL,omitempty"`
	UserEDVCapability string `json:"edvCapability,omitempty"`
	OPSKMSCapability  string `json:"opsKMSCapability,omitempty"` // TODO remove this
	UserEDVServer     string `json:"userEDVServer,omitempty"`
	UserEDVVaultID    string `json:"userEDVVaultID,omitempty"`
	UserEDVEncKID     string `json:"userEDVEncKID,omitempty"`
	UserEDVMACKID     string `json:"userEDVMACKID,omitempty"`
	TokenExpiry       string `json:"tokenExpiry,omitempty"`
}

type userBootstrapData struct {
	Data *BootstrapData `json:"data,omitempty"`
}

type secretRequest struct {
	Secret []byte `json:"secret,omitempty"`
}

type kmsHeader struct {
	userSub     string
	accessToken string
	secretShare []byte
}

type userConfig struct {
	AccessToken string `json:"accessToken,omitempty"`
	SecretShare string `json:"walletSecretShare,omitempty"`
}
