/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cookie

import "net/http"

// Store has session cookies.
type Store interface {
	Open(*http.Request) (CookieJar, error)
}

// CookieJar are contained within a CookieStore.
type CookieJar interface {
	Set(k interface{}, v interface{})
	Get(k interface{}) (interface{}, bool)
	Delete(k interface{})
	Save(*http.Request, http.ResponseWriter) error
}
