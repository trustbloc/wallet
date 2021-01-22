/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chapibridge

import (
	"encoding/json"
	"time"
)

// createInvitationRequest model
//
// Request for creating wallet server invitation.
//
// swagger:parameters createInvitation
type createInvitationRequest struct {
	// in: body
	Body struct {
		// required: true
		UserID string `json:"userID"`
	}
}

// createInvitationResponse model
//
//  Response of out-of-band invitation from wallet server.
//
// swagger:response createInvitationResponse
type createInvitationResponse struct {
	// in: body
	URL string `json:"url"`
}

// applicationProfileRequest model
//
// Request for querying wallet application profile ID for given user from wallet server.
//
// swagger:parameters applicationProfileRequest
type applicationProfileRequest struct { //nolint: unused,deadcode,gocritic // for open API only
	// in: path
	ID string `json:"id"`
}

// applicationProfileResponse model
//
// Response containing wallet application profile of user requested.
//
// swagger:response appProfileResponse
type applicationProfileResponse struct {
	// InvitationID of invitation used to create profile.
	// in: body
	InvitationID string `json:"invitationID"`

	// ConnectionStatus is DIDComm connection status of the profile.
	// in: body
	ConnectionStatus string `json:"status"`
}

// chapiRequest model
//
// CHAPI request to be sent to given wallet application.
//
// swagger:parameters chapiRequest
type chapiRequest struct {
	// required: true
	// in: body
	Body struct {
		// UserID of wallet application profile.
		UserID string `json:"userID"`
		// Request is credential handler request to be sent out.
		Request json.RawMessage `json:"chapiRequest"`
		// Timeout (in milliseconds) waiting for reply.
		Timeout time.Duration `json:"timeout,omitempty"`
	}
}

// chapiResponse model
//
// CHAPI response from requested wallet application.
//
// swagger:response chapiResponse
type chapiResponse struct {
	// in: body
	Body struct {
		Response json.RawMessage `json:"chapiResponse"`
	}
}
