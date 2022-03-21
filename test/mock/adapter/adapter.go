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
	"github.com/square/go-jose/jwt"

	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofband"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofbandv2"
	"github.com/hyperledger/aries-framework-go/pkg/client/presentproof"
	"github.com/hyperledger/aries-framework-go/pkg/common/log"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/decorator"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/issuecredential"
	presentproofsvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/presentproof"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport"
	"github.com/hyperledger/aries-framework-go/pkg/doc/cm"
	"github.com/hyperledger/aries-framework-go/pkg/doc/presexch"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	arieslog "github.com/hyperledger/aries-framework-go/spi/log"
	"github.com/hyperledger/aries-framework-go/spi/storage"
)

var (
	validCredentialManifest []byte
	pdBytes                 []byte

	//go:embed sampledata/sample_prc_fulfillment_unsigned.json
	prcFulFillmentUnsigned []byte

	//go:embed sampledata/sample_prc_fulfillment.json
	prcFulFillment []byte
)

const (
	// issuer html templates
	issuerHTML     = "./templates/issuer/issuer.html"
	waciIssuerHTML = "./templates/issuer/waci-issuer.html"

	// verifier html templates
	verifierHTML     = "./templates/verifier/verifier.html"
	waciVerifierHTML = "./templates/verifier/waci-verifier.html"
	oidcVerifierHTML = "./templates/verifier/oidc-verifier.html"

	// CHAPI html templates
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

	// issuer routes
	router.HandleFunc("/issuer", app.issuer)
	router.HandleFunc("/issuer/waci", app.waciIssuer)
	router.HandleFunc("/issuer/waci-issuance", app.waciIssuance)
	router.HandleFunc("/issuer/waci-issuance-v2", app.waciIssuanceV2)
	router.HandleFunc("/issuer/waci-issuance/{id}", app.waciIssuanceCallback)

	// verifier routes
	router.HandleFunc("/verifier", app.verifier)
	router.HandleFunc("/verifier/waci", app.waciVerifier)
	router.HandleFunc("/verifier/waci-share", app.waciShare)
	router.HandleFunc("/verifier/waci-share-v2", app.waciShareV2)
	router.HandleFunc("/verifier/waci-share/{id}", app.waciShareCallback)
	router.HandleFunc("/verifier/oidc", app.oidcVerifier)
	router.HandleFunc("/verifier/oidc/share", app.oidcShare)
	router.HandleFunc("/verifier/oidc/share/cb", app.oidcShareCallback)

	// CHAPI flow routes
	router.HandleFunc("/web-wallet", app.webWallet)

	return nil
}

// issuer html template endpoints
func (v *adapterApp) issuer(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, issuerHTML, nil)
}

func (v *adapterApp) waciIssuer(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, waciIssuerHTML, nil)
}

// verifier html template endpoints
func (v *adapterApp) verifier(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, verifierHTML, nil)
}

func (v *adapterApp) waciVerifier(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, waciVerifierHTML, nil)
}

func (v *adapterApp) oidcVerifier(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, oidcVerifierHTML, nil)
}

// chapi html template endpoints
func (v *adapterApp) webWallet(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, webWalletHTML, nil)
}

// data endpoints
func (v *adapterApp) waciShare(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// generate OOB invitation
	inv, err := v.agent.OOBClient.CreateInvitation(nil, outofband.WithGoal("share-vp", "streamlined-vp"))
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to create oob invitation : %s", err))

		return
	}

	v.waciInvitationRedirect(w, r, inv)
}

func (v *adapterApp) waciShareV2(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// generate OOB V2 invitation
	inv, err := v.agent.OOBV2Client.CreateInvitation(outofbandv2.WithAccept(transport.MediaTypeDIDCommV2Profile),
		outofbandv2.WithFrom(v.agent.OrbDIDV2), outofbandv2.WithGoal("share-vp", "streamlined-vp"))
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to create oob invitation : %s", err))

		return
	}

	v.waciInvitationRedirect(w, r, inv)
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

	v.waciInvitationRedirect(w, r, inv)
}

func (v *adapterApp) waciIssuanceV2(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// generate OOB V2 invitation
	inv, err := v.agent.OOBV2Client.CreateInvitation(outofbandv2.WithAccept(transport.MediaTypeDIDCommV2Profile),
		outofbandv2.WithFrom(v.agent.OrbDIDV2), outofbandv2.WithGoal("issue-vc", "streamlined-vc"))
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to create oob invitation : %s", err))

		return
	}

	v.waciInvitationRedirect(w, r, inv)
}

