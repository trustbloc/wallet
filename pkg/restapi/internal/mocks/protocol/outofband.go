/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package protocol

import (
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/outofband"
)

// MockOobService is a mock of OobService interface.
type MockOobService struct {
	AcceptInvitationHandle func(*outofband.Invitation, string, []string) (string, error)
	SaveInvitationErr      error
}

// AcceptInvitation mock implementation.
func (m *MockOobService) AcceptInvitation(arg0 *outofband.Invitation, arg1 string, arg2 []string) (string, error) {
	if m.AcceptInvitationHandle != nil {
		return m.AcceptInvitationHandle(arg0, arg1, arg2)
	}

	return "", nil
}

// ActionContinue mock implementation.
func (m *MockOobService) ActionContinue(arg0 string, arg1 outofband.Options) error {
	return nil
}

// ActionStop mock implementation.
func (m *MockOobService) ActionStop(arg0 string, arg1 error) error {
	return nil
}

// Actions mock implementation.
func (m *MockOobService) Actions() ([]outofband.Action, error) {
	return []outofband.Action{}, nil
}

// RegisterActionEvent mock implementation.
func (m *MockOobService) RegisterActionEvent(arg0 chan<- service.DIDCommAction) error {
	return nil
}

// RegisterMsgEvent mock implementation.
func (m *MockOobService) RegisterMsgEvent(arg0 chan<- service.StateMsg) error {
	return nil
}

// SaveInvitation mock implementation.
func (m *MockOobService) SaveInvitation(arg0 *outofband.Invitation) error {
	if m.SaveInvitationErr != nil {
		return m.SaveInvitationErr
	}

	return nil
}

// UnregisterActionEvent mock implementation.
func (m *MockOobService) UnregisterActionEvent(arg0 chan<- service.DIDCommAction) error {
	return nil
}

// UnregisterMsgEvent mock implementation.
func (m *MockOobService) UnregisterMsgEvent(arg0 chan<- service.StateMsg) error {
	return nil
}
