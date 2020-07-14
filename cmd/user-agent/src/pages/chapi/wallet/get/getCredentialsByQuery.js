/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {WalletGet} from "./getCredentials";
import jp from 'jsonpath';
import {PresentationExchange} from '../common/presentationExchange'
import {WalletManager} from "../register/walletManager";
import {getCredentialType} from '../common/util'
import {DIDExchange} from '../common/didExchange'
import {AgentMediator} from '../didcomm/connections'

var uuid = require('uuid/v4')

const responseType = 'VerifiablePresentation'
const manifestType = 'IssuerManifestCredential'

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
                    presentationSubmission = await getConsentCredentials(this.aries, presentationSubmission, this.invitation[0], this.walletManager)
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
            console.error('sending response error', e)
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

async function getConsentCredentials(aries, presentationSubmission, invitation, walletManager) {
    let exchange = new DIDExchange(aries)
    let rpConn = await exchange.connect(invitation)
    let rpDIDDoc = await aries.vdri.resolveDID({id: rpConn.result.TheirDID})

    let acceptCredPool = new Map()
    await Promise.all(presentationSubmission.verifiableCredential.map(async (vc, index) => {
        if (getCredentialType(vc.type) != manifestType) {
            return
        }

        let connection = await walletManager.getConnectionByID(vc.connection)
        let resp = await aries.issuecredential.sendRequest({
            my_did: connection.MyDID,
            their_did: connection.TheirDID,
            request_credential: {
                userDID: rpConn.result.MyDID,
                rpDIDDoc: rpDIDDoc,
            },
        })

        console.log('sent request credential message', resp.piid)

        acceptCredPool.set(resp.piid, {index})
    }))

    console.log(`${acceptCredPool.size} issue credential requests sent`)

    await waitForCredentials(aries, acceptCredPool)

    acceptCredPool.forEach(function (value) {
        presentationSubmission.verifiableCredential[value.index] = value.credential
    })

    // presentationSubmission.verifiableCredential = vcs
    return presentationSubmission
}

async function waitForCredentials(aries, pool) {
    let processed = 0
    return new Promise(async (resolve, reject) => {
        setTimeout(() => reject(new Error("timout waiting for incoming credentials")), 20000)

        for (; ;) {
            let resp = await aries.issuecredential.actions()
            for (let action of resp.actions) {
                if (pool.has(action.PIID)) {
                    let credID = uuid()
                    aries.issuecredential.acceptCredential({
                        piid: action.PIID,
                        names: [credID],
                    })

                    let metadata = await getCredentialMetadata(aries, credID)
                    let credential = await aries.verifiable.getCredential(metadata)

                    pool.get(action.PIID).credential = JSON.parse(credential.verifiableCredential)

                    processed++
                }
            }

            if (processed < pool.size) {
                console.log(`received ${processed} out of ${pool.size} credentials, retrying`)
                await sleep(1000);
                continue
            }

            break
        }

        console.log(`received all ${processed} credentials`)
        resolve()
    })
}


async function getCredentialMetadata(agent, name) {
    const retries = 10;
    for (let i = 0; i < retries; i++) {
        try {
            return await agent.verifiable.getCredentialByName({
                name: name
            })
        } catch (e) {
            console.log(`credential '${name}' not found, retrying`)
        }

        await sleep(500);
    }

    throw new Error("no credential found")
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
