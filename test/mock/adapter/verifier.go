/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/presentproof"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/decorator"
	presentproofsvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/presentproof"
	"github.com/hyperledger/aries-framework-go/pkg/doc/presexch"
	"github.com/hyperledger/aries-framework-go/spi/storage"
	"github.com/trustbloc/edge-core/pkg/log"
)

const (
	verifierHTML = "./templates/verifier.html"
)

var logger = log.New("edge-adapter/rp-operations")

type verifierApp struct {
	agent *didComm
	store storage.Store
}

func startVerifierApp(agent *didComm, router *mux.Router) error {
	prov := mem.NewProvider()

	store, err := prov.OpenStore("verifier")
	if err != nil {
		return fmt.Errorf("failed to create store : %w", err)
	}

	app := verifierApp{agent: agent, store: store}

	actionCh := make(chan service.DIDCommAction)

	err = agent.DIDExchClient.RegisterActionEvent(actionCh)
	if err != nil {
		return fmt.Errorf("failed to register action events on didexchange-client : %w", err)
	}

	err = agent.PresentProofClient.RegisterActionEvent(actionCh)
	if err != nil {
		return fmt.Errorf("failed to register action events on present-proof-client : %w", err)
	}

	go listenForDIDCommMsg(actionCh)

	router.HandleFunc("/verifier", app.login)
	router.HandleFunc("/waci-share", app.waci)
	router.HandleFunc("/waci-share/{id}", app.waciShareCallback)

	return nil
}

func (v *verifierApp) login(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, verifierHTML, nil)
}

func (v *verifierApp) waci(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// generate OOB invitation
	inv, err := v.agent.OOBClient.CreateInvitation(nil)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to create oob invitation : %s", err))

		return
	}

	invBytes, err := json.Marshal(inv)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to unmarshal invitation : %s", err))

		return
	}

	err = v.store.Put(inv.ID, invBytes)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to store invitation data : %s", err))

		return
	}

	callbackURL := os.Getenv(demoExternalURLEnvKey) + "/waci-share/" + inv.ID

	redirectURL := fmt.Sprintf("%s/waci?oob=%s&redirect=%s", r.FormValue("walletURL"),
		base64.URLEncoding.EncodeToString(invBytes),
		base64.URLEncoding.EncodeToString([]byte(callbackURL)))

	logger.Infof("waci redirect : url=%s oob-invitation=%s callbackURL=%s", redirectURL, string(invBytes), callbackURL)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (v *verifierApp) waciShareCallback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := v.store.Get(id)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to get invitation data : %s", err))

		return
	}

	loadTemplate(w, verifierHTML, map[string]interface{}{"Msg": "Successfully Received Presentation"})
}

func listenForDIDCommMsg(actionCh chan service.DIDCommAction) {
	for action := range actionCh {
		logger.Infof("received action message : type=%s", action.Message.Type())

		switch action.Message.Type() {
		case didexchange.RequestMsgType:
			action.Continue(nil)
		case presentproofsvc.ProposePresentationMsgTypeV2:
			continueArg := presentproof.WithRequestPresentation(&presentproof.RequestPresentation{
				Comment: "Request Presentation",
				RequestPresentationsAttach: []decorator.Attachment{
					{
						ID:       uuid.NewString(),
						MimeType: "application/json",
						Data: decorator.AttachmentData{JSON: struct {
							Challenge string                           `json:"challenge"`
							Domain    string                           `json:"domain"`
							PD        *presexch.PresentationDefinition `json:"presentation_definition"`
						}{
							Challenge: uuid.NewString(),
							Domain:    uuid.NewString(),
							PD: &presexch.PresentationDefinition{
								ID:   uuid.NewString(),
								Name: "Demo Verifier",
								InputDescriptors: []*presexch.InputDescriptor{
									{ID: uuid.NewString(), Schema: []*presexch.Schema{{URI: "https://w3id.org/citizenship#PermanentResidentCard"}}},
								},
							}},
						},
					},
				},
			})

			action.Continue(continueArg)
		case presentproofsvc.PresentationMsgTypeV2:
			action.Continue(nil)
		default:
			action.Stop(nil)
		}
	}
}
