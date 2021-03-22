/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package device

import (
	"fmt"

	"github.com/duo-labs/webauthn/protocol"
)

const (
	appleAttestationKey = "apple"
)

// TODO implement full validation of apple attestation object format

// ValidateAppleAttestation validates an Apple attestation object.
func ValidateAppleAttestation(attObj protocol.AttestationObject, _ []byte, //nolint:gocritic //matches library interface
) (string, []interface{}, error) {
	x5c, x509present := attObj.AttStatement["x5c"].([]interface{})
	if !x509present {
		return "", nil, fmt.Errorf("missing x5c")
	}

	return appleAttestationKey, x5c, nil
}
