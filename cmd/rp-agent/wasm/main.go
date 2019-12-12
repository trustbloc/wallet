/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"syscall/js"

	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/decorator"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/ws"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/trustbloc/edge-agent/pkg/didexchange/invitation"
)

func main() {
	fmt.Println("rp wasm start")

	done := make(chan struct{})

	a, err := aries.New(aries.WithTransportReturnRoute(decorator.TransportReturnRouteAll),
		aries.WithOutboundTransports(ws.NewOutbound()))
	if err != nil {
		js.Global().Call("alert", err.Error())
	}

	ctx, err := a.Context()
	if err != nil {
		js.Global().Call("alert", err.Error())
	}

	fmt.Println("Instantiating DID Exchange protocol client")

	c, err := didexchange.New(ctx)
	if err != nil {
		js.Global().Call("alert", err.Error())
	}

	// Configure auto execute
	e := make(chan service.DIDCommAction)
	if err := c.RegisterActionEvent(e); err != nil {
		js.Global().Call("alert", err.Error())
	}

	go func() { service.AutoExecuteActionEvent(e) }()

	// register js callback
	if err := invitation.RegisterHandleInvitationJSCallback(c); err != nil {
		js.Global().Call("alert", err.Error())
	}

	<-done
}
