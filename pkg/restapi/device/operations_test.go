/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package device // nolint:testpackage // changing to different package requires exposing internal REST handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
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

// The mock request body that will be passed by navigator.Credentials.create({publicKey: makeCredentialOptions.publicKey
// to registerFinishPath endpoint
// Sample from : https://github.com/duo-labs/webauthn/blob/master/protocol/attestation_test.go

//nolint:lll //lines are too long
const testCredentialRequestBody = `{ 
		"ID": "FOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmg",
		"rawId": "FOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmg",
		"response": {
			"attestationObject": "o2NmbXRoZmlkby11MmZnYXR0U3RtdKJjc2lnWEYwRAIgfyIhwZj-fkEVyT1GOK8chDHJR2chXBLSRg6bTCjODmwCIHH6GXI_BQrcR-GHg5JfazKVQdezp6_QWIFfT4ltTCO2Y3g1Y4FZAlMwggJPMIIBN6ADAgECAgQSNtF_MA0GCSqGSIb3DQEBCwUAMC4xLDAqBgNVBAMTI1l1YmljbyBVMkYgUm9vdCBDQSBTZXJpYWwgNDU3MjAwNjMxMCAXDTE0MDgwMTAwMDAwMFoYDzIwNTAwOTA0MDAwMDAwWjAxMS8wLQYDVQQDDCZZdWJpY28gVTJGIEVFIFNlcmlhbCAyMzkyNTczNDEwMzI0MTA4NzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABNNlqR5emeDVtDnA2a-7h_QFjkfdErFE7bFNKzP401wVE-QNefD5maviNnGVk4HJ3CsHhYuCrGNHYgTM9zTWriGjOzA5MCIGCSsGAQQBgsQKAgQVMS4zLjYuMS40LjEuNDE0ODIuMS41MBMGCysGAQQBguUcAgEBBAQDAgUgMA0GCSqGSIb3DQEBCwUAA4IBAQAiG5uzsnIk8T6-oyLwNR6vRklmo29yaYV8jiP55QW1UnXdTkEiPn8mEQkUac-Sn6UmPmzHdoGySG2q9B-xz6voVQjxP2dQ9sgbKd5gG15yCLv6ZHblZKkdfWSrUkrQTrtaziGLFSbxcfh83vUjmOhDLFC5vxV4GXq2674yq9F2kzg4nCS4yXrO4_G8YWR2yvQvE2ffKSjQJlXGO5080Ktptplv5XN4i5lS-AKrT5QRVbEJ3B4g7G0lQhdYV-6r4ZtHil8mF4YNMZ0-RaYPxAaYNWkFYdzOZCaIdQbXRZefgGfbMUiAC2gwWN7fiPHV9eu82NYypGU32OijG9BjhGt_aGF1dGhEYXRhWMR0puqSE8mcL3SyJJKzIM9AJiqUwalQoDl_KSULYIQe8EEAAAAAAAAAAAAAAAAAAAAAAAAAAABAFOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmqUBAgMmIAEhWCD_ap3Q9zU8OsGe967t48vyRxqn8NfFTk307mC1WsH2ISJYIIcqAuW3MxhU0uDtaSX8-Ftf_zeNJLdCOEjZJGHsrLxH",
			"clientDataJSON": "eyJjaGFsbGVuZ2UiOiItUmk1TlpUeko4YjZtdlczVFZTY0xvdEVvQUxmZ0JhMkJuNFlTYUlPYkhjIiwib3JpZ2luIjoiaHR0cHM6Ly93ZWJhdXRobi5pbyIsInR5cGUiOiJ3ZWJhdXRobi5jcmVhdGUifQ"
		},
		"type": "public-key"
	}`

// The example from https://github.com/duo-labs/webauthn/blob/master/protocol/assertion_test.go

