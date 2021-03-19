/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package device // nolint:testpackage // changing to different package requires exposing internal REST handlers

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	ariesmem "github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	mockstore "github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	"github.com/stretchr/testify/require"

	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/cookie"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
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

//nolint:lll //lines are too long
const testCredReqReal = `{
	"id":"AdTb0CC_J1vMOVywiVlOXalH-ZxohB5Z6_pNZh-eOaObbCBX_dN6pd-dxyddLHDJ7bkhkwcE9dsG5dnN_265nqvs",
	"rawId":"AdTb0CC_J1vMOVywiVlOXalH-ZxohB5Z6_pNZh-eOaObbCBX_dN6pd-dxyddLHDJ7bkhkwcE9dsG5dnN_265nqvs",
	"type":"public-key",
	"response":{
		"attestationObject":"o2NmbXRmcGFja2VkZ2F0dFN0bXSiY2FsZyZjc2lnWEgwRgIhAJ7qQwHFYvXMBu0T-7O_Y5Cx-OACw46p5mNmgcQAUUFKAiEA83FFc5Dwysa8bthsoWHps-Itdj2G2ijV4tHw-u2clFdoYXV0aERhdGFYxkmWDeWIDoxodDQXD2R2YFuP5K65ooYyx5lc87qDHZdjRWBTctitzgACNbzGCmSLCyXx8FUDAEIB1NvQIL8nW8w5XLCJWU5dqUf5nGiEHlnr-k1mH545o5tsIFf903ql353HJ10scMntuSGTBwT12wbl2c3_brmeq-ylAQIDJiABIVggw5JvLjCxCBS7AHFfnTWPFMxrwaSUZ4VxbWN4ayMU_zoiWCBYvaFdqTTOiF3lZlgjNAlB_oi2eRVq67PmMsdtJyCkZA",
		"clientDataJSON":"eyJ0eXBlIjoid2ViYXV0aG4uY3JlYXRlIiwiY2hhbGxlbmdlIjoiOXhYNk1VcGFtQVk0N29hbXNHYlVOOHBJR1dIOHI1d3B4WEFaUHF6Q0pOOCIsIm9yaWdpbiI6Imh0dHBzOi8vbG9jYWxob3N0OjkwOTgiLCJjcm9zc09yaWdpbiI6ZmFsc2V9"
	}
}`

