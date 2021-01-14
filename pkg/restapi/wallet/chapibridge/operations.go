/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package chapibridge provides wallet server REST features related to credential handler APIs.
package chapibridge

import (
	"net/http"

	"github.com/hyperledger/aries-framework-go/pkg/controller/rest"
	"github.com/trustbloc/edge-core/pkg/log"

	"github.com/trustbloc/edge-agent/pkg/restapi/common"
)

var logger = log.New("wallet/chapi-bridge")

// constants for endpoints of wallet server CHAPI bridge controller.
const (
	CreateInvitationPath   = "/create-invitation"
	RequestCHAPIAppProfile = "/request-app-profile"
	SendCHAPIRequest       = "/send-chapi-request"
)

// Operation is REST service operation controller for CHAPI bridge features.
type Operation struct {
	handlers []rest.Handler
}

// Provider describes dependencies for this command.
type Provider interface {
}

// New returns new CHAPI bridge REST controller instance.
func New(ctx Provider) (*Operation, error) {
	o := &Operation{}

	o.registerHandler()

	return o, nil
}

// GetRESTHandlers get all controller API handler available for this protocol service.
func (c *Operation) GetRESTHandlers() []rest.Handler {
	return c.handlers
}

// registerHandler register handlers to be exposed from this protocol service as REST API endpoints.
func (c *Operation) registerHandler() {
	// Add more protocol endpoints here to expose them as REST controller endpoints
	c.handlers = []rest.Handler{
		common.NewHTTPHandler(CreateInvitationPath, http.MethodGet, c.CreateInvitation),
		common.NewHTTPHandler(RequestCHAPIAppProfile, http.MethodGet, c.RequestApplicationProfile),
		common.NewHTTPHandler(SendCHAPIRequest, http.MethodPost, c.SendCHAPIRequest),
	}
}

// CreateInvitation swagger:route GET /wallet/create-invitation chapi-bridge walletServerInvitationRequest
//
// Creates out-of-band invitation to connect to this wallet server.
// Response contains URL to application with invitation to startup.
//
// Responses:
//    default: genericError
//    200: createInvitationResponse
func (c *Operation) CreateInvitation(rw http.ResponseWriter, req *http.Request) {
	// TODO : to be implemented [#633]
	rw.WriteHeader(http.StatusNotImplemented)

	_, err := rw.Write([]byte("To be implemented !"))
	if err != nil {
		logger.Errorf("Unable to send error response, %s", err)
	}
}

// RequestApplicationProfile swagger:route GET /wallet/request-app-profile chapi-bridge requestAppProfile
//
// Requests wallet application profile of given user.
// Response contains wallet application profile of given user.
//
// Responses:
//    default: genericError
//    200: appProfileResponse
func (c *Operation) RequestApplicationProfile(rw http.ResponseWriter, req *http.Request) {
	// TODO : to be implemented [#633]
	rw.WriteHeader(http.StatusNotImplemented)

	_, err := rw.Write([]byte("To be implemented !"))
	if err != nil {
		logger.Errorf("Unable to send error response, %s", err)
	}
}

// SendCHAPIRequest swagger:route POST /wallet/send-chapi-request chapi-bridge chapiRequest
//
// Sends CHAPI request to given wallet application ID.
// Response contains CHAPI request.
//
// Responses:
//    default: genericError
//    200: chapiResponse
func (c *Operation) SendCHAPIRequest(rw http.ResponseWriter, req *http.Request) {
	// TODO : to be implemented [#633]
	rw.WriteHeader(http.StatusNotImplemented)

	_, err := rw.Write([]byte("To be implemented !"))
	if err != nil {
		logger.Errorf("Unable to send error response, %s", err)
	}
}