//nolint:lll //lines are too long
const testCredentialResponseBody = `{
		"id":"AI7D5q2P0LS-Fal9ZT7CHM2N5BLbUunF92T8b6iYC199bO2kagSuU05-5dZGqb1SP0A0lyTWng",
		"rawId":"AI7D5q2P0LS-Fal9ZT7CHM2N5BLbUunF92T8b6iYC199bO2kagSuU05-5dZGqb1SP0A0lyTWng",
		"type":"public-key",
		"response":{
			"authenticatorData":"dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBFXJJiGa3OAAI1vMYKZIsLJfHwVQMANwCOw-atj9C0vhWpfWU-whzNjeQS21Lpxfdk_G-omAtffWztpGoErlNOfuXWRqm9Uj9ANJck1p6lAQIDJiABIVggKAhfsdHcBIc0KPgAcRyAIK_-Vi-nCXHkRHPNaCMBZ-4iWCBxB8fGYQSBONi9uvq0gv95dGWlhJrBwCsj_a4LJQKVHQ",
			"clientDataJSON":"eyJjaGFsbGVuZ2UiOiJFNFBUY0lIX0hmWDFwQzZTaWdrMVNDOU5BbGdlenROMDQzOXZpOHpfYzlrIiwibmV3X2tleXNfbWF5X2JlX2FkZGVkX2hlcmUiOiJkbyBub3QgY29tcGFyZSBjbGllbnREYXRhSlNPTiBhZ2FpbnN0IGEgdGVtcGxhdGUuIFNlZSBodHRwczovL2dvby5nbC95YWJQZXgiLCJvcmlnaW4iOiJodHRwczovL3dlYmF1dGhuLmlvIiwidHlwZSI6IndlYmF1dGhuLmdldCJ9",
			"signature":"MEUCIBtIVOQxzFYdyWQyxaLR0tik1TnuPhGVhXVSNgFwLmN5AiEAnxXdCq0UeAVGWxOaFcjBZ_mEZoXqNboY5IkQDdlWZYc"
			}
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
	t.Run("error if cannot open permanent store", func(t *testing.T) {
		expected := errors.New("test")
		config := config(t)
		config.Storage.Storage = &mockstore.Provider{
			ErrOpenStoreHandle: expected,
		}
		_, err := New(config)
		require.Error(t, err)
		require.True(t, errors.Is(err, expected))
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

		o, sessionData := prepareRequest(t, userSub, o)

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

		o, sessionData := prepareRequest(t, userSub, o)
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

func TestBeginLogin(t *testing.T) {
	o, err := New(config(t))
	require.NoError(t, err)

	webAuthConfig := &webauthn.Config{
		RPDisplayName: "test",
		RPID:          "test_121",
	}

	o.webauthn = &webauthn.WebAuthn{
		Config: webAuthConfig,
	}
	userSub := uuid.New().String()
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
	t.Run("begin login - success", func(t *testing.T) {
		w := httptest.NewRecorder()
		cred := webauthn.Credential{}
		// create device for user
		deviceInfo := NewDevice(&usr)
		deviceInfo.AddCredential(cred)

		err := o.saveDeviceInfo(deviceInfo)
		require.NoError(t, err)

		o.beginLogin(w, newDeviceLoginRequest(userSub))
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("begin login - failed to begin login no Credentials for user", func(t *testing.T) {
		w := httptest.NewRecorder()
		deviceInfo := NewDevice(&usr)

		err := o.saveDeviceInfo(deviceInfo)
		require.NoError(t, err)

		o.beginLogin(w, newDeviceLoginRequest(userSub))
		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "failed to begin login: Found no credentials for user")
	})
	t.Run("begin login - failed to get user data", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
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

		w := httptest.NewRecorder()

		o.beginLogin(w, newDeviceLoginRequest(userSub))
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "failed to get device data")
	})
	t.Run("finish login - failed to get user data", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		w := httptest.NewRecorder()

		o.beginLogin(w, newDeviceLoginRequest(userSub))
		require.Equal(t, http.StatusNotFound, w.Code)
		require.Contains(t, w.Body.String(), "missing device user session cookie")
	})
}

func TestFinishLogin(t *testing.T) {
	o, err := New(config(t))
	require.NoError(t, err)

	webAuthConfig := &webauthn.Config{
		RPDisplayName: "https://webauthn.io",
		RPID:          "webauthn.io",
		RPOrigin:      "https://webauthn.io",
	}

	o.webauthn = &webauthn.WebAuthn{
		Config: webAuthConfig,
	}
	userSub := uuid.New().String()
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
	t.Run("finish login - success", func(t *testing.T) {
		loginChallenge := "E4PTcIH_HfX1pC6Sigk1SC9NAlgeztN0439vi8z_c9k"
		w := httptest.NewRecorder()
		//nolint:lll //lines are too long
		credByteID, err := base64.RawURLEncoding.DecodeString("AI7D5q2P0LS-Fal9ZT7CHM2N5BLbUunF92T8b6iYC199bO2kagSuU05-5dZGqb1SP0A0lyTWng")
		require.NoError(t, err)
		//nolint:lll //lines are too long
		byteCredentialPubKey, err := base64.RawURLEncoding.DecodeString("pQMmIAEhWCAoCF-x0dwEhzQo-ABxHIAgr_5WL6cJceREc81oIwFn7iJYIHEHx8ZhBIE42L26-rSC_3l0ZaWEmsHAKyP9rgslApUdAQI")
		require.NoError(t, err)

		cred := webauthn.Credential{
			ID:        credByteID,
			PublicKey: byteCredentialPubKey,
		}
		// create device for user
		deviceInfo := NewDevice(&usr)
		deviceInfo.AddCredential(cred)

		r := newDeviceLoginRequest(userSub)

		err = o.saveDeviceInfo(deviceInfo)
		require.NoError(t, err)

		sessionData := &webauthn.SessionData{
			Challenge:        loginChallenge,
			UserID:           deviceInfo.WebAuthnID(),
			UserVerification: protocol.VerificationDiscouraged,
		}

		err = o.store.session.SaveWebauthnSession(userSub, sessionData, r, w)
		require.NoError(t, err)
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(testCredentialResponseBody)))

		o.finishLogin(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("finish login - failed to get web auth login sessionr", func(t *testing.T) {
		w := httptest.NewRecorder()
		deviceInfo := NewDevice(&usr)

		err := o.saveDeviceInfo(deviceInfo)
		require.NoError(t, err)

		o.finishLogin(w, newDeviceLoginRequest(userSub))
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "failed to get web auth login session: error unmarshaling data")
	})
	t.Run("finish login - failed to get device data", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
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

		w := httptest.NewRecorder()

		o.finishLogin(w, newDeviceLoginRequest(userSub))
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "failed to get device data")
	})
	t.Run("finish login - failed to get user data", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		w := httptest.NewRecorder()

		o.finishLogin(w, newDeviceLoginRequest(userSub))
		require.Equal(t, http.StatusNotFound, w.Code)
		require.Contains(t, w.Body.String(), "missing device user session cookie")
	})
}

func Test_GetDeviceInfo(t *testing.T) {
	o, err := New(config(t))
	require.NoError(t, err)
	config := config(t)
	config.Storage.Storage = &mockstore.Provider{
		Store: &mockstore.MockStore{
			ErrGet: errors.New("test"),
		},
	}

	userSub := uuid.New().String()
	usr := user.User{Sub: userSub}
	err = o.store.users.Save(&usr)
	require.NoError(t, err)
	t.Run("get device info - failed to get user", func(t *testing.T) {
		d, err := o.getDeviceInfo("")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get user")
		require.Nil(t, d)
	})
	t.Run("get device info - store error", func(t *testing.T) {
		d, err := o.getDeviceInfo(userSub)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to fetch device data")
		require.Nil(t, d)
	})
}

func Test_SaveDeviceInfo(t *testing.T) {
	config := config(t)
	config.Storage.Storage = &mockstore.Provider{
		Store: &mockstore.MockStore{
			Store: make(map[string][]byte),
		},
	}
	o, err := New(config)
	require.NoError(t, err)

	d := &Device{}

	t.Run("save device info - failed to save the device data", func(t *testing.T) {
		err := o.saveDeviceInfo(d)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to save the device data")
	})
}

func Test_JsonResponse(t *testing.T) {
	w := httptest.NewRecorder()
	jsonResponse(w, make(chan int), 0)
	require.Equal(t, w.Code, http.StatusInternalServerError)
}

func newDeviceRegistrationRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/register/device", nil)
}

func newDeviceLoginRequest(username string) *http.Request {
	return httptest.NewRequest(http.MethodGet, fmt.Sprintf("/login/begin?username=%s", username), nil)
}

func prepareRequest(t *testing.T, userSub string, o *Operation) (*Operation, *webauthn.SessionData) {
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