//nolint:lll //lines are too long
const testCredReqApple = `{
	"id":"FOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmg",
	"type":"public-key",
	"rawId":"FOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmg",
	"response":{
		"clientDataJSON":"eyJjaGFsbGVuZ2UiOiItUmk1TlpUeko4YjZtdlczVFZTY0xvdEVvQUxmZ0JhMkJuNFlTYUlPYkhjIiwib3JpZ2luIjoiaHR0cHM6Ly93ZWJhdXRobi5pbyIsInR5cGUiOiJ3ZWJhdXRobi5jcmVhdGUifQ",
		"attestationObject":"pGhBdXRoRGF0YaVkcnBpZFggdKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBlZmxhZ3MYQWpzaWduX2NvdW50AGhhdHRfZGF0YaNmYWFndWlkUAAAAAAAAAAAAAAAAAAAAABtY3JlZGVudGlhbF9pZFhAFOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmmpwdWJsaWNfa2V5WE2lAQIDJiABIVgg_2qd0Pc1PDrBnveu7ePL8kcap_DXxU5N9O5gtVrB9iEiWCCHKgLltzMYVNLg7Wkl_PhbX_83jSS3QjhI2SRh7Ky8R2hleHRfZGF0YfZoYXV0aERhdGFYxHSm6pITyZwvdLIkkrMgz0AmKpTBqVCgOX8pJQtghB7wQQAAAAAAAAAAAAAAAAAAAAAAAAAAAEAU7Fyayo8s0Ie3Igu9s2Su0cx0oB6pIldhkNt4V3SRzkSebwqXprQmnSzOAuKUPDVLgKNmF71nbrXURQZuPISapQECAyYgASFYIP9qndD3NTw6wZ73ru3jy_JHGqfw18VOTfTuYLVawfYhIlgghyoC5bczGFTS4O1pJfz4W1__N40kt0I4SNkkYeysvEdjZm10ZWFwcGxlZ2F0dFN0bXSiY3NpZ1hGMEQCIH8iIcGY_n5BFck9RjivHIQxyUdnIVwS0kYOm0wozg5sAiBx-hlyPwUK3Efhh4OSX2sylUHXs6ev0FiBX0-JbUwjtmN4NWOBWQJTMIICTzCCATegAwIBAgIEEjbRfzANBgkqhkiG9w0BAQsFADAuMSwwKgYDVQQDEyNZdWJpY28gVTJGIFJvb3QgQ0EgU2VyaWFsIDQ1NzIwMDYzMTAgFw0xNDA4MDEwMDAwMDBaGA8yMDUwMDkwNDAwMDAwMFowMTEvMC0GA1UEAwwmWXViaWNvIFUyRiBFRSBTZXJpYWwgMjM5MjU3MzQxMDMyNDEwODcwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATTZakeXpng1bQ5wNmvu4f0BY5H3RKxRO2xTSsz-NNcFRPkDXnw-Zmr4jZxlZOBydwrB4WLgqxjR2IEzPc01q4hozswOTAiBgkrBgEEAYLECgIEFTEuMy42LjEuNC4xLjQxNDgyLjEuNTATBgsrBgEEAYLlHAIBAQQEAwIFIDANBgkqhkiG9w0BAQsFAAOCAQEAIhubs7JyJPE-vqMi8DUer0ZJZqNvcmmFfI4j-eUFtVJ13U5BIj5_JhEJFGnPkp-lJj5sx3aBskhtqvQfsc-r6FUI8T9nUPbIGyneYBtecgi7-mR25WSpHX1kq1JK0E67Ws4hixUm8XH4fN71I5joQyxQub8VeBl6tuu-MqvRdpM4OJwkuMl6zuPxvGFkdsr0LxNn3yko0CZVxjudPNCrabaZb-VzeIuZUvgCq0-UEVWxCdweIOxtJUIXWFfuq-GbR4pfJheGDTGdPkWmD8QGmDVpBWHczmQmiHUG10WXn4Bn2zFIgAtoMFje34jx1fXrvNjWMqRlN9jooxvQY4Rrfw"
	}
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
		config.Storage.Storage = &mockstore.MockStoreProvider{
			FailNamespace: user.StoreName,
		}
		_, err := New(config)
		require.Error(t, err)
	})
	t.Run("error if cannot open permanent store", func(t *testing.T) {
		expected := errors.New("test")
		config := config(t)
		config.Storage.Storage = &mockstore.MockStoreProvider{
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
					userSubCookieName: userSub,
				},
			},
		}
		usr := user.User{Sub: userSub}
		err = o.store.users.Save(&usr)
		require.NoError(t, err)

		userData, proceed := o.getUserData(w, newDeviceRegistrationRequest(), userSubCookieName)
		require.NotNil(t, userData)
		require.True(t, proceed)
		require.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("get user data - missing device user session cookie", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		w := httptest.NewRecorder()
		userData, proceed := o.getUserData(w, newDeviceRegistrationRequest(), "")
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
					userSubCookieName: userSub,
				},
			},
		}
		userData, proceed := o.getUserData(w, newDeviceRegistrationRequest(), userSubCookieName)
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
		o.getUserData(w, newDeviceRegistrationRequest(), "")
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
					userSubCookieName: userSub,
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
					userSubCookieName: userSub,
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
		testCases := []struct {
			name           string
			challenge      string
			webAuthNConfig *webauthn.Config
			requestBody    []byte
			expectedStatus int
		}{
			{
				"success - direct self attestation",
				challenge,
				&webauthn.Config{
					RPDisplayName: "https://webauthn.io",
					RPID:          "webauthn.io",
					RPOrigin:      "https://webauthn.io",
				},
				[]byte(testCredentialRequestBody),
				http.StatusOK,
			},
			{
				"success - 'apple' credential",
				challenge,
				&webauthn.Config{
					RPDisplayName: "https://webauthn.io",
					RPID:          "webauthn.io",
					RPOrigin:      "https://webauthn.io",
				},
				[]byte(testCredReqApple),
				http.StatusOK,
			},
			{
				"success - real data",
				"9xX6MUpamAY47oamsGbUN8pIGWH8r5wpxXAZPqzCJN8", &webauthn.Config{
					RPDisplayName: "trustbloc",
					RPID:          "localhost",
					RPOrigin:      "https://localhost:9098",
				},
				[]byte(testCredReqReal),
				http.StatusOK,
			},
		}

		for _, testCase := range testCases {
			testCase := testCase // pin

			t.Run(testCase.name, func(t *testing.T) {
				o, err := New(config(t))
				require.NoError(t, err)
				userSub := uuid.New().String()
				w := httptest.NewRecorder()
				o.store.cookies = &cookie.MockStore{
					Jar: &cookie.MockJar{
						Cookies: map[interface{}]interface{}{
							userSubCookieName: userSub,
						},
					},
				}
				r := newDeviceRegistrationRequest()

				o, sessionData := prepareRequest(t, userSub, testCase.challenge, testCase.webAuthNConfig, o)

				dummyServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
					writer.WriteHeader(http.StatusOK)
				}))

				defer func() { dummyServer.Close() }()

				mockHubAuth := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
					http.Redirect(writer, req, dummyServer.URL, http.StatusFound)
				}))

				defer func() { mockHubAuth.Close() }()

				o.hubAuthURL = mockHubAuth.URL

				err = o.store.session.SaveWebauthnSession(userSub, sessionData, r, w)
				require.NoError(t, err)

				r.Header = w.Header()
				r.Body = ioutil.NopCloser(bytes.NewReader(testCase.requestBody))

				o.finishRegistration(w, r)
				require.Equal(t, http.StatusOK, w.Code)
			})
		}
	})

	t.Run("finish registration - failed to create credential", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()

		w := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: userSub,
				},
			},
		}
		r := newDeviceRegistrationRequest()

		o, sessionData := prepareRequest(t, userSub, challenge, &webauthn.Config{
			RPDisplayName: "https://webauthn.io",
			RPID:          "webauthn.io",
		}, o)

		err = o.store.session.SaveWebauthnSession(userSub, sessionData, r, w)
		require.NoError(t, err)

		r.Header = w.Header()
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(testCredentialRequestBody)))

		o.finishRegistration(w, r)
		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), `Error validating origin`)
	})

	t.Run("finish registration - failed to get web auth session", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()

		w := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: userSub,
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
					userSubCookieName: userSub,
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

	t.Run("finish registration - failed direct attestation, hub-auth rejected", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()
		w := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: userSub,
				},
			},
		}
		r := newDeviceRegistrationRequest()

		o, sessionData := prepareRequestWithDefaults(t, userSub, o)

		mockHubAuth := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			writer.WriteHeader(http.StatusInternalServerError)
		}))

		defer func() { mockHubAuth.Close() }()

		o.hubAuthURL = mockHubAuth.URL

		err = o.store.session.SaveWebauthnSession(userSub, sessionData, r, w)
		require.NoError(t, err)

		r.Header = w.Header()
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(testCredentialRequestBody)))

		o.finishRegistration(w, r)
		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "failed response from hub-auth")
	})
	t.Run("finish registration - failed direct attestation, invalid cert", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()
		w := httptest.NewRecorder()
		o.store.cookies = &cookie.MockStore{
			Jar: &cookie.MockJar{
				Cookies: map[interface{}]interface{}{
					userSubCookieName: userSub,
				},
			},
		}
		r := newDeviceRegistrationRequest()

		o, sessionData := prepareRequestWithDefaults(t, userSub, o)

		err = o.store.session.SaveWebauthnSession(userSub, sessionData, r, w)
		require.NoError(t, err)

		r.Header = w.Header()
		r.Body = ioutil.NopCloser(bytes.NewReader(generateBadCredentialResponse(t)))

		o.finishRegistration(w, r)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestOperation_requestDeviceValidation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()

		dummyServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			writer.WriteHeader(http.StatusOK)
		}))

		defer func() { dummyServer.Close() }()

		mockHubAuth := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			http.Redirect(writer, req, dummyServer.URL, http.StatusFound)
		}))

		defer func() { mockHubAuth.Close() }()

		o.hubAuthURL = mockHubAuth.URL

		certASN := makeCert(t)
		certPEM := pem.EncodeToMemory(&pem.Block{Bytes: certASN, Type: "CERTIFICATE"})

		var certArray []interface{}
		certArray = append(certArray, certPEM)

		err = o.requestDeviceValidation(context.TODO(), userSub, "AAAAAAA", certArray)
		require.NoError(t, err)
	})

	t.Run("failure - missing cert", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)

		err = o.requestDeviceValidation(context.TODO(), "userSub", "AAAAAAA", nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "missing")
	})

	t.Run("failure - cert data isn't []byte", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)

		var certs []interface{}
		certs = append(certs, struct {
			a struct{}
		}{
			a: struct{}{},
		})

		err = o.requestDeviceValidation(context.TODO(), "userSub", "AAAAAAA", certs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "can't cast")
	})

	t.Run("failure - bad response from hub-auth", func(t *testing.T) {
		o, err := New(config(t))
		require.NoError(t, err)
		userSub := uuid.New().String()

		mockHubAuth := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			writer.WriteHeader(http.StatusForbidden)
		}))

		defer func() { mockHubAuth.Close() }()

		o.hubAuthURL = mockHubAuth.URL

		certASN := makeCert(t)
		certPEM := pem.EncodeToMemory(&pem.Block{Bytes: certASN, Type: "CERTIFICATE"})

		var certArray []interface{}
		certArray = append(certArray, certPEM)

		err = o.requestDeviceValidation(context.TODO(), userSub, "AAAAAAA", certArray)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed response from hub-auth/device endpoint")
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
				userSubCookieName: userSub,
				deviceCookieName:  userSub,
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
					userSubCookieName: userSub,
					deviceCookieName:  userSub,
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
	t.Run("begin login - failed to get user data", func(t *testing.T) {
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
				userSubCookieName: userSub,
				deviceCookieName:  userSub,
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
		o.saveCookie(w, r, userSub, deviceCookieName)

		o.finishLogin(w, r)
		require.Equal(t, http.StatusFound, w.Code)
	})
	t.Run("finish login - failed to get web auth login session", func(t *testing.T) {
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
					userSubCookieName: userSub,
					deviceCookieName:  userSub,
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
	config.Storage.Storage = &mockstore.MockStoreProvider{
		Store: &mockstore.MockStore{},
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
	config.Storage.Storage = &mockstore.MockStoreProvider{
		Store: &mockstore.MockStore{
			Store: make(map[string]mockstore.DBEntry),
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

func prepareRequestWithDefaults(t *testing.T, userSub string, o *Operation) (*Operation, *webauthn.SessionData) {
	return prepareRequest(t, userSub, challenge, &webauthn.Config{
		RPDisplayName: "https://webauthn.io",
		RPID:          "webauthn.io",
		RPOrigin:      "https://webauthn.io",
	}, o)
}

func prepareRequest(t *testing.T, userSub, challenge string, webauthnConf *webauthn.Config, o *Operation,
) (*Operation, *webauthn.SessionData) {
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

	o.webauthn = &webauthn.WebAuthn{
		Config: webauthnConf,
	}

	return o, sessionData
}

func config(t *testing.T) *Config {
	t.Helper()

	return &Config{
		Storage: &StorageConfig{
			Storage:      ariesmem.NewProvider(),
			SessionStore: ariesmem.NewProvider(),
		},
		Cookie: &cookie.Config{
			AuthKey: key(t),
			EncKey:  key(t),
			MaxAge:  900,
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

func makeCert(t *testing.T) []byte {
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		PublicKeyAlgorithm:    x509.ECDSA,
	}

	secret, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &secret.PublicKey, secret)
	require.NoError(t, err)

	certASN, err := asn1.Marshal(cert)
	require.NoError(t, err)

	return certASN
}

func Test_generateTestData(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/register/finish", bytes.NewReader([]byte(testCredentialRequestBody)))

	parsedResponse, err := protocol.ParseCredentialCreationResponse(r)
	require.NoError(t, err)

	parsedResponse.Response.AttestationObject.Format = "apple"

	data := marshalCCR(t, &parsedResponse.Raw, &parsedResponse.Response.AttestationObject)

	fmt.Printf("%s\n", string(data))
}

func marshalCCR(t *testing.T, out *protocol.CredentialCreationResponse, att *protocol.AttestationObject) []byte {
	attObjBytes, err := cbor.Marshal(att)
	require.NoError(t, err)

	out.AttestationResponse.AttestationObject = attObjBytes

	outBytes, err := json.Marshal(out)
	require.NoError(t, err)

	return outBytes
}

// generates a credential creation response without any x509 certificates.
func generateBadCredentialResponse(t *testing.T) []byte {
	wrongCertASN := makeCert(t)

	return marshalCCR(t, &protocol.CredentialCreationResponse{
		PublicKeyCredential: protocol.PublicKeyCredential{
			Credential: protocol.Credential{
				ID:   "FOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmg",
				Type: "public-key",
			},
			RawID:      fromb64("FOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmg=="), // nolint:lll // test data
			Extensions: nil,
		},
		AttestationResponse: protocol.AuthenticatorAttestationResponse{
			AuthenticatorResponse: protocol.AuthenticatorResponse{
				ClientDataJSON: fromb64("eyJjaGFsbGVuZ2UiOiItUmk1TlpUeko4YjZtdlczVFZTY0xvdEVvQUxmZ0JhMkJuNFlTYUlPYkhjIiwib3JpZ2luIjoiaHR0cHM6Ly93ZWJhdXRobi5pbyIsInR5cGUiOiJ3ZWJhdXRobi5jcmVhdGUifQ=="), // nolint:lll // test data
			},
		},
	}, &protocol.AttestationObject{
		AuthData: protocol.AuthenticatorData{
			RPIDHash: fromb64("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvA="),
			Flags:    0x41,
			Counter:  0x0,
			AttData: protocol.AttestedCredentialData{
				AAGUID:              fromb64("AAAAAAAAAAAAAAAAAAAAAA=="),
				CredentialID:        fromb64("FOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmg=="),                 // nolint:lll // test data
				CredentialPublicKey: fromb64("pSABIVgg/2qd0Pc1PDrBnveu7ePL8kcap/DXxU5N9O5gtVrB9iEiWCCHKgLltzMYVNLg7Wkl/PhbX/83jSS3QjhI2SRh7Ky8RwECAyY="), // nolint:lll // test data
			},
			ExtData: []uint8(nil),
		},
		RawAuthData: fromb64("dKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAQBTsXJrKjyzQh7ciC72zZK7RzHSgHqkiV2GQ23hXdJHORJ5vCpemtCadLM4C4pQ8NUuAo2YXvWdutdRFBm48hJqlAQIDJiABIVgg/2qd0Pc1PDrBnveu7ePL8kcap/DXxU5N9O5gtVrB9iEiWCCHKgLltzMYVNLg7Wkl/PhbX/83jSS3QjhI2SRh7Ky8Rw=="), // nolint:lll // test data
		Format:      "fido-u2f",
		AttStatement: map[string]interface{}{
			"sig": fromb64("MEQCIH8iIcGY/n5BFck9RjivHIQxyUdnIVwS0kYOm0wozg5sAiBx+hlyPwUK3Efhh4OSX2sylUHXs6ev0FiBX0+JbUwjtg=="), // nolint:lll // test data
			"x5c": []interface{}{wrongCertASN},
		},
	})
}

func fromb64(s string) []byte {
	b, _ := base64.StdEncoding.DecodeString(s) // nolint:errcheck // only used for test data

	return b
}
