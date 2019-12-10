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

type didexchangeClient interface {
	RegisterMsgEvent(ch chan<- service.StateMsg) error
	HandleInvitation(invitation *didexchange.Invitation) (string, error)
}

// RegisterHandleInvitationJSCallback register handle invitation js callback
func RegisterHandleInvitationJSCallback(didexchangeClient didexchangeClient) error {
	statusCh := make(chan service.StateMsg)
	if err := didexchangeClient.RegisterMsgEvent(statusCh); err != nil {
		return fmt.Errorf("failed to register msg event: %w", err)
	}

	c := callback{didexchangeClient: didexchangeClient, statusCh: statusCh, jsDoc: js.Global().Get("document")}
	js.Global().Set("acceptInvitation", js.FuncOf(c.handleInvitation))

	return nil
}

type callback struct {
	didexchangeClient didexchangeClient
	statusCh          chan service.StateMsg
	jsDoc             js.Value
}

func (c *callback) handleInvitation(this js.Value, inputs []js.Value) interface{} {
	// https://github.com/golang/go/issues/26382
	go func() {
		invitationData := c.jsDoc.
			Call("getElementById", inputs[0].String()).Get("value").String()
		if invitationData == "" {
			js.Global().Call("alert", "invitation data is empty")
			return
		}
		var invitation *didexchange.Invitation
		if err := json.Unmarshal([]byte(invitationData), &invitation); err != nil {
			js.Global().Call("alert", fmt.Errorf("failed to unmarshal invitation: %s", err).Error())
			return
		}
		fmt.Printf("handle invitation %s\n", invitationData)
		connID, err := c.didexchangeClient.HandleInvitation(invitation)
		if err != nil {
			js.Global().Call("alert", fmt.Errorf("failed to handle invitation: %s", err).Error())
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
