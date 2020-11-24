/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

var uuid = require('uuid/v4')

const didDocReqMsgType = 'https://trustbloc.dev/blinded-routing/1.0/diddoc-req'
const didDocReqRespType = 'https://trustbloc.dev/blinded-routing/1.0/diddoc-resp'
const registerRouteReqType = 'https://trustbloc.dev/blinded-routing/1.0/register-route-req'
const registerRouteRespType = 'https://trustbloc.dev/blinded-routing/1.0/register-route-resp'

/**
 * BlindedRouter provides blinded routing features
 * @param agent instance
 * @class
 */
export class BlindedRouter {
    constructor(agent, opts) {
        if (!opts.blindedRouting) {
            return
        }

        this.agent = agent
    }

    async sharePeerDID(connection) {
        if (!this.agent) {
            console.debug('Not sharing router peer DID since blinded routing is disabled !')
            return
        }

        let {ConnectionID} = connection

        console.debug('Sending DID Doc request')
        // request peer DID from other party
        let payload
        const retries = 5;
        // TODO this retry logic to be removed once integration issue with did-exchange state complete update delay is resolved
        for (let i = 1; i <= retries; i++) {
            try {
                payload = await this.agent.messaging.send({
                    "connection_ID": ConnectionID,
                    "message_body": {
                        "@id": uuid(),
                        "@type": didDocReqMsgType,
                        "sent_time": new Date().toJSON(),
                    },
                    "await_reply": {messageType: didDocReqRespType, timeout: 2000000000},
                })
            } catch (e) {
                if (!e.message.includes("failed to get reply") || i === retries) {
                    throw e
                }
                await sleep(3000);
                continue
            }
            break
        }

        if(!payload.response) {
            throw 'no response DID found in did doc response'
        }

        console.debug('requesting peer DID from wallet')

        let {message} = payload.response
        let peerDID = _parseResponseDID(message)
        if (!peerDID) {
            console.error('failed to get peerDID from inviter, could not find peer DID in response message.')
            throw 'failed to get peer DID from inviter'
        }

        // request wallet peer DID from router by sending peer DID from other party
        let walletDID = await requestDIDFromMediator(this.agent, peerDID)

        console.log('sharing wallet peer DID to inviter')
        // share wallet peer DID to other party
        await this.agent.messaging.reply({
            "message_ID": message['@id'],
            "message_body": {
                "@id": uuid(),
                "@type": registerRouteReqType,
                data: {didDoc: walletDID},
                "sent_time": new Date().toJSON(),
            },
            "start_new_thread": true,
            "await_reply": {messageType: registerRouteRespType, timeout: 2000000000}
        })

    }
}

let _parseResponseDID = (response) => response.data.didDoc ? response.data.didDoc : undefined

async function requestDIDFromMediator(agent, reqDoc) {
    let res = await agent.mediatorclient.sendCreateConnectionRequest({
        didDoc: reqDoc
    })

    if (res.payload && res.payload.message) {
        let response = res.payload.message
        // TODO currently getting routerDIDDoc as byte[], to be fixed
        if (response.data.didDoc && response.data.didDoc.length > 0) {
            return JSON.parse(String.fromCharCode.apply(String, response.data.didDoc))
        }
    }

    console.error('failed to request DID from router, failed to get connection response')
    throw 'failed to request DID from router, failed to get connection response'
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

