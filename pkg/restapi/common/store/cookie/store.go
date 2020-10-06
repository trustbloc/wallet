/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cookie

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	// StoreName is the name of the cookie store.
	StoreName = "edgeagent_wallet"
)

// NewStore returns a new CookieStore.
func NewStore(authKey, encKey []byte) *CookieStore {
	// TODO make session cookies max age configurable: https://github.com/trustbloc/edge-agent/issues/388
	cookieStore := sessions.NewCookieStore(authKey, encKey)
	cookieStore.MaxAge(900) // 15 mins

	return &CookieStore{cs: cookieStore}
}

// CookieStore holds CookieJars.
type CookieStore struct {
	cs *sessions.CookieStore
}

// Open opens the CookieJar.
func (cs *CookieStore) Open(r *http.Request) (CookieJar, error) {
	s, err := cs.cs.Get(r, StoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch session cookies %s: %w", StoreName, err)
	}

	return &Session{s: s}, nil
}

// Session is a CookieJar holding cookies.
type Session struct {
	s *sessions.Session
}

// Set the cookie.
func (s *Session) Set(k interface{}, v interface{}) {
	s.s.Values[k] = v
}

// Get the cookie.
func (s *Session) Get(k interface{}) (interface{}, bool) {
	v, found := s.s.Values[k]

	return v, found
}

// Delete the cookie.
func (s *Session) Delete(k interface{}) {
	delete(s.s.Values, k)
}

// Save changes to the CookieJar.
func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {
	return s.s.Save(r, w)
}