func (v *adapterApp) waciInvitationRedirect(w http.ResponseWriter, r *http.Request, inv interface{}) {
	r.ParseForm()

	invBytes, err := json.Marshal(inv)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to unmarshal invitation : %s", err))

		return
	}

	redirectURL := fmt.Sprintf("%s/waci?oob=%s", r.FormValue("walletURL"),
		base64.URLEncoding.EncodeToString(invBytes))

	validCredentialManifest = []byte(r.FormValue("credManifest"))
	pdBytes = []byte(r.FormValue("pEx"))

	logger.Infof("waci redirect : url=%s oob-invitation=%s credentialManifest=%s presentationExchange=%s",
		redirectURL, string(invBytes), string(validCredentialManifest), string(pdBytes))

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

	loadTemplate(w, waciVerifierHTML, map[string]interface{}{"Msg": "Successfully Received Presentation"})
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

	loadTemplate(w, waciIssuerHTML, map[string]interface{}{"Msg": "Successfully Sent Credential to holder"})
}

func (v *adapterApp) oidcShare(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	walletAuthURL := r.FormValue("walletAuthURL")
	pdBytes = []byte(r.FormValue("pEx"))

	var pd *presexch.PresentationDefinition
	err := json.Unmarshal(pdBytes, &pd)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to unmarshal presentation definition : %s", err))

		return
	}

	authClaims := &OIDCAuthClaims{
		VPToken: &VPToken{
			PresDef: pd,
		},
	}

	claimsBytes, err := json.Marshal(authClaims)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to unmarshal invitation : %s", err))

		return
	}

	state := uuid.NewString()

	// TODO: use OIDC client library
	// construct wallet auth req with PEx
	req, err := http.NewRequest("GET", walletAuthURL, nil)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to get interaction data : %s", err))

		return
	}

	q := req.URL.Query()
	q.Add("client_id", "demo-verifier")
	q.Add("redirect_uri", os.Getenv(demoExternalURLEnvKey)+"/verifier/oidc/share/cb")
	q.Add("scope", "openid")
	q.Add("state", state)
	q.Add("claims", string(claimsBytes))

	req.URL.RawQuery = q.Encode()

	redirectURL := req.URL.String()

	logger.Infof("oidc share redirect : url=%s claims=%s", redirectURL, string(claimsBytes))

	err = v.store.Put(state, pdBytes)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to save state data : %s", err))

		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (v *adapterApp) oidcShareCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")

	pdBytes, err := v.store.Get(state)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to get oidc state data : %s", err))

		return
	}

	var pd *presexch.PresentationDefinition

	err = json.Unmarshal(pdBytes, &pd)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to unmarshal presentation definition : %s", err))

		return
	}

	idToken := r.URL.Query().Get("id_token")
	vpToken := r.URL.Query().Get("vp_token")

	logger.Infof("oidc share callback : id_token=%s vp_token=%s",
		idToken, vpToken)

	var claims *OIDCTokenCliams

	token, err := jwt.ParseSigned(idToken)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to parsed token : %s", err))

		return
	}

	err = token.UnsafeClaimsWithoutVerification(&claims)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to convert to claim object : %s", err))

		return
	}

	presSubBytes, err := json.Marshal(claims.VPToken)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to marshal _vp_token : %s", err))

		return
	}

	logger.Infof("oidc share callback : _vp_token=%v vp_token=%s", string(presSubBytes), vpToken)

	_, err = verifiable.ParsePresentation([]byte(vpToken), verifiable.WithPresJSONLDDocumentLoader(ld.NewDefaultDocumentLoader(nil)))
	if err != nil {
		loadTemplate(w, oidcVerifierHTML,
			map[string]interface{}{
				"ErrMsg":                    fmt.Sprintf("ERROR: failed to validate presentation : %s", err),
				"ID_TOKEN":                  "ID_TOKEN:\n" + idToken,
				"DECODED_VPDEF_IN_ID_TOKEN": "DECODED_VPDEF_IN_ID_TOKEN: " + string(presSubBytes),
				"VP_TOKEN":                  "VP_TOKEN: " + string(vpToken),
			},
		)

		return
	}

	loadTemplate(w, oidcVerifierHTML,
		map[string]interface{}{
			"Msg":                       "Successfully Received Presentation",
			"ID_TOKEN":                  "ID_TOKEN:\n" + idToken,
			"DECODED_VPDEF_IN_ID_TOKEN": "DECODED_VPDEF_IN_ID_TOKEN: " + string(presSubBytes),
			"VP_TOKEN":                  "VP_TOKEN: " + string(vpToken),
		},
	)
}

