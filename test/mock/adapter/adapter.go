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
	"net/url"
	"os"
	"strings"
	"time"

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
	issuerHTML          = "./templates/issuer/issuer.html"
	waciIssuerHTML      = "./templates/issuer/waci-issuer.html"
	oidcIssuerHTML      = "./templates/issuer/oidc-issuer.html"
	oidcIssuerLoginHTML = "./templates/issuer/oidc-login.html"

	// verifier html templates
	verifierHTML     = "./templates/verifier/verifier.html"
	waciVerifierHTML = "./templates/verifier/waci-verifier.html"
	oidcVerifierHTML = "./templates/verifier/oidc-verifier.html"

	// CHAPI html templates
	webWalletHTML = "./templates/webWallet.html"
)

var logger = log.New("mock-adapter")

type issuerConfiguration struct {
	Issuer                string          `json:"issuer"`
	AuthorizationEndpoint string          `json:"authorization_endpoint"`
	CredentialManifests   json.RawMessage `json:"credential_manifests"`
}

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
	router.HandleFunc("/issuer/oidc", app.oidcIssuer)
	router.HandleFunc("/issuer/oidc/login", app.oidcIssuerLogin)
	router.HandleFunc("/issuer/oidc/issuance", app.initiateIssuance).Methods(http.MethodPost)
	router.HandleFunc("/{id}/.well-known/openid-configuration", app.wellKnownConfiguration).Methods(http.MethodGet)
	router.HandleFunc("/{id}/issuer/oidc/authorize", app.issuerAuthorize).Methods(http.MethodGet)
	router.HandleFunc("/issuer/oidc/authorize-request", app.issuerSendAuthorizeResponse).Methods(http.MethodPost)

	// verifier routes
	router.HandleFunc("/verifier", app.verifier)
	router.HandleFunc("/verifier/waci", app.waciVerifier)
	router.HandleFunc("/verifier/waci-share", app.waciShare)
	router.HandleFunc("/verifier/waci-share-v2", app.waciShareV2)
	router.HandleFunc("/verifier/waci-share/{id}", app.waciShareCallback)
	router.HandleFunc("/verifier/oidc", app.oidcVerifier)
	router.HandleFunc("/verifier/oidc/share", app.oidcShare)
	router.HandleFunc("/verifier/oidc/share/cb", app.oidcShareCallback)
	router.HandleFunc("/verifier/oidc/share/registration", app.oidcShareRegistration)


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

func (v *adapterApp) oidcIssuer(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, oidcIssuerHTML, nil)
}

func (v *adapterApp) oidcIssuerLogin(w http.ResponseWriter, r *http.Request) {
	loadTemplate(w, oidcIssuerLoginHTML, nil)
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
	redirectURLParam := os.Getenv(demoExternalURLEnvKey) + "/verifier/oidc/share/cb"
	registrationURLParam := os.Getenv(demoExternalURLEnvKey) + "/verifier/oidc/share/registration"
	// To follow the specification while we don't have a signed token the client_id = redirect_uri
	q.Add("client_id", redirectURLParam)
	q.Add("redirect_uri", redirectURLParam)
	q.Add("scope", "openid")
	q.Add("state", state)
	q.Add("claims", string(claimsBytes))
	// Remove the following param if we are using a signed token and pre registration
	q.Add("registration_uri", registrationURLParam)

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

	var claims *OIDCTokenClaims

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

	_, err = verifiable.ParsePresentation([]byte(vpToken), verifiable.WithPresJSONLDDocumentLoader(ld.NewDefaultDocumentLoader(nil)), verifiable.WithPresDisabledProofCheck())
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

func (v *adapterApp) oidcShareRegistration(w http.ResponseWriter, r *http.Request) {
    data := RegistrationMetadata{"client_id_from_rp"}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(data)
}

func (v *adapterApp) initiateIssuance(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	walletURL := r.FormValue("walletInitIssuanceURL")
	credentialTypes := strings.Split(r.FormValue("credentialTypes"), ",")
	manifestIDs := strings.Split(r.FormValue("manifestIDs"), ",")
	issuerURL := r.FormValue("issuerURL")
	credManifest := r.FormValue("credManifest")

	key := uuid.NewString()
	issuer := issuerURL + "/" + key
	issuerConf, err := json.MarshalIndent(&issuerConfiguration{
		Issuer:                issuer,
		AuthorizationEndpoint: issuer + "/issuer/oidc/authorize",
		CredentialManifests:   []byte(credManifest),
	}, "", "	")
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to prepare issuer wellknown configuration : %s", err))

		return
	}

	err = v.store.Put(key, issuerConf)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to server configuration : %s", err))

		return
	}

	u, err := url.Parse(walletURL)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to parse wallet init issuance URL : %s", err))

		return
	}

	q := u.Query()
	q.Set("issuer", issuer)

	for _, credType := range credentialTypes {
		q.Add("credential_type", credType)
	}

	for _, manifestID := range manifestIDs {
		q.Add("manifest_id", manifestID)
	}

	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusFound)
}

