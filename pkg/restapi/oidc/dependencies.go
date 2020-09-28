/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"context"
	oidc2 "github.com/trustbloc/edge-agent/pkg/restapi/common/oidc"
	"golang.org/x/oauth2"
)

// OIDCClient is capable of formatting authorization requests, exchanging the token grant for an access_token
// and id_token, and verifying id_tokens.
type OIDCClient interface {
	FormatRequest(state string) string
	Exchange(c context.Context, code string) (*oauth2.Token, error)
	VerifyIDToken(c context.Context, oauthToken oidc2.OAuth2Token) (oidc2.IDToken, error)
}
