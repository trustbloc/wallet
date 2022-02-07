/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/piprate/json-gold/ld"

	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofband"
	"github.com/hyperledger/aries-framework-go/pkg/client/presentproof"
	"github.com/hyperledger/aries-framework-go/pkg/common/log"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/decorator"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/issuecredential"
	presentproofsvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/presentproof"
	"github.com/hyperledger/aries-framework-go/pkg/doc/cm"
	"github.com/hyperledger/aries-framework-go/pkg/doc/presexch"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	arieslog "github.com/hyperledger/aries-framework-go/spi/log"
	"github.com/hyperledger/aries-framework-go/spi/storage"
)

var (
	//go:embed sampledata/credential_manifest.json
	validCredentialManifest []byte

	//go:embed sampledata/fulfillment_dl_vp.json
	fulfillmentDrivingLicenseVP []byte

	//go:embed sampledata/sample_udc_fulfillment.json
	universityDegreeFulFillment []byte
)

const (
	verifierHTML  = "./templates/verifier.html"
	issuerHTML    = "./templates/issuer.html"
	webWalletHTML = "./templates/webWallet.html"
)

var logger = log.New("mock-adapter")

type adapterApp struct {
	agent *didComm
	store storage.Store
}

func startAdapterApp(agent *didComm, router *mux.Router) error {
	log.SetLevel("", arieslog.DEBUG)

	prov := mem.NewProvider()

	store, err := prov.OpenStore("verifier")
	if err != nil {
		return fmt.Errorf("failed to create store : %w", err)
	}

	app := adapterApp{agent: agent, store: store}

	actionCh := make(chan service.DIDCommAction)

	err = agent.DIDExchClient.RegisterActionEvent(actionCh)
	if err != nil {
		return fmt.Errorf("failed to register action events on didexchange-client : %w", err)
	}

	err = agent.PresentProofClient.RegisterActionEvent(actionCh)
	if err != nil {
		return fmt.Errorf("failed to register action events on present-proof-client : %w", err)
	}

	err = agent.IssueCredentialClient.RegisterActionEvent(actionCh)
	if err != nil {
		return fmt.Errorf("failed to register action events on issue-credential-client : %w", err)
	}

	go listenForDIDCommMsg(actionCh, store)

	router.HandleFunc("/verifier", app.verifier)
	router.HandleFunc("/issuer", app.issuer)
	router.HandleFunc("/web-wallet", app.webWallet)
	router.HandleFunc("/waci-share", app.waciShare)
	router.HandleFunc("/waci-share/{id}", app.waciShareCallback)
	router.HandleFunc("/waci-issuance", app.waciIssuance)
	router.HandleFunc("/waci-issuance/{id}", app.waciIssuanceCallback)

	return nil
}

func (v *adapterApp) verifier(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, verifierHTML, nil)
}

func (v *adapterApp) issuer(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, issuerHTML, nil)
}

func (v *adapterApp) webWallet(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, webWalletHTML, nil)
}

func (v *adapterApp) waciShare(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// generate OOB invitation
	inv, err := v.agent.OOBClient.CreateInvitation(nil, outofband.WithGoal("share-vp", "streamlined-vp"))
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

func (v *adapterApp) waciIssuance(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// generate OOB invitation
	inv, err := v.agent.OOBClient.CreateInvitation(nil, outofband.WithGoal("issue-vc", "streamlined-vc"))
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

func (v *adapterApp) waciShareCallback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := v.store.Get(id)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to get interaction data : %s", err))

		return
	}

	loadTemplate(w, verifierHTML, map[string]interface{}{"Msg": "Successfully Received Presentation"})
}

