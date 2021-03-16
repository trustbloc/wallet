/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {getDomainAndChallenge} from "../common/util";
import {WalletManager} from "../register/walletManager";
import {fetchCredentials, filterCred} from "./getCredentialsByFrame";
import {DIDManager} from "..";

const jsonld = require('jsonld');
var uuid = require('uuid/v4')
var flatten = require('flat')

const QUERY_TYPES = ["QueryByFrame", "QueryByExample"]
const CHAPI_RESPONSE_TYPE = 'VerifiablePresentation'

/**
 * MultipleQuery provides support for multiple CHAPI get credential queries.
 * Supported queries - QueryByExample & QueryByFrame
 *
 * @param agent instance & credential event
 * @class
 */
export class MultipleQuery {
    constructor(agent, credEvent) {
        this.agent = agent
        this.credEvent = credEvent

        let {query} = credEvent.credentialRequestOptions.web.VerifiablePresentation
        this.query = Array.isArray(query) ? query : [query]

        if (!this.query.every(q => QUERY_TYPES.includes(q.type))) {
            console.error(`Invalid query types found in request, supported query types ${JSON.stringify(QUERY_TYPES)}`)
            throw 'invalid query'
        }

        const {domain, challenge} = getDomainAndChallenge(credEvent);
        this.domain = domain
        this.challenge = challenge

        this.walletManager = new WalletManager(agent)
        this.didManager = new DIDManager(agent)
    }

    async queryCredentials() {

        let {result} = await this.agent.verifiable.getCredentials()
        if (!result) {
            return []
        }


        return await _mixedQuery(this.agent, result, this.query)
    }

    async generatePresentation(user, selections) {
        try {
            const _getProof = async ({frame, credential}) => {
                if (frame) {
                    let response = await this.agent.verifiable.deriveCredential({
                        credential,
                        frame,
                        skipVerify: true,
                        nonce: uuid()
                    })
                    return JSON.parse(response.verifiableCredential)
                }

                return credential
            }

            const allProofs = await Promise.all(selections.map(_getProof));
            // temp fix, should find metadata from user preference
            let walletMetadata = this._getDIDForSigning(user, selections[0].credential)

            let vcs = allProofs.reduce((acc, val) => acc.concat(val), [])
            const {did, signatureType} = await walletMetadata
            console.log(`presenting with ${did}`)

            let resp = await this.agent.verifiable.generatePresentation({
                verifiableCredential: vcs,
                domain: this.domain,
                challenge: this.challenge,
                skipVerify: true,
                did, signatureType
            })

            this.sendResponse(CHAPI_RESPONSE_TYPE, resp.verifiablePresentation)
        } catch (e) {
            console.error('sending response error', e)
            this.sendResponse("error", e)
        }
    }

    // TODO temp fix, should always read DID from user preference settings
    async _getDIDForSigning(user, credential) {
        try {
            if (credential.credentialSubject.id) {
                let metadata = await this.didManager.getDIDMetadata(credential.credentialSubject.id)
                let {id, signatureType} = metadata
                return {did: id, signatureType}
            }
        } catch (e) {
            console.error('failed to get did from credential subject ID, switching to default DID', e)
        }

        let walletMetadata = await this.walletManager.getWalletMetadata(user)
        let {signatureType, did} = walletMetadata
        return {did, signatureType}
    }

    cancel() {
        this.sendResponse("error", "operation cancelled by user")
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

async function _mixedQuery(agent, vcs, query) {
    const _query = async ({type: queryType, credentialQuery}) => {
        const {example, frame, reason} = credentialQuery
        // filter cred records by example 'type & context' if provided, or else filter by frame 'type & context'
        let records = filterCred(vcs, example ? {
            types: example.type,
            contexts: example['@context']
        } : {types: frame.type, contexts: frame['@context']});

        // fetch VCs
        let creds = await fetchCredentials(agent, records)

        // query by example
        const _getByExample = async (credential) => {
            return {credential, reason}
        };

        // query by frame.
        const _getByFrame = async (credential) => {
            // if no frame, then show all credential details
            let framed = frame ? await jsonld.frame(credential, frame) : credential
            let output = flatCredentialSubject(framed.credentialSubject);
            return {credential, reason, frame, output}
        };

        return queryType == 'QueryByExample' ? await Promise.all(creds.map(_getByExample)) : await Promise.all(creds.map(_getByFrame));
    }

    const requiredResults = await Promise.all(query.map(_query));

    // flatten results
    return requiredResults
        .reduce((acc, val) => acc.concat(val), []);
}

function flatCredentialSubject(subj) {
    return flatten(subj, {
        transformKey: function (key) {
            let parts = key.split('#')
            return parts[parts.length - 1]
        }
    })
}
