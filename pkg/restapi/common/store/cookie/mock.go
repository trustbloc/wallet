/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cookie

import "net/http"

// MockStore is a mock cookie Store.
type MockStore struct {
	Jar     CookieJar
	OpenErr error
}

// Open returns a CookieJar.
func (m *MockStore) Open(_ *http.Request) (CookieJar, error) {
	if m.Jar != nil || m.OpenErr != nil {
		return m.Jar, m.OpenErr
	}

	return &MockJar{}, nil
}

// MockJar is a mock CookieJar.
type MockJar struct {
	Cookies map[interface{}]interface{}
	SaveErr error
}

// Set the cookie.
func (m *MockJar) Set(k interface{}, v interface{}) {
	if m.Cookies == nil {
		m.Cookies = make(map[interface{}]interface{})
	}

	m.Cookies[k] = v
}

// Get the cookie.
func (m *MockJar) Get(k interface{}) (interface{}, bool) {
	v, ok := m.Cookies[k]

	return v, ok
}

// Delete the cookie.
func (m *MockJar) Delete(k interface{}) {
	delete(m.Cookies, k)
}

// Save changes to the CookieJar.
func (m *MockJar) Save(_ *http.Request, _ http.ResponseWriter) error {
	return m.SaveErr
}
