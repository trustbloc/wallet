/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/trustbloc/edge-agent/pkg/restapi/common"
)

type signer interface {
	Sign(data []byte) ([]byte, error)
}

type kmsSigner struct {
	baseURL    string
	httpClient common.HTTPClient
	keystoreID string
	keyID      string
	header     *hubKMSHeader
}

func newKMSSigner(baseURL, keystoreID, keyID string, h *hubKMSHeader, httpClient common.HTTPClient) *kmsSigner {
	return &kmsSigner{
		baseURL:    baseURL,
		httpClient: httpClient,
		keystoreID: keystoreID,
		keyID:      keyID,
		header:     h,
	}
}

func (a *kmsSigner) Sign(data []byte) ([]byte, error) {
	reqBytes, err := json.Marshal(signReq{
		Message: base64.URLEncoding.EncodeToString(data),
	})
	if err != nil {
		return nil, fmt.Errorf("marshal create sign req : %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(),
		http.MethodPost, a.baseURL+fmt.Sprintf(signEndpoint, a.keystoreID, a.keyID), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}

	addAuthZKMSHeaders(req, a.header)

	resp, _, err := common.SendHTTPRequest(req, a.httpClient, http.StatusOK, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to sign from kms: %w", err)
	}

	var parsedResp signResp

	if errUnmarshal := json.Unmarshal(resp, &parsedResp); errUnmarshal != nil {
		return nil, fmt.Errorf("failed to unmarshal sign resp: %w", errUnmarshal)
	}

	signatureBytes, err := base64.URLEncoding.DecodeString(parsedResp.Signature)
	if err != nil {
		return nil, err
	}

	return signatureBytes, nil
}
