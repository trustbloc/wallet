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

type OIDCProvider interface {
	Endpoint() oauth2.Endpoint
	Verifier(*oidc.Config) verifier
}

type verifier interface {
	Verify(context.Context, string) (idToken, error)
}

type oidcProviderImpl struct {
	op *oidc.Provider
}

func (o *oidcProviderImpl) Verifier(config *oidc.Config) verifier {
	return &verifierImpl{v: o.op.Verifier(config)}
}

type verifierImpl struct {
	v *oidc.IDTokenVerifier
}

func (v *verifierImpl) Verify(ctx context.Context, token string) (idToken, error) {
	return v.v.Verify(ctx, token)
}

func (o *oidcProviderImpl) Endpoint() oauth2.Endpoint {
	return o.op.Endpoint()
}

type idToken interface {
	Claims(interface{}) error
}
