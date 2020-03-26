/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package controller

import (
	"github.com/trustbloc/edge-agent/pkg/controller/command"
	didclientcmd "github.com/trustbloc/edge-agent/pkg/controller/command/didclient"
)

type allOpts struct {
	blocDomain string
}

// Opt represents a controller option.
type Opt func(opts *allOpts)

// WithBlocDomain is an option allowing for the trustbloc domain to be set.
func WithBlocDomain(blocDomain string) Opt {
	return func(opts *allOpts) {
		opts.blocDomain = blocDomain
	}
}

// GetCommandHandlers returns all command handlers provided by controller.
func GetCommandHandlers(opts ...Opt) ([]command.Handler, error) {
	cmdOpts := &allOpts{}
	// Apply options
	for _, opt := range opts {
		opt(cmdOpts)
	}

	// did client command operation
	didClientCmd := didclientcmd.New(cmdOpts.blocDomain)

	var allHandlers []command.Handler
	allHandlers = append(allHandlers, didClientCmd.GetHandlers()...)

	return allHandlers, nil
}
