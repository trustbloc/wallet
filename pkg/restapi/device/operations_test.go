/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package device // nolint:testpackage // changing to different package requires exposing internal REST handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/cookie"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
	"github.com/trustbloc/edge-core/pkg/storage/memstore"
	"github.com/trustbloc/edge-core/pkg/storage/mockstore"
)

// The mock request body that will be passed by navigator.credentials.create({publicKey: makeCredentialOptions.publicKey
// to registerFinishPath endpoint
// Sample from : https://github.com/duo-labs/webauthn/blob/master/protocol/attestation_test.go

//nolint:lll //lines are too long
const testCredentialRequestBody = `{ 
		"id": "FOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmg",
		"rawId": "FOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmg",
		"response": {
			"attestationObject": "o2NmbXRoZmlkby11MmZnYXR0U3RtdKJjc2lnWEYwRAIgfyIhwZj-fkEVyT1GOK8chDHJR2chXBLSRg6bTCjODmwCIHH6GXI_BQrcR-GHg5JfazKVQdezp6_QWIFfT4ltTCO2Y3g1Y4FZAlMwggJPMIIBN6ADAgECAgQSNtF_MA0GCSqGSIb3DQEBCwUAMC4xLDAqBgNVBAMTI1l1YmljbyBVMkYgUm9vdCBDQSBTZXJpYWwgNDU3MjAwNjMxMCAXDTE0MDgwMTAwMDAwMFoYDzIwNTAwOTA0MDAwMDAwWjAxMS8wLQYDVQQDDCZZdWJpY28gVTJGIEVFIFNlcmlhbCAyMzkyNTczNDEwMzI0MTA4NzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABNNlqR5emeDVtDnA2a-7h_QFjkfdErFE7bFNKzP401wVE-QNefD5maviNnGVk4HJ3CsHhYuCrGNHYgTM9zTWriGjOzA5MCIGCSsGAQQBgsQKAgQVMS4zLjYuMS40LjEuNDE0ODIuMS41MBMGCysGAQQBguUcAgEBBAQDAgUgMA0GCSqGSIb3DQEBCwUAA4IBAQAiG5uzsnIk8T6-oyLwNR6vRklmo29yaYV8jiP55QW1UnXdTkEiPn8mEQkUac-Sn6UmPmzHdoGySG2q9B-xz6voVQjxP2dQ9sgbKd5gG15yCLv6ZHblZKkdfWSrUkrQTrtaziGLFSbxcfh83vUjmOhDLFC5vxV4GXq2674yq9F2kzg4nCS4yXrO4_G8YWR2yvQvE2ffKSjQJlXGO5080Ktptplv5XN4i5lS-AKrT5QRVbEJ3B4g7G0lQhdYV-6r4ZtHil8mF4YNMZ0-RaYPxAaYNWkFYdzOZCaIdQbXRZefgGfbMUiAC2gwWN7fiPHV9eu82NYypGU32OijG9BjhGt_aGF1dGhEYXRhWMR0puqSE8mcL3SyJJKzIM9AJiqUwalQoDl_KSULYIQe8EEAAAAAAAAAAAAAAAAAAAAAAAAAAABAFOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmqUBAgMmIAEhWCD_ap3Q9zU8OsGe967t48vyRxqn8NfFTk307mC1WsH2ISJYIIcqAuW3MxhU0uDtaSX8-Ftf_zeNJLdCOEjZJGHsrLxH",
			"clientDataJSON": "eyJjaGFsbGVuZ2UiOiItUmk1TlpUeko4YjZtdlczVFZTY0xvdEVvQUxmZ0JhMkJuNFlTYUlPYkhjIiwib3JpZ2luIjoiaHR0cHM6Ly93ZWJhdXRobi5pbyIsInR5cGUiOiJ3ZWJhdXRobi5jcmVhdGUifQ"
		},
		"type": "public-key"
	}`

// challenge which is associated with the above testCredentialRequestBody.
const challenge = "-Ri5NZTzJ8b6mvW3TVScLotEoALfgBa2Bn4YSaIObHc"

func TestNew(t *testing.T) {
	t.Run("returns an instance", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		require.NotNil(t, o)
	})

	t.Run("error if cannot open user store", func(t *testing.T) {
		config := config(t)
		config.Storage.Storage = &mockstore.Provider{
			FailNameSpace: user.StoreName,
		}
		_, err := New(config)
		require.Error(t, err)
	})
}

func TestOperation_GetRESTHandlers(t *testing.T) {
	o, err := New(config(t))
	require.NoError(t, err)

	require.NotEmpty(t, o.GetRESTHandlers())
}

