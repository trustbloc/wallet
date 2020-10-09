/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cookie

import "net/http"

// Store contains cookie jars.
type Store interface {
	Open(*http.Request) (Jar, error)
}

// Jar is a container of cookies from a Jars.
type Jar interface {
	Set(k interface{}, v interface{})
	Get(k interface{}) (interface{}, bool)
	Delete(k interface{})
	Save(*http.Request, http.ResponseWriter) error
}
