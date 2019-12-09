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
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
)

func TestRegisterHandleInvitationJSCallback(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		require.NoError(t, RegisterHandleInvitationJSCallback(&mockDidExchangeClient{handleInvitationValue: "conn1"}))
		require.True(t, js.Global().Get("acceptInvitation").Truthy())
	})

	t.Run("test error from register msg event", func(t *testing.T) {
		err := RegisterHandleInvitationJSCallback(&mockDidExchangeClient{
			registerMsgEventErr: fmt.Errorf("error")})
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to register msg event")
	})
}

func TestHandleInvitationSuccess(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		statusCh := make(chan service.StateMsg, 1)
		c := callback{didexchangeClient: &mockDidExchangeClient{handleInvitationValue: "conn1"},
			statusCh: statusCh}
		w := js.ValueOf("element1")
		js.Global().Set("getElementById", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			m := make(map[string]interface{})
			invitationData, err := json.Marshal(didexchange.Invitation{})
			require.NoError(t, err)
			m["value"] = js.ValueOf(string(invitationData))
			return m
		}))
		msg := make(chan string, 2)
		js.Global().Set("alert", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			msg <- args[0].String()
			return nil
		}))

		c.handleInvitation(js.Value{}, []js.Value{w})
		statusCh <- service.StateMsg{Properties: &mockDidExchangeEvent{connID: "conn1"},
			Type: service.PostState, StateID: "completed"}
		time.Sleep(500 * time.Millisecond)
		select {
		case m := <-msg:
			require.Contains(t, m, "DID exchange is establish for connection ID:conn1")
		case <-time.After(1 * time.Second):
			require.Fail(t, "timeout waiting for error")
		}
	})
}

func TestHandleInvitationError(t *testing.T) {
	t.Run("test empty invitation data", func(t *testing.T) {
		statusCh := make(chan service.StateMsg, 1)
		c := callback{didexchangeClient: &mockDidExchangeClient{}, statusCh: statusCh}
		w := js.ValueOf("element1")
		js.Global().Set("getElementById", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			m := make(map[string]interface{})
			m["value"] = js.ValueOf("")
			return m
		}))
		errFlag := make(chan string, 2)
		js.Global().Set("alert", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			errFlag <- args[0].String()
			return nil
		}))

		c.handleInvitation(js.Value{}, []js.Value{w})
		time.Sleep(500 * time.Millisecond)
		select {
		case err := <-errFlag:
			require.Contains(t, err, "invitation data is empty")
		case <-time.After(1 * time.Second):
			require.Fail(t, "timeout waiting for error")
		}
	})
	t.Run("test error unmarshal invitation data", func(t *testing.T) {
		statusCh := make(chan service.StateMsg, 1)
		c := callback{didexchangeClient: &mockDidExchangeClient{}, statusCh: statusCh}
		w := js.ValueOf("element1")
		js.Global().Set("getElementById", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			m := make(map[string]interface{})
			m["value"] = js.ValueOf("invalidInvitation")
			return m
		}))
		errFlag := make(chan string, 2)
		js.Global().Set("alert", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			errFlag <- args[0].String()
			return nil
		}))

		c.handleInvitation(js.Value{}, []js.Value{w})
		time.Sleep(500 * time.Millisecond)
		select {
		case err := <-errFlag:
			require.Contains(t, err, "failed to unmarshal invitation")
		case <-time.After(1 * time.Second):
			require.Fail(t, "timeout waiting for error")
		}
	})
	t.Run("test error handle invitation", func(t *testing.T) {
		statusCh := make(chan service.StateMsg, 1)
		c := callback{didexchangeClient: &mockDidExchangeClient{handleInvitationErr: fmt.Errorf("error")},
			statusCh: statusCh}
		w := js.ValueOf("element1")
		js.Global().Set("getElementById", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			m := make(map[string]interface{})
			invitationData, err := json.Marshal(didexchange.Invitation{})
			require.NoError(t, err)
			m["value"] = js.ValueOf(string(invitationData))
			return m
		}))
		errFlag := make(chan string, 2)
		js.Global().Set("alert", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			errFlag <- args[0].String()
			return nil
		}))

		c.handleInvitation(js.Value{}, []js.Value{w})
		time.Sleep(500 * time.Millisecond)
		select {
		case err := <-errFlag:
			require.Contains(t, err, "failed to handle invitation")
		case <-time.After(1 * time.Second):
			require.Fail(t, "timeout waiting for error")
		}
	})

	t.Run("test error from service processing", func(t *testing.T) {
		statusCh := make(chan service.StateMsg, 1)
		c := callback{didexchangeClient: &mockDidExchangeClient{handleInvitationValue: "conn1"},
			statusCh: statusCh}
		w := js.ValueOf("element1")
		js.Global().Set("getElementById", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			m := make(map[string]interface{})
			invitationData, err := json.Marshal(didexchange.Invitation{})
			require.NoError(t, err)
			m["value"] = js.ValueOf(string(invitationData))
			return m
		}))
		errFlag := make(chan string, 2)
		js.Global().Set("alert", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			errFlag <- args[0].String()
			return nil
		}))

		c.handleInvitation(js.Value{}, []js.Value{w})
		statusCh <- service.StateMsg{Properties: fmt.Errorf("error")}
		time.Sleep(500 * time.Millisecond)
		select {
		case err := <-errFlag:
			require.Contains(t, err, "service processing failed")
		case <-time.After(1 * time.Second):
			require.Fail(t, "timeout waiting for error")
		}
	})
}

type mockDidExchangeClient struct {
	handleInvitationErr   error
	handleInvitationValue string
	registerMsgEventErr   error
}

func (m *mockDidExchangeClient) RegisterMsgEvent(ch chan<- service.StateMsg) error {
	return m.registerMsgEventErr
}

func (m *mockDidExchangeClient) HandleInvitation(invitation *didexchange.Invitation) (string, error) {
	return m.handleInvitationValue, m.handleInvitationErr
}

type mockDidExchangeEvent struct {
	connID string
}

// connection ID
func (m *mockDidExchangeEvent) ConnectionID() string {
	return m.connID
}

// invitation ID
func (m *mockDidExchangeEvent) InvitationID() string {
	return ""
}
