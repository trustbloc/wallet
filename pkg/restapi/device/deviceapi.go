/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package device

import (
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/store/user"
)

// Device represents the device model.
type Device struct {
	id          string
	name        string
	displayName string
	credentials []webauthn.Credential
}

// NewDevice creates and returns a new Device.
func NewDevice(usr *user.User) *Device {
	device := &Device{}
	device.id = usr.Sub
	device.name = usr.Name
	device.displayName = usr.FamilyName
	device.credentials = []webauthn.Credential{}

	return device
}

// WebAuthnID returns the user's ID.
func (d Device) WebAuthnID() []byte {
	return []byte(d.id)
}

// WebAuthnName returns the user's username.
func (d Device) WebAuthnName() string {
	return d.name
}

// WebAuthnDisplayName returns the user's display name.
func (d Device) WebAuthnDisplayName() string {
	return d.displayName
}

// WebAuthnIcon is not (yet) implemented.
func (d Device) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user.
// nolint:gocritic
func (d *Device) AddCredential(cred webauthn.Credential) {
	d.credentials = append(d.credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user.
func (d Device) WebAuthnCredentials() []webauthn.Credential {
	return d.credentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled.
// with all the user's credentials.
func (d Device) CredentialExcludeList() []protocol.CredentialDescriptor {
	credentialExcludeList := []protocol.CredentialDescriptor{}

	for _, cred := range d.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}
