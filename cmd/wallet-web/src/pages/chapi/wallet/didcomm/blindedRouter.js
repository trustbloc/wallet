/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

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

        // request peer DID from other party
        console.debug('Sending DID Doc request')
        let response = await this.agent.blindedrouting.sendDIDDocRequest({connectionID: ConnectionID})
        console.log("payload from did doc request ", JSON.stringify(response, null ,2))
        if (!response.payload) {
            throw 'no response DID found in did doc response'
        }

        let {message} = response.payload
        let peerDID = _parseResponseDID(message)
        if (!peerDID) {
            console.error('failed to get peerDID from inviter, could not find peer DID in response message.')
            throw 'failed to get peer DID from inviter'
        }

        // request wallet peer DID from router by sending peer DID from other party
        console.debug('requesting peer DID from wallet')
        let walletDID = await requestDIDFromMediator(this.agent, peerDID)
        console.debug(`walletDID: ${JSON.stringify(walletDID, null, 2)}`)

        console.log('sharing wallet peer DID to inviter')
        // share wallet peer DID to other party
        await this.agent.blindedrouting.sendRegisterRouteRequest({
            messageID: message['@id'],
            didDoc: walletDID,
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