func (v *adapterApp) wellKnownConfiguration(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	issuerConf, err := v.store.Get(id)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to read wellknown configuration : %s", err))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(issuerConf)
}
func (v *adapterApp) issuerAuthorize(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(w, http.StatusBadRequest,
			fmt.Sprintf("failed to parse request : %s", err))

		return
	}

	claims, err := url.PathUnescape(r.Form.Get("claims"))
	if err != nil {
		handleError(w, http.StatusBadRequest,
			fmt.Sprintf("failed to read claims : %s", err))

		return
	}

	redirectURI, err := url.PathUnescape(r.Form.Get("redirect_uri"))
	if err != nil {
		handleError(w, http.StatusBadRequest,
			fmt.Sprintf("failed to read redirect URI : %s", err))

		return
	}

	scope := r.Form.Get("scope")
	state := r.Form.Get("state")
	responseType := r.Form.Get("response_type")
	clientID := r.Form.Get("client_id")

	// basic validation only.
	if claims == "" || redirectURI == "" || clientID == "" {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Invalid Request"))

		return
	}

	authState := uuid.NewString()
	authRequest, err := json.Marshal(map[string]string{
		"claims":        claims,
		"scope":         scope,
		"state":         state,
		"response_type": responseType,
		"client_id":     clientID,
		"redirect_uri":  redirectURI,
	})
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to process authorization request : %s", err))

		return
	}

	err = v.store.Put(getAuthStateKeyPrefix(authState), authRequest)
	if err != nil {
		handleError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to save state : %s", err))

		return
	}

	authStateCookie := http.Cookie{
		Name:    "state",
		Value:   authState,
		Expires: time.Now().Add(5 * time.Minute),
		Path:    "/",
	}

	http.SetCookie(w, &authStateCookie)
	http.Redirect(w, r, "/issuer/oidc/login", http.StatusFound)
}

func (v *adapterApp) issuerSendAuthorizeResponse(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("state")
	if err != nil {
		handleError(w, http.StatusForbidden, "invalid state")

		return
	}

	authRqstBytes, err := v.store.Get(getAuthStateKeyPrefix(stateCookie.Value))
	if err != nil {
		handleError(w, http.StatusBadRequest, "invalid request")

		return
	}

	var authRequest map[string]string
	err = json.Unmarshal(authRqstBytes, &authRequest)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "failed to read request")

		return
	}

	redirectURI, ok := authRequest["redirect_uri"]
	if !ok {
		handleError(w, http.StatusInternalServerError, "failed to redirect, invalid URL")

		return
	}

	// TODO process credential types or manifests from claims and prepare credential endpoint with credential to be issued.
	// TODO Provide authorize response along with redirect

	http.Redirect(w, r, redirectURI, http.StatusFound)
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

func getAuthStateKeyPrefix(key string) string {
	return fmt.Sprintf("authstate_%s", key)
}
