/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package oidc

type createKeystoreReq struct {
	Controller         string `json:"controller,omitempty"`
	OperationalVaultID string `json:"operationalVaultID,omitempty"`
}

// TODO delete during completion of https://github.com/trustbloc/edge-agent/issues/489.
type todoDeleteThisModel struct {
	UserEDVVaultURL  string `json:"edvVaultURL"`
	OpsEDVVaultURL   string `json:"opsVaultURL"`
	AuthzKeyStoreURL string `json:"authKeyStoreURL"`
	OpsKeyStoreURL   string `json:"opsKeyStoreURL"`
}
