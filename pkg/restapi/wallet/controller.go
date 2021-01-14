/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package wallet provides handlers for wallet server operations.
package wallet

import (
	"github.com/hyperledger/aries-framework-go/pkg/controller/rest"
	"github.com/hyperledger/aries-framework-go/pkg/framework/context"

	"github.com/trustbloc/edge-agent/pkg/restapi/wallet/chapibridge"
)

// GetRESTHandlers gets all REST handlers provided by wallet controller.
func GetRESTHandlers(ctx *context.Provider) ([]rest.Handler, error) { //nolint:interfacer,gocritic
	// chapiBridge REST controller operations,
	chapiBridge, err := chapibridge.New(ctx)
	if err != nil {
		return nil, err
	}

	// create handlers from all REST operations.
	var allHandlers []rest.Handler
	allHandlers = append(allHandlers, chapiBridge.GetRESTHandlers()...)

	return allHandlers, nil
}
