/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package controller

import (
	"fmt"

	ariesstorage "github.com/hyperledger/aries-framework-go/pkg/storage"

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
func GetCommandHandlers(storeProvider ariesstorage.Provider, opts ...Opt) ([]command.Handler, error) {
	cmdOpts := &allOpts{}
	// Apply options
	for _, opt := range opts {
		opt(cmdOpts)
	}

	// did client command operation
	didClientCmd, err := didclientcmd.New(storeProvider, cmdOpts.blocDomain)
	if err != nil {
		return nil, fmt.Errorf("failed initialized didclient command: %w", err)
	}

	var allHandlers []command.Handler
	allHandlers = append(allHandlers, didClientCmd.GetHandlers()...)

	return allHandlers, nil
}
