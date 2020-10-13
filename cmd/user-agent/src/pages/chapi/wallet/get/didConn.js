/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {WalletManager} from '../register/walletManager'
import {WalletStore} from '../store/saveCredential'
import {DIDExchange} from '../common/didExchange'
import {BlindedRouter} from '../didcomm/blindedRouter'
import {getCredentialType, filterCredentialsByType} from "..";

const manifestCredType = "IssuerManifestCredential"
const governanceCredType = "GovernanceCredential"

var uuid = require('uuid/v4')

/**
 * DIDConn provides CHAPI did connection/exchange features
 * @param aries instance & credential event
 * @class
 */
export class DIDConn {
    constructor(aries, trustblocAgent, trustblocStartupOpts, credEvent, walletUser) {
        this.aries = aries
        this.walletUser = walletUser
        this.walletManager = new WalletManager()
        this.walletStore = new WalletStore(aries, trustblocAgent, trustblocStartupOpts, credEvent, walletUser)
        this.exchange = new DIDExchange(aries)
        this.blindedRouter = new BlindedRouter(aries, trustblocStartupOpts)
        this.credEvent = credEvent

        const {domain, challenge, invitation, credentials} = getRequestParams(credEvent);
        this.domain = domain
        this.challenge = challenge
        this.invitation = invitation
        this.credentials = credentials
    }

    getUserCredentials() {
        return this.credentials ? filterCredentialsByType(this.credentials, [manifestCredType, governanceCredType]) : []
    }

    getGovernanceCredential() {
        let govnVCs = this.credentials ? filterCredentialsByType(this.credentials, [governanceCredType], true) : []

        /*
           TODO:
            * current assumption - expecting only one governance VC in request, may be support for multiple.
            * correlate governance VC with requesting party so that consent for trust gets shown only once.
            * verify governance VC proof.
            * verify requesting party in governance framework to make sure this party of behaving properly.
            * request party to get challenged to produce a VP that the governance credential agency has accredited them.
          */
        return govnVCs.length > 0 ? govnVCs[0] : undefined
    }

    async connect() {
        // perform did exchange
        let connection = await this.exchange.connect(this.invitation)

        // share peer DID with inviter for blinded routing
        await this.blindedRouter.sharePeerDID(connection.result)

        // save wallet metadata
        let walletMetadata = await this.walletManager.getWalletMetadata(this.walletUser)
        if (!walletMetadata.connections) {
            walletMetadata.connections = []
        }
        walletMetadata.connections.push(connection.result.ConnectionID)
        await this.walletManager.storeWalletMetadata(walletMetadata, walletMetadata)

        // save connection
        await this.walletManager.storeConnection(connection.result.ConnectionID, connection.result)

        // save credentials
        if (this.credentials) {
            for (let credential of this.credentials) {
                if (getCredentialType(credential.type) == manifestCredType) {
                    await this.walletManager.storeManifest(connection.result.ConnectionID, credential)
                } else {
                    await this.walletStore.save(uuid(), credential)
                }
            }
        }

        let responseData = await this._didConnResponse(walletMetadata, connection.result)
        this.sendResponse("VerifiablePresentation", responseData)
    }


    async _didConnResponse(walletMetadata, connection) {

        let credential = {
            "@context": [
                "https://www.w3.org/2018/credentials/v1",
                "https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld"
            ],
            issuer: walletMetadata.did,
            issuanceDate: new Date(),
            type: ["VerifiableCredential", "DIDConnection"],
            credentialSubject: {
                id: connection.ConnectionID,
                threadID: connection.ThreadID,
                inviteeDID: connection.MyDID,
                inviterDID: connection.TheirDID,
                inviterLabel: connection.TheirLabel,
                connectionState: connection.State,
            }
        }

        // create did connection VC
        let vc, failure
        await this.aries.verifiable.signCredential({
            credential: credential,
            did: walletMetadata.did,
            signatureType: walletMetadata.signatureType
        }).then(resp => {
                if (!resp.verifiableCredential) {
                    failure = "failed to create did connection credential"
                    return
                }

                vc = resp.verifiableCredential
            }
        ).catch(err => {
            failure = err
        })

        if (failure) {
            console.error("failed to create didconnection credential", failure)
            return failure
        }

        // create did connection response VP
        let presentation = {
            "@context": [
                "https://www.w3.org/2018/credentials/v1"
            ],
            type: "VerifiablePresentation",
            holder: walletMetadata.did,
            verifiableCredential: [vc]
        }

        let data
        await this.aries.verifiable.generatePresentation({
            presentation: presentation,
            domain: this.domain,
            challenge: this.challenge,
            did: walletMetadata.did,
            signatureType: walletMetadata.signatureType,
            skipVerify: true,
        }).then(resp => {
                if (!resp.verifiablePresentation) {
                    data = "failed to create did connection presentation"
                    return
                }

                data = resp.verifiablePresentation
            }
        ).catch(err => {
            data = err
            console.log('failed to create presentation, errMsg:', err)
        })

        return data
    }

    cancel() {
        this.sendResponse("Response", "permission denied")
    }

    sendResponse(type, data) {
        this.credEvent.respondWith(new Promise(function (resolve) {
            return resolve({
                dataType: type,
                data: data
            });
        }))
    }
}

function getRequestParams(credEvent) {
    if (!credEvent.credentialRequestOptions.web.VerifiablePresentation) {
        throw "invitation not found in did connect credential event"
    }

    const verifiable = credEvent.credentialRequestOptions.web.VerifiablePresentation

    let {challenge, domain, query, invitation, credentials} = verifiable;

    if (query && query.challenge) {
        challenge = query.challenge;
    }

    if (query && query.domain) {
        domain = query.domain;
    }

    if (!domain && credEvent.credentialRequestOrigin) {
        domain = credEvent.credentialRequestOrigin.split('//').pop()
    }

    return {domain, challenge, invitation, credentials}
}
