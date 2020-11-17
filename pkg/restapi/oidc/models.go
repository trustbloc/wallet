/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package oidc

type createKeystoreReq struct {
	Controller         string `json:"controller,omitempty"`
	OperationalVaultID string `json:"operationalVaultID,omitempty"`
}
