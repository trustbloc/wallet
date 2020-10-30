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
import {Messenger} from '../common/messaging'
import {AgentMediator} from '../didcomm/mediator'
import {BlindedRouter} from '../didcomm/blindedRouter'
import {DIDManager} from '../didmgmt/didManager'

var uuid = require('uuid/v4')

const responseType = 'VerifiablePresentation'
const manifestType = 'IssuerManifestCredential'
const didDocReqMsgType = 'https://trustbloc.dev/adapter/1.0/diddoc-req'
const didDocResTopic = 'diddoc-res'

/**
 * WalletGetByQuery provides CHAPI get vp features
 * @param agent instance & credential event
 * @class
 */
export class WalletGetByQuery extends WalletGet {
    constructor(agent, credEvent, opts) {
        super(agent, credEvent);

        // validate query and init Presentation Exchange
        let query = jp.query(credEvent, '$..credentialRequestOptions.web.VerifiablePresentation.query[?(@.type=="PresentationDefinitionQuery")]');
        if (query.length == 0) {
            throw "invalid request, incorrect query type"
        }

        this.exchange = new PresentationExchange(query[0].presentationDefinitionQuery)

        this.invitation = jp.query(credEvent, '$..credentialRequestOptions.web.VerifiablePresentation.query[?(@.type=="DIDConnect")].invitation');

        /*
          TODO:
           * current assumption - expecting only one governance VC in request, may be support for multiple.
           * correlate governance VC with requesting party so that consent for trust gets shown only once.
           * verify governance VC proof.
           * verify requesting party in governance framework to make sure this party of behaving properly.
           * request party to get challenged to produce a VP that the governance credential agency has accredited them.
         */
        this.govnVC = jp.query(credEvent, '$..credentialRequestOptions.web.VerifiablePresentation.query[?(@.type=="DIDConnect")].credentials[?(@.type[0]=="GovernanceCredential" || @.type[1]=="GovernanceCredential")]');

        this.walletManager = new WalletManager()
        this.mediator = new AgentMediator(agent)
        this.blindedRouter = new BlindedRouter(agent, opts)
        this.didManager = new DIDManager(agent, opts)
        this.didExchange = new DIDExchange(agent)
        this.messenger = new Messenger(agent)

        // TODO below line to be remove after #434
        this.betaFeature = opts.betaFeature
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
            const resp = await this.agent.verifiable.getCredential({
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
            // remove unselected VCs from final presentation submission and get authorization credentials for matched manifests.
            if (selectedIndexes && selectedIndexes.length > 0) {
                presentationSubmission = retainOnlySelected(presentationSubmission, selectedIndexes)

                if (this.invitation.length > 0) {
                    presentationSubmission = await this._getAuthorizationCredentials(presentationSubmission)
                }
            }

            let walletMetadata = await this.walletManager.getWalletMetadata(walletUser)

            let data
            await this.agent.verifiable.generatePresentation({
                presentation: presentationSubmission,
                did: walletMetadata.did,
                domain: this.domain,
                challenge: this.challenge,
                skipVerify: true,
                signatureType: walletMetadata.signatureType
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


    async _getAuthorizationCredentials(presentationSubmission) {
        let rpConn = await this.didExchange.connect(this.invitation[0])

        if (this.betaFeature) {
            // share peer DID with RP for blinded routing
            await this.blindedRouter.sharePeerDID(rpConn.result)
        }

        // request new peer DID from RP
        let didDocRes = await this.messenger.send(rpConn.result.ConnectionID, {
            "@id": uuid(),
            "@type": didDocReqMsgType,
            "sent_time": new Date().toJSON(),
        }, {
            replyTopic: didDocResTopic,
            timeout: 5000,
            retry: {
                attempts: 5,
                err: 'time out waiting reply for topic'
            }
        })

        // response could be byte array from RP
        let rpDIDDoc =  Array.isArray(didDocRes.data.didDoc) ?
            JSON.parse(String.fromCharCode.apply(String, didDocRes.data.didDoc)) : didDocRes.data.didDoc

        let peerDID = (await this.didManager.createPeerDID()).DID
        let agent = this.agent
        let walletManager = this.walletManager
        let acceptCredPool = new Map()

        await Promise.all(presentationSubmission.verifiableCredential.map(async (vc, index) => {
            if (getCredentialType(vc.type) != manifestType) {
                return
            }

            let connection = await walletManager.getConnectionByID(vc.connection)
            // TODO `request_credential.requests~attach.data.json.subjectDID` to be removed once adapters are updated
            let resp = await agent.issuecredential.sendRequest({
                my_did: connection.MyDID,
                their_did: connection.TheirDID,
                request_credential: {
                    "requests~attach": [
                        {
                            "lastmod_time": new Date(),
                            data: {
                                json: {
                                    subjectDID: peerDID.id,
                                    subjectDIDDoc: {
                                        id: peerDID.id,
                                        doc: peerDID
                                    },
                                    requestingPartyDIDDoc: {
                                        id: rpDIDDoc.id,
                                        doc: rpDIDDoc
                                    },
                                }
                            }
                        }
                    ]
                }
            })

            console.log('sent request credential message', resp.piid)

            acceptCredPool.set(resp.piid, {index})
        }))

        console.log(`${acceptCredPool.size} issue credential requests sent`)

        await waitForCredentials(agent, acceptCredPool)

        acceptCredPool.forEach(function (value) {
            presentationSubmission.verifiableCredential[value.index] = value.credential
        })

        return presentationSubmission
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

                let vcDescrs = jp.query(presentationSubmission, `$.presentation_submission.descriptor_map[?(@.path=="$.verifiableCredential[${vcIndex}]")].id`)
                vcDescrs.forEach(function (id) {
                    descriptors.push({id, path: `$.verifiableCredential[${vcCount}]`})
                })

                vcCount++
            }
        })
    })

    presentationSubmission.verifiableCredential = vcs
    presentationSubmission.presentation_submission.descriptor_map = descriptors

    return presentationSubmission
}

async function waitForCredentials(agent, pool) {
    let processed = 0
    return new Promise(async (resolve, reject) => {
        setTimeout(() => reject(new Error("timout waiting for incoming credentials")), 20000)

        for (; ;) {
            let resp = await agent.issuecredential.actions()
            for (let action of resp.actions) {
                if (pool.has(action.PIID)) {
                    let credID = uuid()
                    agent.issuecredential.acceptCredential({
                        piid: action.PIID,
                        names: [credID],
                    })

                    let metadata = await getCredentialMetadata(agent, credID)
                    let credential = await agent.verifiable.getCredential(metadata)

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
