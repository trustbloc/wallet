/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package mocks provides useful mocks for testing different command controller features.
package mocks

// NewMockNotifier returns mock notifier implementation.
func NewMockNotifier() *Notifier {
	return &Notifier{}
}

// Notifier is mock implementation of notifier.
type Notifier struct {
	NotifyFunc func(topic string, message []byte) error
}

// Notify is mock implementation of notifier Notify().
func (n *Notifier) Notify(topic string, message []byte) error {
	if n.NotifyFunc != nil {
		return n.NotifyFunc(topic, message)
	}

	return nil
}