func (v *adapterApp) waciIssuanceCallback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := v.store.Get(id)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to get interaction data : %s", err))

		return
	}

	loadTemplate(w, issuerHTML, map[string]interface{}{"Msg": "Successfully Sent Credential to holder"})
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
				Attachments: []decorator.GenericAttachment{
					{
						ID:        uuid.NewString(),
						MediaType: "application/json",
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
		case issuecredential.ProposeCredentialMsgTypeV2:
			thID, err := action.Message.ThreadID()
			if err != nil {
				logger.Errorf("failed to get thread ID", err)
				action.Stop(nil)
			}

			err = store.Put(thID, []byte(thID))
			if err != nil {
				logger.Errorf("failed to save interaction data", err)
				action.Stop(nil)
			}

			offerCredMsg, err := createOfferCredentialMsg(validCredentialManifest, fulfillmentDrivingLicenseVP, thID)
			if err != nil {
				logger.Errorf("failed to prepare offer credential message", err)
				action.Stop(nil)
			}

			action.Continue(issuecredential.WithOfferCredential(offerCredMsg))
		case issuecredential.RequestCredentialMsgTypeV2:
			thID, err := action.Message.ThreadID()
			if err != nil {
				logger.Errorf("failed to get thread ID", err)
				action.Stop(nil)
			}

			issueCredMsg, err := createIssueCredentialMsg(universityDegreeFulFillment, os.Getenv(demoExternalURLEnvKey)+"/waci-issuance/"+thID, thID)
			if err != nil {
				logger.Errorf("failed to prepare issue credential message", err)
				action.Stop(nil)
			}

			action.Continue(issuecredential.WithIssueCredential(issueCredMsg))
		default:
			action.Stop(nil)
		}
	}
}

func createOfferCredentialMsg(manifest, fulfillmentVP []byte, thID string) (*issuecredential.OfferCredentialParams, error) {
	var credentialManifest cm.CredentialManifest

	err := json.Unmarshal(manifest, &credentialManifest)
	if err != nil {
		return nil, err
	}

	// change manifest ID and also update fulfillment with new ID (to make test run multiple times)
	credentialManifest.ID = thID
	fulfillmentVPStr := strings.ReplaceAll(string(fulfillmentVP), `"manifest_id": "dcc75a16-19f5-4273-84ce-4da69ee2b7fe"`,
		fmt.Sprintf(`"manifest_id": "%s"`, thID))

	vp, err := verifiable.ParsePresentation([]byte(fulfillmentVPStr), verifiable.WithPresJSONLDDocumentLoader(ld.NewDefaultDocumentLoader(nil)))
	if err != nil {
		return nil, err
	}

	attachID1, attachID2 := uuid.New().String(), uuid.New().String()
	format1, format2 := "dif/credential-manifest/manifest@v1.0", "dif/credential-manifest/fulfillment@v1.0"

	return &issuecredential.OfferCredentialParams{
		Type:    issuecredential.OfferCredentialMsgTypeV2,
		Comment: "Offer to issue University Degree Credential for Mr.Smith",
		Formats: []issuecredential.Format{{
			AttachID: attachID1,
			Format:   format1,
		}, {
			AttachID: attachID2,
			Format:   format2,
		}},
		Attachments: []decorator.GenericAttachment{
			{
				ID:        attachID1,
				MediaType: "application/json",
				Data: decorator.AttachmentData{
					JSON: struct {
						Manifest cm.CredentialManifest `json:"credential_manifest,omitempty"`
					}{
						Manifest: credentialManifest,
					},
				},
			},
			{
				ID:        attachID2,
				MediaType: "application/json",
				Data: decorator.AttachmentData{
					JSON: vp,
				},
			},
		},
	}, nil
}

func createIssueCredentialMsg(vp []byte, redirect, thID string) (*issuecredential.IssueCredentialParams, error) {
	attachID := uuid.New().String()

	// change ID & manifest ID
	vpStr := strings.ReplaceAll(string(vp), "http://example.gov/credentials/3732", "http://example.gov/credentials/"+uuid.NewString())
	vpStr = strings.ReplaceAll(string(vpStr), `"manifest_id": "dcc75a16-19f5-4273-84ce-4da69ee2b7fe"`,
		fmt.Sprintf(`"manifest_id": "%s"`, thID))

	presentation, err := verifiable.ParsePresentation([]byte(vpStr), verifiable.WithPresDisabledProofCheck(),
		verifiable.WithPresJSONLDDocumentLoader(ld.NewDefaultDocumentLoader(nil)))
	if err != nil {
		return nil, err
	}

	return &issuecredential.IssueCredentialParams{
		Type: issuecredential.IssueCredentialMsgTypeV2,
		Formats: []issuecredential.Format{{
			AttachID: attachID,
			Format:   "dif/credential-manifest/fulfillment@v1.0",
		}},
		Attachments: []decorator.GenericAttachment{{
			ID:        attachID,
			MediaType: "application/ld+json",
			Data: decorator.AttachmentData{
				JSON: presentation,
			},
		}},
		WebRedirect: &decorator.WebRedirect{
			Status: "OK",
			URL:    redirect,
		},
	}, nil
}
