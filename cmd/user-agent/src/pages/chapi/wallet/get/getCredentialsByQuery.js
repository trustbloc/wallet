/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {WalletGet} from "./getCredentials";
import jp from 'jsonpath';
import {PresentationExchange} from '../common/presentationExchange'
import {WalletManager} from "../register/walletManager";
import {getCredentialType, waitForNotification} from '../common/util'
import {AgentMediator} from '../didcomm/connections'

const responseType = "VerifiablePresentation"
// TODO shouldn't hardcode expected credential type (https://github.com/hyperledger/aries-framework-go/issues/2003)
const consentCredentialType = 'ConsentCredential'

/**
 * WalletGetByQuery provides CHAPI get vp features
 * @param aries instance & credential event
 * @class
 */
export class WalletGetByQuery extends WalletGet {
    constructor(aries, credEvent) {
        super(aries, credEvent);

        // validate query and init Presentation Exchange
        let query = jp.query(credEvent, '$..credentialRequestOptions.web.VerifiablePresentation.query[?(@.type=="PresentationDefinitionQuery")]');
        if (query.length == 0) {
            throw "invalid request, incorrect query type"
        }

        this.exchange = new PresentationExchange(query[0].presentationDefinitionQuery)

        this.invitation = jp.query(credEvent, '$..credentialRequestOptions.web.VerifiablePresentation.query[?(@.type=="DIDConnect")].invitation');

        this.walletManager = new WalletManager()

        this.mediator = new AgentMediator(aries)
    }

    requirementDetails() {
        return this.exchange.requirementDetails()
    }

    async connect() {
        // make sure mediator is connected
        await this.mediator.reconnect()
    }

    async getPresentationSubmission() {
        let credentials = await super.getCredentialRecords()

        let vcs = []
        for (let credential of credentials) {
            const resp = await this.aries.verifiable.getCredential({
                id: credential.key
            })

            vcs.push(JSON.parse(resp.verifiableCredential))
        }

        let manifests = await this.walletManager.getAllManifests()
        let submission = this.exchange.createPresentationSubmission(vcs, manifests)

        return submission
    }

    async createAndSendPresentation(walletUser, presentationSubmission, selectedIndexes) {
        try {
            // remove unselected VCs from final presentation submission and get consent credentials for matched manifests.
            if (selectedIndexes && selectedIndexes.length > 0) {
                presentationSubmission = retainOnlySelected(presentationSubmission, selectedIndexes)

                if (this.invitation.length > 0) {
                    presentationSubmission = await getConsentCredentials(this.aries, presentationSubmission, this.walletManager)
                }
            }


            let data
            await this.aries.verifiable.generatePresentation({
                presentation: presentationSubmission,
                did: walletUser.did,
                domain: this.domain,
                challenge: this.challenge,
                skipVerify: true,
                signatureType: walletUser.signatureType
            }).then(resp => {
                    data = resp.verifiablePresentation
                }
            )

            this.sendResponse(responseType, data)
        } catch (e) {
            this.sendResponse("error", e)
        }

    }

    cancel() {
        this.sendResponse("error", "wallet declined credential share")
    }

    sendNoCredntials() {
        this.sendResponse("error", "no credentials found for given presentation exchange request")
    }
}

// retainOnlySelected retain only selected VCs and their respective descriptors
function retainOnlySelected(presentationSubmission, selectedIndexes) {
    let descriptors = []
    let vcs = []

    let vcCount = 0
    selectedIndexes.forEach(function (selected, index) {
        presentationSubmission.verifiableCredential.forEach(function (vc, vcIndex) {
            if (selected && index == vcIndex) {
                vcs.push(vc)

                let vcDescrs = jp.query(presentationSubmission, `$.presentation_submission.descriptor_map[?(@.path=="$.verifiableCredential.[${vcIndex}]")].id`)
                vcDescrs.forEach(function (id) {
                    descriptors.push({id, path: `$.verifiableCredential.[${vcCount}]`})
                })

                vcCount++
            }
        })
    })

    presentationSubmission.verifiableCredential = vcs
    presentationSubmission.presentation_submission.descriptor_map = descriptors

    return presentationSubmission
}

async function getConsentCredentials(aries, presentationSubmission, walletManager) {
    let vcs = []
    // TODO DID exchange with verifier not needed for now
    // let exchange = new DIDExchange(aries)
    // let connection = await exchange.connect(invitation)
    // console.log(`got connection with ${connection.result.State} status`, connection)

    //TODO to improve performance all request should be sent to all issuers in paralle and need to establish
    // cooreltation between incoming actions
    for (let vc of presentationSubmission.verifiableCredential) {
        if (getCredentialType(vc.type) != 'IssuerManifestCredential') {
            vcs.push(vc)
            return
        }

        let connection = await walletManager.getConnectionByID(vc.connection)

        //TODO pass user DID & RP DID as ~attachments
        aries.issuecredential.sendRequest({
            my_did: connection.MyDID,
            their_did: connection.TheirDID,
            request_credential: {},
        })

        let event = await waitForNotification(aries, ["issue-credential_states"], "post_state")
        if (event.StateID != 'request-sent') {
            throw 'failed to send credential request to issuer'
        }
        console.log("sent credential request state notification", event.Message['@id'])

        let piid = event.Message['@id']
        let action = await getAction(aries, piid)
        aries.issuecredential.acceptCredential({
            piid: action.PIID,
            names: [consentCredentialType],
        })

        // TODO shouldn't be 2 call to get credentials (https://github.com/hyperledger/aries-framework-go/issues/2003)
        let metadata = await getCredentialMetadata(aries, consentCredentialType)
        let credential = await aries.verifiable.getCredential(metadata)
        vcs.push(JSON.parse(credential))
    }

    presentationSubmission.verifiableCredential = vcs
    return presentationSubmission
}

// TODO all retry logics below to be replaced by Promise based retry
const retries = 10;

async function getAction(agent, piid) {
    for (let i = 0; i < retries; i++) {
        let resp = await agent.issuecredential.actions()
        if (resp.actions.length > 0) {
            for (let action of resp.actions) {
                if (action.PIID == piid) {
                    return action
                }
            }
        }

        await sleep(1000);
    }

    throw new Error("no action found")
}

async function getCredentialMetadata(agent, name) {
    for (let i = 0; i < retries; i++) {
        try {
            return await agent.verifiable.getCredentialByName({
                name: name
            })
        } catch (e) {
            console.log(`credential '${name}' not found, retrying`)
        }

        await sleep(1000);
    }

    throw new Error("no credential found")
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