func Test_GetUserData(t *testing.T) {
	t.Run("get user data", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()

		w := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					deviceCookieName: userSub,
				},
			},
		}
		usr := user.User{Sub: userSub}
		err = o.store.users.Save(&usr)
		require.NoError(t, err)

		userData, proceed := o.getUserData(w, newDeviceRegistrationRequest())
		require.NotNil(t, userData)
		require.True(t, proceed)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("get user data - missing device user session cookie", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		w := httptest.NewRecorder()
		userData, proceed := o.getUserData(w, newDeviceRegistrationRequest())
		require.Nil(t, userData)
		require.False(t, proceed)
		require.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("get user data failed - failed to read user session cookie", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		w := httptest.NewRecorder()
		userSub := uuid.New().String()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					deviceCookieName: userSub,
				},
			},
		}
		userData, proceed := o.getUserData(w, newDeviceRegistrationRequest())
		require.Nil(t, userData)
		require.False(t, proceed)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("error internal server error if cannot fetch the cookies", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		o.store.cookies = &cookie.MockStore{
			OpenErr: errors.New("test"),
		}
		w := httptest.NewRecorder()
		o.getUserData(w, newDeviceRegistrationRequest())
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRegistration_BeginRegistration(t *testing.T) {
	t.Run("begin registration", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()

		w := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					deviceCookieName: userSub,
				},
			},
		}
		usr := user.User{Sub: userSub}
		err = o.store.users.Save(&usr)
		require.NoError(t, err)

		config := &webauthn.Config{
			RPDisplayName: "test",
			RPID:          "test_121",
		}

		o.webauthn = &webauthn.WebAuthn{
			Config: config,
		}

		o.beginRegistration(w, newDeviceRegistrationRequest())
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("begin registration - failed to get user data", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()

		w := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					deviceCookieName: userSub,
				},
			},
		}

		config := &webauthn.Config{
			RPDisplayName: "test",
			RPID:          "test_121",
		}

		o.webauthn = &webauthn.WebAuthn{
			Config: config,
		}

		o.beginRegistration(w, newDeviceRegistrationRequest())
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRegistration_FinishRegistration(t *testing.T) {
	t.Run("finish registration", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()
		w := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					deviceCookieName: userSub,
				},
			},
		}
		r := newDeviceRegistrationRequest()

		o, sessionData := prepareFinishRegistration(t, userSub, o)

		err = o.store.session.SaveWebauthnSession(userSub, sessionData, r, w)
		require.NoError(t, err)

		r.Header = w.Header()
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(testCredentialRequestBody)))

		o.finishRegistration(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("finish registration - failed to create credential", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()

		w := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					deviceCookieName: userSub,
				},
			},
		}
		r := newDeviceRegistrationRequest()

		o, sessionData := prepareFinishRegistration(t, userSub, o)
		config := &webauthn.Config{
			RPDisplayName: "https://webauthn.io",
			RPID:          "webauthn.io",
		}
		o.webauthn.Config = config

		err = o.store.session.SaveWebauthnSession(userSub, sessionData, r, w)
		require.NoError(t, err)

		r.Header = w.Header()
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(testCredentialRequestBody)))

		o.finishRegistration(w, r)
		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), `{"errMessage":"failed to finish registration: Error validating origin"}`)
	})

	t.Run("finish registration - failed to get web auth session", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()

		w := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					deviceCookieName: userSub,
				},
			},
		}
		usr := user.User{Sub: userSub}
		err = o.store.users.Save(&usr)
		require.NoError(t, err)

		config := &webauthn.Config{
			RPDisplayName: "test",
			RPID:          "test_121",
		}

		o.webauthn = &webauthn.WebAuthn{
			Config: config,
		}

		o.finishRegistration(w, newDeviceRegistrationRequest())
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("finish registration - failed to get user data", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()

		w := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					deviceCookieName: userSub,
				},
			},
		}

		config := &webauthn.Config{
			RPDisplayName: "test",
			RPID:          "test_121",
		}

		o.webauthn = &webauthn.WebAuthn{
			Config: config,
		}

		o.finishRegistration(w, newDeviceRegistrationRequest())
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRegistration_MarshallResponse(t *testing.T) {
	t.Run("marshall response error ", func(t *testing.T) {
		value := make(chan int)
		_, err := json.Marshal(value)
		require.Error(t, err)
		require.EqualError(t, err, "json: unsupported type: chan int")
	})
}

func newDeviceRegistrationRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/register/device", nil)
}

func prepareFinishRegistration(t *testing.T, userSub string, o *Operation) (*Operation, *webauthn.SessionData) {
	usr := user.User{Sub: userSub}
	err := o.store.users.Save(&usr)
	require.NoError(t, err)

	device := NewDevice(&usr)

	sessionData := &webauthn.SessionData{
		Challenge:        challenge,
		UserID:           device.WebAuthnID(),
		UserVerification: protocol.VerificationDiscouraged,
	}
	o.webauthn = &webauthn.WebAuthn{}
	config := &webauthn.Config{
		RPDisplayName: "https://webauthn.io",
		RPID:          "webauthn.io",
		RPOrigin:      "https://webauthn.io",
	}

	o.webauthn = &webauthn.WebAuthn{
		Config: config,
	}

	return o, sessionData
}

func config(t *testing.T) *Config {
	t.Helper()

	return &Config{
		Storage: &StorageConfig{
			Storage:      memstore.NewProvider(),
			SessionStore: memstore.NewProvider(),
		},
		Keys: &KeyConfig{
			Auth: key(t),
			Enc:  key(t),
		},
	}
}

func key(t *testing.T) []byte {
	t.Helper()

	key := make([]byte, 32)

	n, err := rand.Reader.Read(key)
	require.NoError(t, err)
	require.Equal(t, 32, n)

	return key
}
