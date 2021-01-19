/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package protocol

import (
	"sync"

	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/mock/didcomm/protocol"
	mocksvc "github.com/hyperledger/aries-framework-go/pkg/mock/didcomm/service"
	"github.com/hyperledger/aries-framework-go/pkg/storage"
)

// MockProvider mock implementation of provider needed by sdk command controller.
type MockProvider struct {
	*protocol.MockProvider
	StoreProvider        storage.Provider
	ServiceEndpointValue string
	CustomMessenger      service.Messenger
}

// NewMockProvider returns mock implementation of basic provider.
func NewMockProvider() *MockProvider {
	return &MockProvider{MockProvider: &protocol.MockProvider{}}
}

// ServiceEndpoint returns the service endpoint.
func (p *MockProvider) ServiceEndpoint() string {
	return p.ServiceEndpointValue
}

// StorageProvider returns the storage provider.
func (p *MockProvider) StorageProvider() storage.Provider {
	return p.StoreProvider
}

// Messenger return mock messenger.
func (p *MockProvider) Messenger() service.Messenger {
	if p.CustomMessenger != nil {
		return p.CustomMessenger
	}

	return &mocksvc.MockMessenger{}
}

// NewMockMessenger returns new mock messenger.
func NewMockMessenger() *MockMessenger {
	return &MockMessenger{MockMessenger: &mocksvc.MockMessenger{}}
}

// MockMessenger mock implementation of messenger.
type MockMessenger struct {
	*mocksvc.MockMessenger
	lastID string
	lock   sync.RWMutex
}

// Send mock messenger Send.
func (m *MockMessenger) Send(msg service.DIDCommMsgMap, myDID, theirDID string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.lastID = msg.ID()

	return nil
}

// ReplyToNested mock messenger ReplyToNested.
func (m *MockMessenger) ReplyToNested(msg service.DIDCommMsgMap, opts *service.NestedReplyOpts) error {
	if m.ErrReplyToNested != nil {
		return m.ErrReplyToNested
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	m.lastID = msg.ID()

	return nil
}

// GetLastID returns ID of the last message received.
func (m *MockMessenger) GetLastID() string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.lastID
}
