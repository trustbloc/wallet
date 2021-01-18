/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chapibridge

import "encoding/json"

// createInvitationResponse model
//
//  Response of out-of-band invitation from wallet server.
//
// swagger:response createInvitationResponse
type createInvitationResponse struct {
	// in: body
	URL string `json:"url"`
}

// requestApplicationProfile model
//
// Request for querying wallet application profile ID for given user from wallet server.
//
// swagger:parameters requestAppProfile
type requestApplicationProfile struct { //nolint: unused,deadcode,gocritic // for open API only
	// in: body
	// required: true
	ID string `json:"id"`
}

// applicationProfileResponse model
//
// Response containing wallet application profile of user requested.
//
// swagger:response appProfileResponse
type applicationProfileResponse struct { //nolint: unused,deadcode,gocritic // for open API only
	// in: body
	ProfileID string `json:"profile_id"`

	// TODO: add more wallet application profile details [#633]
}

// chapiRequest model
//
// CHAPI request to be sent to given wallet application.
//
// swagger:parameters chapiRequest
type chapiRequest struct { //nolint: unused,deadcode,gocritic // for open API only
	// required: true
	// in: body
	Body struct {
		ProfileID string          `json:"profile_id"`
		Request   json.RawMessage `json:"chapi_request"`
	}
}

// chapiResponse model
//
// CHAPI response from requested wallet application.
//
// swagger:response chapiResponse
type chapiResponse struct { //nolint: unused,deadcode,gocritic // for open API only
	// in: body
	Body struct {
		ProfileID string `json:"profile_id"`

		Response json.RawMessage `json:"chapi_response"`
	}
}
