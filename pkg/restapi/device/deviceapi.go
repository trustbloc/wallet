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
	ID          string
	Name        string
	DisplayName string
	Credentials []webauthn.Credential
}

// NewDevice creates and returns a new Device.
func NewDevice(usr *user.User) *Device {
	device := &Device{}
	device.ID = usr.Sub
	device.Name = usr.Name
	device.DisplayName = usr.FamilyName
	device.Credentials = []webauthn.Credential{}

	return device
}

// WebAuthnID returns the user's ID.
func (d Device) WebAuthnID() []byte {
	return []byte(d.ID)
}

// WebAuthnName returns the user's username.
func (d Device) WebAuthnName() string {
	return d.Name
}

// WebAuthnDisplayName returns the user's display Name.
func (d Device) WebAuthnDisplayName() string {
	return d.DisplayName
}

// WebAuthnIcon is not (yet) implemented.
func (d Device) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user.
// nolint:gocritic
func (d *Device) AddCredential(cred webauthn.Credential) {
	d.Credentials = append(d.Credentials, cred)
}

// WebAuthnCredentials returns Credentials owned by the user.
func (d Device) WebAuthnCredentials() []webauthn.Credential {
	return d.Credentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled.
// with all the user's Credentials.
func (d Device) CredentialExcludeList() []protocol.CredentialDescriptor {
	credentialExcludeList := []protocol.CredentialDescriptor{}

	for _, cred := range d.Credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}
