/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package chapibridge provides wallet server REST features related to credential handler APIs.
package chapibridge

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyperledger/aries-framework-go/pkg/client/outofband"
	"github.com/hyperledger/aries-framework-go/pkg/controller/rest"
	didexchangesvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/trustbloc/edge-core/pkg/log"

	"github.com/trustbloc/edge-agent/pkg/restapi/common"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/logutil"
)

var logger = log.New("wallet/chapi-bridge")

// constants for endpoints of wallet server CHAPI bridge controller.
const (
	commandName            = "/wallet"
	CreateInvitationPath   = "/create-invitation"
	RequestCHAPIAppProfile = "/request-app-profile"
	SendCHAPIRequest       = "/send-chapi-request"
)

// Operation is REST service operation controller for CHAPI bridge features.
type Operation struct {
	handlers     []rest.Handler
	agentLabel   string
	walletAppURL string
	outOfBand    *outofband.Client
}

// Provider describes dependencies for this command.
type Provider interface {
	ServiceEndpoint() string
	Service(id string) (interface{}, error)
	KMS() kms.KeyManager
}

// New returns new CHAPI bridge REST controller instance.
func New(p Provider, defaultLabel, walletAppURL string) (*Operation, error) {
	outOfBandClient, err := outofband.New(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create out-of-band client : %w", err)
	}

	o := &Operation{
		agentLabel:   defaultLabel,
		walletAppURL: walletAppURL,
		outOfBand:    outOfBandClient,
	}

	o.registerHandler()

	return o, nil
}

// GetRESTHandlers get all controller API handler available for this protocol service.
func (o *Operation) GetRESTHandlers() []rest.Handler {
	return o.handlers
}

// registerHandler register handlers to be exposed from this protocol service as REST API endpoints.
func (o *Operation) registerHandler() {
	// Add more protocol endpoints here to expose them as REST controller endpoints
	o.handlers = []rest.Handler{
		common.NewHTTPHandler(CreateInvitationPath, http.MethodGet, o.CreateInvitation),
		common.NewHTTPHandler(RequestCHAPIAppProfile, http.MethodGet, o.RequestApplicationProfile),
		common.NewHTTPHandler(SendCHAPIRequest, http.MethodPost, o.SendCHAPIRequest),
	}
}

// CreateInvitation swagger:route GET /wallet/create-invitation chapi-bridge walletServerInvitationRequest
//
// Creates out-of-band invitation to connect to this wallet server.
// Response contains URL to application with invitation to load during startup.
//
// Responses:
//    default: genericError
//    200: createInvitationResponse
func (o *Operation) CreateInvitation(rw http.ResponseWriter, req *http.Request) {
	// TODO : public DIDs in request parameters - [Issue#645]
	invitation, err := o.outOfBand.CreateInvitation([]string{didexchangesvc.PIURI},
		outofband.WithLabel(o.agentLabel))
	if err != nil {
		logutil.LogError(logger, commandName, CreateInvitationPath, err.Error())
		common.WriteErrorResponsef(rw, logger, http.StatusInternalServerError, err.Error())

		return
	}

	invitationBytes, err := json.Marshal(invitation)
	if err != nil {
		logutil.LogError(logger, commandName, CreateInvitationPath, err.Error())
		common.WriteErrorResponsef(rw, logger, http.StatusInternalServerError, err.Error())

		return
	}

	common.WriteResponse(rw, logger, &createInvitationResponse{
		URL: fmt.Sprintf("%s?oob=%s", o.walletAppURL, base64.StdEncoding.EncodeToString(invitationBytes)),
	})

	logutil.LogInfo(logger, commandName, CreateInvitationPath, "created oob invitation successfully")
}

// RequestApplicationProfile swagger:route GET /wallet/request-app-profile chapi-bridge requestAppProfile
//
// Requests wallet application profile of given user.
// Response contains wallet application profile of given user.
//
// Responses:
//    default: genericError
//    200: appProfileResponse
func (o *Operation) RequestApplicationProfile(rw http.ResponseWriter, req *http.Request) {
	// TODO : to be implemented [#633]
	common.WriteErrorResponsef(rw, logger, http.StatusNotImplemented, "To be implemented !")
}

// SendCHAPIRequest swagger:route POST /wallet/send-chapi-request chapi-bridge chapiRequest
//
// Sends CHAPI request to given wallet application ID.
// Response contains CHAPI request.
//
// Responses:
//    default: genericError
//    200: chapiResponse
func (o *Operation) SendCHAPIRequest(rw http.ResponseWriter, req *http.Request) {
	// TODO : to be implemented [#633]
	common.WriteErrorResponsef(rw, logger, http.StatusNotImplemented, "To be implemented !")
}
