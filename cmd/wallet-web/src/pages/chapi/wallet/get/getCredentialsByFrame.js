/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// import {getCredentialType} from '../common/util.js';
import {getDomainAndChallenge} from "../common/util";
import {WalletManager} from "../register/walletManager";
import {DIDManager} from "..";

const jsonld = require('jsonld');
var uuid = require('uuid/v4')

const QUERY_TYPE = "QueryByFrame"
const CHAPI_RESPONSE_TYPE = 'VerifiablePresentation'

/**
 * SelectiveDisclosure provides CHAPI get credential 'QueryByFrame' features
 * @param agent instance & credential event
 * @class
 */
export class SelectiveDisclosure {
    constructor(agent, credEvent) {
        this.agent = agent
        this.credEvent = credEvent

        const {query} = credEvent.credentialRequestOptions.web.VerifiablePresentation
        const {type, credentialQuery} = Array.isArray(query) ? query[0] : query

        if (type != QUERY_TYPE) {
            console.error(`Invalid query type '${type}', selective disclosure is supported only for ${QUERY_TYPE}`)
            throw 'invalid query type'
        }

        const {domain, challenge} = getDomainAndChallenge(credEvent);

        this.credentialQuery = credentialQuery
        this.domain = domain
        this.challenge = challenge

        this.walletManager = new WalletManager(agent)
        this.didManager = new DIDManager(agent)
    }

    async queryByFrame() {
        let records = []
        await this.agent.verifiable.getCredentials().then(
            resp => {
                if (resp.result) {
                    resp.result.forEach((item) => {
                        records.push(item)
                    })
                }
            })
            .catch(err => {
                throw err
            })

        return await _queryByFrame(this.agent, records, this.credentialQuery)
    }

    async generatePresentation(user, selections) {
        try {
            const _getProof = async ({frame, credential}) => {
                return await this.agent.verifiable.deriveCredential({
                    credential,
                    frame,
                    skipVerify: true,
                    nonce: uuid()
                })
            }

            const allProofs = await Promise.all(selections.map(_getProof));
            // temp fix, should find metadata from user preference
            let walletMetadata = this._getDIDForSigning(user, selections[0].credential)

            let vcs = allProofs.reduce((acc, val) => acc.concat(JSON.parse(val.verifiableCredential)), [])
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

async function _queryByFrame(agent, vcs, credentialQuery) {
    let query = Array.isArray(credentialQuery) ? credentialQuery : [credentialQuery]

    const _query = async ({example, frame}) => {
        // filter cred records by example 'type & context' if provided, or else filter by frame 'type & context'
        let records = filterCred(vcs, example ? {
            types: example.type,
            contexts: example['@context']
        } : {types: frame.type, contexts: frame['@context']});

        // fetch VCs
        let creds = await fetchCredentials(agent, records)

        // get output frame from each credentials.
        const _getFrame = async (credential) => {
            // if no frame, then show all credential details
            let output = frame ? await jsonld.frame(credential, frame) : credential
            return {credential, frame, output}
        };

        return await Promise.all(creds.map(_getFrame));
    }

    const requiredResults = await Promise.all(query.map(_query));

    // flatten results
    return requiredResults
        .reduce((acc, val) => acc.concat(val), []);
}

// TODO currently supporting search by contexts & types only.
export function filterCred(vcs, {types = [], contexts = []}) {
    return vcs.filter(vc => contexts.every(ctx => vc.context.includes(ctx)))
        .filter(vc => types.every(type => vc.type.includes(type)))
}

export async function fetchCredentials(agent, records) {
    const _fetchvc = record => {
        return agent.verifiable.getCredential({
            id: record.id
        })
    }

    // fetch all credentials concurrently
    const requiredCredentials = await Promise.all(records.map(_fetchvc));
    // flatten the results
    return requiredCredentials.reduce((acc, val) => acc.concat(JSON.parse(val.verifiableCredential)), []);
}
