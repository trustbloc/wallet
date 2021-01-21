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

	"github.com/gorilla/mux"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofband"
	"github.com/hyperledger/aries-framework-go/pkg/controller/command"
	"github.com/hyperledger/aries-framework-go/pkg/controller/rest"
	"github.com/hyperledger/aries-framework-go/pkg/controller/webnotifier"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	didexchangesvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/hyperledger/aries-framework-go/pkg/storage"
	"github.com/trustbloc/edge-core/pkg/log"

	"github.com/trustbloc/edge-agent/pkg/restapi/common"
	"github.com/trustbloc/edge-agent/pkg/restapi/common/logutil"
)

var logger = log.New("wallet/chapi-bridge")

// constants for endpoints of wallet server CHAPI bridge controller.
const (
	commandName            = "/wallet"
	CreateInvitationPath   = "/create-invitation"
	RequestCHAPIAppProfile = "/{id}/request-app-profile"
	SendCHAPIRequest       = "/send-chapi-request"

	invalidIDErr = "invalid ID"

	_actions = "_actions"
	_states  = "_states"
)

// Operation is REST service operation controller for CHAPI bridge features.
type Operation struct {
	handlers     []rest.Handler
	agentLabel   string
	walletAppURL string
	store        *walletAppProfileStore
	outOfBand    *outofband.Client
	didExchange  *didexchange.Client
}

// Provider describes dependencies for this command.
type Provider interface {
	ServiceEndpoint() string
	Service(id string) (interface{}, error)
	KMS() kms.KeyManager
	StorageProvider() storage.Provider
	ProtocolStateStorageProvider() storage.Provider
}

// New returns new CHAPI bridge REST controller instance.
func New(p Provider, notifier command.Notifier, defaultLabel, walletAppURL string) (*Operation, error) {
	store, err := newWalletAppProfileStore(p.StorageProvider())
	if err != nil {
		return nil, fmt.Errorf("failed to open wallet profile store : %w", err)
	}

	outOfBandClient, err := outofband.New(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create out-of-band client : %w", err)
	}

	didExchangeClient, err := didexchange.New(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create did-exchange client : %w", err)
	}

	o := &Operation{
		agentLabel:   defaultLabel,
		walletAppURL: walletAppURL,
		store:        store,
		outOfBand:    outOfBandClient,
		didExchange:  didExchangeClient,
	}

	err = o.setupEventHandlers(notifier)
	if err != nil {
		return nil, fmt.Errorf("failed to register events : %w", err)
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
		common.NewHTTPHandler(CreateInvitationPath, http.MethodPost, o.CreateInvitation),
		common.NewHTTPHandler(RequestCHAPIAppProfile, http.MethodGet, o.RequestApplicationProfile),
		common.NewHTTPHandler(SendCHAPIRequest, http.MethodPost, o.SendCHAPIRequest),
	}
}

// CreateInvitation swagger:route POST /wallet/create-invitation chapi-bridge createInvitation
//
// Creates out-of-band invitation to connect to this wallet server.
// Response contains URL to application with invitation to load during startup.
//
// Responses:
//    default: genericError
//    200: createInvitationResponse
func (o *Operation) CreateInvitation(rw http.ResponseWriter, req *http.Request) {
	var request createInvitationRequest

	err := json.NewDecoder(req.Body).Decode(&request.Body)
	if err != nil {
		logutil.LogError(logger, commandName, CreateInvitationPath, err.Error())
		common.WriteErrorResponsef(rw, logger, http.StatusBadRequest, err.Error())

		return
	}

	if request.Body.UserID == "" {
		logutil.LogError(logger, commandName, CreateInvitationPath, invalidIDErr)
		common.WriteErrorResponsef(rw, logger, http.StatusBadRequest, invalidIDErr)

		return
	}

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

	err = o.store.SaveProfile(request.Body.UserID, &walletAppProfile{InvitationID: invitation.ID})
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

// RequestApplicationProfile swagger:route GET /wallet/{id}/request-app-profile chapi-bridge applicationProfileRequest
//
// Requests wallet application profile of given user.
// Response contains wallet application profile of given user.
//
// Responses:
//    default: genericError
//    200: appProfileResponse
func (o *Operation) RequestApplicationProfile(rw http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["id"]
	if userID == "" {
		logutil.LogError(logger, commandName, RequestCHAPIAppProfile, invalidIDErr)
		common.WriteErrorResponsef(rw, logger, http.StatusBadRequest, invalidIDErr)

		return
	}

	profile, err := o.store.GetProfileByUserID(userID)
	if err != nil {
		logutil.LogError(logger, commandName, RequestCHAPIAppProfile, err.Error())
		common.WriteErrorResponsef(rw, logger, http.StatusInternalServerError, err.Error())

		return
	}

	var status string
	if profile.ConnectionID != "" {
		status = didexchangesvc.StateIDCompleted
	}

	common.WriteResponse(rw, logger, &applicationProfileResponse{profile.InvitationID, status})

	logutil.LogInfo(logger, commandName, CreateInvitationPath, "sent profile info successfully")
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

func (o *Operation) setupEventHandlers(notifier command.Notifier) error {
	// creates action channel
	actions := make(chan service.DIDCommAction)
	// registers action channel to listen for events
	if err := o.didExchange.RegisterActionEvent(actions); err != nil {
		return fmt.Errorf("register action event: %w", err)
	}

	// create state channel subscribers
	states := []chan service.StateMsg{
		make(chan service.StateMsg), make(chan service.StateMsg),
	}

	// registers state channels to listen for events
	for _, state := range states {
		if err := o.didExchange.RegisterMsgEvent(state); err != nil {
			return fmt.Errorf("register msg event: %w", err)
		}
	}

	subscribers := []chan service.DIDCommAction{
		make(chan service.DIDCommAction), make(chan service.DIDCommAction),
	}

	go service.AutoExecuteActionEvent(subscribers[1])
	go o.stateMsgListener(states[1])

	go func() {
		for action := range actions {
			for i := range subscribers {
				action.Message = action.Message.Clone()
				subscribers[i] <- action
			}
		}
	}()

	obs := webnotifier.NewObserver(notifier)
	obs.RegisterAction(didexchangesvc.DIDExchange+_actions, subscribers[0])
	obs.RegisterStateMsg(didexchangesvc.DIDExchange+_states, states[0])

	return nil
}

func (o *Operation) stateMsgListener(ch <-chan service.StateMsg) {
	for msg := range ch {
		if msg.Type != service.PostState || msg.StateID != didexchangesvc.StateIDCompleted {
			continue
		}

		var event didexchange.Event

		switch p := msg.Properties.(type) {
		case didexchange.Event:
			event = p
		default:
			logger.Warnf("failed to cast didexchange event properties")

			continue
		}

		logger.Debugf(
			"Received connection complete event for invitationID=%s connectionID=%s",
			event.InvitationID(), event.ConnectionID())

		err := o.store.UpdateProfile(&walletAppProfile{
			InvitationID: event.InvitationID(),
			ConnectionID: event.ConnectionID(),
		})
		if err != nil {
			logger.Warnf("Failed to update wallet application profile: %w", err)
		}
	}
}
