/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// import {getCredentialType} from '../common/util.js';
import {getDomainAndChallenge} from "../common/util";
import {WalletManager} from "../register/walletManager";

const jsonld = require('jsonld');

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
                return await this.agent.verifiable.deriveCredential({credential, frame, skipVerify: true,})
            }

            const allProofs = await Promise.all(selections.map(_getProof));
            let walletMetadata = this.walletManager.getWalletMetadata(user)

            let vcs = allProofs.reduce((acc, val) => acc.concat(val.verifiableCredential), [])
            const {did, signatureType} = await walletMetadata

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
        // filter cred records by query by example
        let records = filterCred(vcs, {types: example.type, contexts: example['@context']});

        // fetch VCs
        let creds = await _fetchCredentials(agent, records)

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
function filterCred(vcs, {types = [], contexts = []}) {
    return vcs.filter(vc => contexts.every(ctx => vc.context.includes(ctx)))
        .filter(vc => types.every(type => vc.type.includes(type)))
}

async function _fetchCredentials(agent, records) {
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