func listenForDIDCommMsg(actionCh chan service.DIDCommAction, store storage.Store) {
	for action := range actionCh {
		logger.Infof("received action message : type=%s", action.Message.Type())

		switch action.Message.Type() {
		case didexchange.RequestMsgType:
			action.Continue(nil)
		case presentproofsvc.ProposePresentationMsgTypeV2, presentproofsvc.ProposePresentationMsgTypeV3:
			thID, err := action.Message.ThreadID()
			if err != nil {
				logger.Errorf("failed to get thread ID", err)
				action.Stop(nil)
			}

			var pd *presexch.PresentationDefinition

			err = json.Unmarshal(pdBytes, &pd)
			if err != nil {
				logger.Errorf("failed to unmarshal presentation definition", err)
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
							PD:        pd,
						},
						},
					},
				},
				WillConfirm: true,
			})

			action.Continue(continueArg)
		case presentproofsvc.PresentationMsgTypeV2, presentproofsvc.PresentationMsgTypeV3:
			thID, err := action.Message.ThreadID()
			if err != nil {
				logger.Errorf("failed to get thread ID", err)
				action.Stop(nil)
			}

			action.Continue(presentproofsvc.WithProperties(
				map[string]interface{}{
					"~web-redirect": &decorator.WebRedirect{
						Status: "OK",
						URL:    os.Getenv(demoExternalURLEnvKey) + "/verifier/waci-share/" + thID,
					},
				},
			))
		case issuecredential.ProposeCredentialMsgTypeV2, issuecredential.ProposeCredentialMsgTypeV3:
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

			offerCredMsg, err := createOfferCredentialMsg(validCredentialManifest, prcFulFillmentUnsigned)
			if err != nil {
				logger.Errorf("failed to prepare offer credential message", err)
				action.Stop(nil)
			}

			action.Continue(issuecredential.WithOfferCredential(offerCredMsg))
		case issuecredential.RequestCredentialMsgTypeV2, issuecredential.RequestCredentialMsgTypeV3:
			thID, err := action.Message.ThreadID()
			if err != nil {
				logger.Errorf("failed to get thread ID", err)
				action.Stop(nil)
			}

			issueCredMsg, err := createIssueCredentialMsg(prcFulFillment, os.Getenv(demoExternalURLEnvKey)+"/issuer/waci-issuance/"+thID)
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

func createOfferCredentialMsg(manifest, fulfillmentVP []byte) (*issuecredential.OfferCredentialParams, error) {
	var credentialManifest cm.CredentialManifest

	err := json.Unmarshal(manifest, &credentialManifest)
	if err != nil {
		return nil, err
	}

	vp, err := verifiable.ParsePresentation([]byte(fulfillmentVP), verifiable.WithPresJSONLDDocumentLoader(ld.NewDefaultDocumentLoader(nil)))
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
				Format:    format1,
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
				Format:    format2,
				MediaType: "application/json",
				Data: decorator.AttachmentData{
					JSON: vp,
				},
			},
		},
	}, nil
}

func createIssueCredentialMsg(vp []byte, redirect string) (*issuecredential.IssueCredentialParams, error) {
	attachID := uuid.New().String()

	// change credential ID
	vpStr := strings.ReplaceAll(string(vp), "http://example.gov/credentials/3732", "http://example.gov/credentials/"+uuid.NewString())

	presentation, err := verifiable.ParsePresentation([]byte(vpStr), verifiable.WithPresDisabledProofCheck(),
		verifiable.WithPresJSONLDDocumentLoader(ld.NewDefaultDocumentLoader(nil)))
	if err != nil {
		return nil, err
	}

	format := "dif/credential-manifest/fulfillment@v1.0"

	return &issuecredential.IssueCredentialParams{
		Type: issuecredential.IssueCredentialMsgTypeV2,
		Formats: []issuecredential.Format{{
			AttachID: attachID,
			Format:   format,
		}},
		Attachments: []decorator.GenericAttachment{{
			ID:        attachID,
			Format:    format,
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
