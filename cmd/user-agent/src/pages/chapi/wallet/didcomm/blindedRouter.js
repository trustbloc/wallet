/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {Messenger} from '../common/messaging'
import {AgentMediator} from './mediator'

var uuid = require('uuid/v4')

const peerDIDRequestMsgType = 'https://trustbloc.github.io/blinded-routing/1.0/diddoc-req'
const peerDIDResponseTopic = 'diddoc-resp'
const sharePeerDIDReqType = 'https://trustbloc.github.io/blinded-routing/1.0/share-diddoc-req'
const sharePeerDIDResTopic = 'share-diddoc-res'

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

        this.messenger = new Messenger(agent)
        this.mediator = new AgentMediator(agent)
        this.mediatorEndpoint = opts.walletMediatorURL
    }

    async sharePeerDID(connection) {
        if (!this.messenger) {
            console.debug('Not sharing router peer DID since blinded routing is disabled !')
            return
        }

        let {ConnectionID} = connection

        // request peer DID from other party
        let responseMsg = await this.messenger.send(ConnectionID, {
            "@id": uuid(),
            "@type": peerDIDRequestMsgType,
            "sent_time": new Date().toJSON(),
        }, {replyTopic: peerDIDResponseTopic})

        let peerDID = _parseResponseDID(responseMsg)
        if (!peerDID) {
            console.error('failed to get peerDID from inviter, could not find peer DID in response message.')
            throw 'failed to get peer DID from inviter'
        }

        // request wallet peer DID from router by sending peer DID from other party
        let walletDID = await this.mediator.requestDID(peerDID)

        // share wallet peer DID to other party
        await this.messenger.reply(responseMsg['@id'], {
            "@id": uuid(),
            "@type": sharePeerDIDReqType,
            data: {didDoc: walletDID},
            "sent_time": new Date().toJSON(),
        }, {replyTopic: sharePeerDIDResTopic, startNewThread: true})

    }
}

let _parseResponseDID = (response) => response.data.didDoc ? response.data.didDoc : undefined
