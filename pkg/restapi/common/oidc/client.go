/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// Provider provides discovery of OIDC provider endpoints and also verifies id_tokens.
type Provider interface {
	Endpoint() oauth2.Endpoint
	Verifier(*oidc.Config) Verifier
}

// ProviderAdapter adapts an *oidc.Provider into an OIDCProvider.
type ProviderAdapter struct {
	OP *oidc.Provider
}

// Endpoint returns the OIDC endpoints.
func (o *ProviderAdapter) Endpoint() oauth2.Endpoint {
	return o.OP.Endpoint()
}

// Verifier returns an OIDC verifier.
func (o *ProviderAdapter) Verifier(config *oidc.Config) Verifier {
	return &verifierAdapter{v: o.OP.Verifier(config)}
}

// Verifier parses and verifies a raw id_token.
type Verifier interface {
	Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}

type verifierAdapter struct {
	v *oidc.IDTokenVerifier
}

func (v *verifierAdapter) Verify(ctx context.Context, token string) (*oidc.IDToken, error) {
	return v.v.Verify(ctx, token)
}

type oauth2Config interface {
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
	Exchange(context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

type oauth2ConfigImpl struct {
	oc *oauth2.Config
}

func (o *oauth2ConfigImpl) AuthCodeURL(state string, options ...oauth2.AuthCodeOption) string {
	return o.oc.AuthCodeURL(state, options...)
}

func (o *oauth2ConfigImpl) Exchange(
	ctx context.Context, code string, options ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return o.oc.Exchange(ctx, code, options...)
}

// BasicClient for OIDC.
type BasicClient struct {
	provider             Provider
	oauth2ConfigSupplier func() oauth2Config
	clientID             string
	tlsConfig            *tls.Config
}

// Config defines configuration for oidc client.
type Config struct {
	TLSConfig    *tls.Config
	Provider     Provider
	CallbackURL  string
	ClientID     string
	ClientSecret string
	Scopes       []string
}

// NewClient returns new BasicClient instance.
func NewClient(config *Config) *BasicClient {
	return &BasicClient{
		provider: config.Provider,
		oauth2ConfigSupplier: func() oauth2Config {
			return &oauth2ConfigImpl{oc: &oauth2.Config{
				ClientID:     config.ClientID,
				ClientSecret: config.ClientSecret,
				Endpoint:     config.Provider.Endpoint(),
				RedirectURL:  config.CallbackURL,
				Scopes:       config.Scopes,
			}}
		},
		clientID:  config.ClientID,
		tlsConfig: config.TLSConfig,
	}
}

// FormatRequest returns a correctly-formatted OIDC request.
func (c *BasicClient) FormatRequest(state string) string {
	return c.oauth2ConfigSupplier().AuthCodeURL(state)
}

// Exchange the auth code for the OAuth2 token.
func (c *BasicClient) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.oauth2ConfigSupplier().Exchange(
		context.WithValue(
			ctx,
			oauth2.HTTPClient,
			&http.Client{Transport: &http.Transport{TLSClientConfig: c.tlsConfig}},
		),
		code,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	if !token.Valid() {
		return nil, fmt.Errorf("server returned an invalid token")
	}

	return token, nil
}

// VerifyIDToken parses the id_token within the OAuth2 token and verifies it.
func (c *BasicClient) VerifyIDToken(ctx context.Context, oauthToken OAuth2Token) (IDToken, error) {
	rawIDToken, found := oauthToken.Extra("id_token").(string)
	if !found {
		return nil, fmt.Errorf("missing id_token")
	}

	idToken, err := c.provider.Verifier(&oidc.Config{ClientID: c.clientID}).Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify id_token: %w", err)
	}

	return idToken, nil
}
