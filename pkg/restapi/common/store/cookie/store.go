/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cookie

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	// StoreName is the name of the cookie store.
	StoreName = "edgeagent_wallet"
)

// Config the cookie store.
type Config struct {
	AuthKey []byte
	EncKey  []byte
	MaxAge  int
}

// NewStore returns a new CookieStore.
func NewStore(config *Config) *Jars {
	cs := sessions.NewCookieStore(config.AuthKey, config.EncKey)
	cs.MaxAge(config.MaxAge)

	return &Jars{cs: cs}
}

// Jars is a collection of cookie Jars.
type Jars struct {
	cs *sessions.CookieStore
}

// Open the Jar.
func (cs *Jars) Open(r *http.Request) (Jar, error) {
	s, err := cs.cs.Get(r, StoreName)

	// the session 's' is returned even when error is not nil - pass it back to the caller and let the caller decide.
	return &Session{s: s}, err
}

// Session is a Jar holding cookies.
type Session struct {
	s *sessions.Session
}

// Set the cookie.
func (s *Session) Set(k, v interface{}) {
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

// Save changes to the Jar.
func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {
	s.s.Options = &sessions.Options{
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/",
	}

	return s.s.Save(r, w)
}
