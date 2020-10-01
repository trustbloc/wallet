/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/trustbloc/edge-core/pkg/storage"
	"net/http"
)

func openStore(p storage.Provider, name string) (storage.Store, error) {
	err := p.CreateStore(name)
	if err != nil && !errors.Is(err, storage.ErrDuplicateStore) {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return p.OpenStore(name)
}

// validation rules on the received user claims from the OIDC provider go here
func evaluateClaims(u *endUser) error {
	if u.Sub == "" {
		return fmt.Errorf("empty 'sub' in end user claims")
	}

	return nil
}

type cookieStoreAdapter struct {
	cookieStore *sessions.CookieStore
}

func (c *cookieStoreAdapter) Get(r *http.Request, name string) (Cookies, error) {
	session, err := c.cookieStore.Get(r, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session '%s': %w", name, err)
	}

	return &sessionsAdapter{session: session}, nil
}

type sessionsAdapter struct {
	session *sessions.Session
}

func (s *sessionsAdapter) Set(k interface{}, v interface{}) {
	s.session.Values[k] = v
}

func (s *sessionsAdapter) Get(k interface{}) (interface{}, bool) {
	v, found := s.session.Values[k]

	return v, found
}

func (s *sessionsAdapter) Delete(k interface{}) {
	delete(s.session.Values, k)
}

func (s *sessionsAdapter) Save(r *http.Request, w http.ResponseWriter) error {
	return s.session.Save(r, w)
}
