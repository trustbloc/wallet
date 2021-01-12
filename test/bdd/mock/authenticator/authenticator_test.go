/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package authenticator_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/stretchr/testify/require"

	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
	"github.com/trustbloc/edge-agent/pkg/restapi/device"
	"github.com/trustbloc/edge-agent/test/bdd/mock/authenticator"
)

func TestAuthenticate(t *testing.T) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "rp display name",       // Display Name for your site
		RPID:          "rp-example.com",        // Generally the domain name for your site
		RPOrigin:      "http://rp-example.com", // The origin URL for WebAuthn requests
	})
	require.NoError(t, err)

	defaultUser := user.User{
		Sub:        "user_sub",
		Name:       "John Doe",
		GivenName:  "John",
		FamilyName: "Doe",
	}

	defaultDevice := device.NewDevice(&defaultUser)

	webAuthnUser := protocol.UserEntity{
		ID:          defaultDevice.WebAuthnID(),
		DisplayName: defaultDevice.WebAuthnDisplayName(),
		CredentialEntity: protocol.CredentialEntity{
			Name: defaultDevice.WebAuthnName(),
		},
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.User = webAuthnUser
		credCreationOpts.CredentialExcludeList = defaultDevice.CredentialExcludeList()
		credCreationOpts.Attestation = "direct"
	}

	t.Run("webauthn test flow", func(t *testing.T) {
		creationParams, registerSessionData, err := webAuthn.BeginRegistration(
			defaultDevice,
			registerOptions,
		)
		require.NoError(t, err)

		fmt.Printf("credParams: %#v\n", creationParams)

		ma := authenticator.New()

		authData, err := ma.Authenticate(
			"../../fixtures/keys/device/ec-cacert.pem",
			"../../fixtures/keys/device/ec-cakey.pem",
			webAuthn.Config.RPOrigin,
			creationParams,
		)
		require.NoError(t, err)

		fmt.Printf("AuthData: %s\n", string(authData))

		parsedResponse, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(authData))
		fmt.Printf("pR err: %#v\n", err)
		require.NoError(t, err)

		fmt.Printf("parsedResponse: %#v\n", parsedResponse)

		credential, err := webAuthn.CreateCredential(defaultDevice, *registerSessionData, parsedResponse)
		fmt.Printf("cred err: %#v\n", err)
		require.NoError(t, err)

		fmt.Printf("credential: %#v\n", credential)

		defaultDevice.AddCredential(*credential)

		fmt.Printf("### LOGIN ###\n")

		// generate PublicKeyCredentialRequestOptions, session data
		assertionParams, loginSessionData, err := webAuthn.BeginLogin(defaultDevice)
		require.NoError(t, err)

		fmt.Printf("session: %#v\n", loginSessionData)

		assertData, err := ma.Assert(webAuthn.Config.RPOrigin, &(assertionParams.Response))
		require.NoError(t, err)

		fmt.Printf("AssertData: %s\n", string(assertData))

		parsedCRResponse, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(assertData))
		fmt.Printf("pR err: %#v\n", err)
		require.NoError(t, err)
		fmt.Printf("parsedCRResponse: %#v\n", parsedCRResponse)

		credential, err = webAuthn.ValidateLogin(defaultDevice, *loginSessionData, parsedCRResponse)
		fmt.Printf("cred err: %#v\n", err)
		require.NoError(t, err)

		fmt.Printf("credential: %#v\n", credential)
	})
}
