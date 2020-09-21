/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"context"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// OIDCProvider is the OIDC identity provider.
type OIDCProvider interface {
	Endpoint() oauth2.Endpoint
	Verifier(*oidc.Config) Verifier
}

// Verifier verifies id_tokens.
type Verifier interface {
	Verify(context.Context, string) (idToken, error)
}

// OIDCProviderImpl adapts an *oidc.Provider into an OIDCProvider.
type OIDCProviderImpl struct {
	OP *oidc.Provider
}

// Verifier returns a Verifier.
func (o *OIDCProviderImpl) Verifier(config *oidc.Config) Verifier {
	return &verifierImpl{v: o.OP.Verifier(config)}
}

// Endpoint returns the OIDC provider's endpoints.
func (o *OIDCProviderImpl) Endpoint() oauth2.Endpoint {
	return o.OP.Endpoint()
}

type verifierImpl struct {
	v *oidc.IDTokenVerifier
}

func (v *verifierImpl) Verify(ctx context.Context, token string) (idToken, error) {
	return v.v.Verify(ctx, token)
}

type idToken interface {
	Claims(interface{}) error
}
