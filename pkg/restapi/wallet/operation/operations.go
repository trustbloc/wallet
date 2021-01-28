/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package operation provides wallet server REST features related to wallet operations.
package operation

import (
	"github.com/hyperledger/aries-framework-go/pkg/controller/command"
	"github.com/hyperledger/aries-framework-go/pkg/controller/rest"
)

// constants for endpoints of wallet server wallet  controller.
const ()

// Operation is REST service operation controller for wallet  features.
type Operation struct {
}

// Provider describes dependencies for this command.
type Provider interface {
}

// New returns new wallet  REST controller instance.
func New(p Provider, notifier command.Notifier, msgHandler command.MessageHandler) (*Operation, error) {
	return &Operation{}, nil
}

// GetRESTHandlers get all controller API handler available for this protocol service.
func (o *Operation) GetRESTHandlers() []rest.Handler {
	return []rest.Handler{}
}
