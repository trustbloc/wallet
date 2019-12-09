// +build js,wasm

/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package invitation

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
)

// RegisterHandleInvitationJSCallback register handle invitation js callback
func RegisterHandleInvitationJSCallback(didexchangeClient *didexchange.Client) {
	statusCh := make(chan service.StateMsg, 10)
	if err := didexchangeClient.RegisterMsgEvent(statusCh); err != nil {
		js.Global().Call("alert", err.Error())
	}

	c := callback{didexchangeClient: didexchangeClient, statusCh: statusCh}
	js.Global().Set("acceptInvitation", js.FuncOf(c.handleInvitation))
}

type callback struct {
	didexchangeClient *didexchange.Client
	statusCh          chan service.StateMsg
}

func (c *callback) handleInvitation(this js.Value, inputs []js.Value) interface{} {
	// https://github.com/golang/go/issues/26382
	go func() {
		invitationData := js.Global().Get("document").
			Call("getElementById", inputs[0].String()).Get("value").String()
		if invitationData == "" {
			js.Global().Call("alert", "invitation data is empty")
			return
		}
		var invitation *didexchange.Invitation
		if err := json.Unmarshal([]byte(invitationData), &invitation); err != nil {
			js.Global().Call("alert", err.Error())
			return
		}
		fmt.Printf("handle invitation %s\n", invitationData)
		connID, err := c.didexchangeClient.HandleInvitation(invitation)
		if err != nil {
			js.Global().Call("alert", err.Error())
			return
		}
		fmt.Printf("handle invitation is successful connenction id :%s\n", connID)
		completedFlag := make(chan struct{})
		errFlag := make(chan error)
		go msgEventListener(c.statusCh, connID, errFlag, completedFlag)
		select {
		case <-completedFlag:
			js.Global().Call("alert", "DID exchange is establish for connection ID:"+connID)
		case err = <-errFlag:
			js.Global().Call("alert", err.Error())
		case <-time.After(5 * time.Second):
			js.Global().Call("alert", "timeout waiting for did exchange completed status")
		}
	}()

	return nil
}

func msgEventListener(statusCh chan service.StateMsg, connID string, errFlag chan error, completedFlag chan struct{}) {
	for e := range statusCh {
		switch v := e.Properties.(type) {
		case didexchange.Event:
			if e.Type == service.PostState && e.StateID == "completed" && v.ConnectionID() == connID {
				close(completedFlag)
			}
		case error:
			errFlag <- fmt.Errorf("service processing failed : %w", v)
		}
	}
}
