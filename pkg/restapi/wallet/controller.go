/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package wallet provides handlers for wallet server operations.
package wallet

import (
	"github.com/hyperledger/aries-framework-go/pkg/controller/command"
	"github.com/hyperledger/aries-framework-go/pkg/controller/rest"
	"github.com/hyperledger/aries-framework-go/pkg/controller/webnotifier"
	"github.com/hyperledger/aries-framework-go/pkg/framework/context"

	"github.com/trustbloc/edge-agent/pkg/restapi/wallet/chapibridge"
)

const wsPath = "/ws"

type allOpts struct {
	webhookURLs  []string
	defaultLabel string
	msgHandler   command.MessageHandler
	notifier     command.Notifier
	walletAppURL string
}

// Opt represents a controller option.
type Opt func(opts *allOpts)

// WithWebhookURLs is an option for setting up a webhook dispatcher which will notify clients of events.
func WithWebhookURLs(webhookURLs ...string) Opt {
	return func(opts *allOpts) {
		opts.webhookURLs = webhookURLs
	}
}

// WithNotifier is an option for setting up a notifier which will notify clients of events.
func WithNotifier(notifier command.Notifier) Opt {
	return func(opts *allOpts) {
		opts.notifier = notifier
	}
}

// WithDefaultLabel is an option allowing for the defaultLabel to be set.
func WithDefaultLabel(defaultLabel string) Opt {
	return func(opts *allOpts) {
		opts.defaultLabel = defaultLabel
	}
}

// WithMessageHandler is an option allowing for the message handler to be set.
func WithMessageHandler(handler command.MessageHandler) Opt {
	return func(opts *allOpts) {
		opts.msgHandler = handler
	}
}

// WithWalletAppURL is an option for setting up wallet APP URL for wallet server.
func WithWalletAppURL(walletApp string) Opt {
	return func(opts *allOpts) {
		opts.walletAppURL = walletApp
	}
}

// GetRESTHandlers gets all REST handlers provided by wallet controller.
func GetRESTHandlers(ctx *context.Provider, opts ...Opt) ([]rest.Handler, error) { //nolint:interfacer,gocritic
	restAPIOpts := &allOpts{}
	// Apply options
	for _, opt := range opts {
		opt(restAPIOpts)
	}

	notifier := restAPIOpts.notifier
	if notifier == nil {
		notifier = webnotifier.New(wsPath, restAPIOpts.webhookURLs)
	}

	// chapiBridge REST controller operations,
	chapiBridge, err := chapibridge.New(ctx, notifier, restAPIOpts.msgHandler,
		restAPIOpts.defaultLabel, restAPIOpts.walletAppURL)
	if err != nil {
		return nil, err
	}

	// create handlers from all REST operations.
	var allHandlers []rest.Handler
	allHandlers = append(allHandlers, chapiBridge.GetRESTHandlers()...)

	return allHandlers, nil
}
