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
	"github.com/hyperledger/aries-framework-go/pkg/common/log"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/decorator"
	presentproofsvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/presentproof"
	"github.com/hyperledger/aries-framework-go/pkg/doc/presexch"
	arieslog "github.com/hyperledger/aries-framework-go/spi/log"
	"github.com/hyperledger/aries-framework-go/spi/storage"
)

const (
	verifierHTML  = "./templates/verifier.html"
	webWalletHTML = "./templates/webWallet.html"
)

var logger = log.New("mock-adapter")

type verifierApp struct {
	agent *didComm
	store storage.Store
}

func startVerifierApp(agent *didComm, router *mux.Router) error {
	log.SetLevel("", arieslog.DEBUG)

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

	go listenForDIDCommMsg(actionCh, store)

	router.HandleFunc("/verifier", app.login)
	router.HandleFunc("/web-wallet", app.webWallet)
	router.HandleFunc("/waci-share", app.waci)
	router.HandleFunc("/waci-share/{id}", app.waciShareCallback)

	return nil
}

func (v *verifierApp) login(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, verifierHTML, nil)
}

func (v *verifierApp) webWallet(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, webWalletHTML, nil)
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

	redirectURL := fmt.Sprintf("%s/waci?oob=%s", r.FormValue("walletURL"),
		base64.URLEncoding.EncodeToString(invBytes))

	logger.Infof("waci redirect : url=%s oob-invitation=%s", redirectURL, string(invBytes))

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

func listenForDIDCommMsg(actionCh chan service.DIDCommAction, store storage.Store) {
	for action := range actionCh {
		logger.Infof("received action message : type=%s", action.Message.Type())

		switch action.Message.Type() {
		case didexchange.RequestMsgType:
			action.Continue(nil)
		case presentproofsvc.ProposePresentationMsgTypeV2:
			thID, err := action.Message.ThreadID()
			if err != nil {
				logger.Errorf("failed to get thread ID", err)
				action.Stop(nil)
			}

			pd := &presexch.PresentationDefinition{
				ID:   uuid.NewString(),
				Name: "Demo Verifier",
				InputDescriptors: []*presexch.InputDescriptor{
					{ID: uuid.NewString(), Schema: []*presexch.Schema{{URI: "https://w3id.org/citizenship#PermanentResidentCard"}}},
				},
			}

			pdBytes, err := json.Marshal(pd)
			if err != nil {
				logger.Errorf("unable to marshal presentation definition bytes", err)
				action.Stop(nil)
			}

			err = store.Put(thID, pdBytes)
			if err != nil {
				logger.Errorf("failed to save presentation definition", err)
				action.Stop(nil)
			}

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
							PD:        pd},
						},
					},
				},
				WillConfirm: true,
			})

			action.Continue(continueArg)
		case presentproofsvc.PresentationMsgTypeV2:
			thID, err := action.Message.ThreadID()
			if err != nil {
				logger.Errorf("failed to get thread ID", err)
				action.Stop(nil)
			}

			action.Continue(presentproofsvc.WithProperties(
				map[string]interface{}{
					"~web-redirect": &decorator.WebRedirect{
						Status: "OK",
						URL:    os.Getenv(demoExternalURLEnvKey) + "/waci-share/" + thID,
					},
				},
			))
		default:
			action.Stop(nil)
		}
	}
}
