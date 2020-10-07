/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {Messenger} from '../common/messaging'
import {AgentMediator} from './mediator'

var uuid = require('uuid/v4')

const requestPeerDIDMsgType = 'request-peer-did'
const sendPeerDIDMsgType = 'send-peer-did'

/**
 * BlindedRouter provides blinded routing features
 * @param aries agent instance
 * @class
 */
export class BlindedRouter {
    constructor(aries, opts) {
        console.log('opts in blinded router', opts)
        if (!opts.blindedRouting) {
            return
        }

        this.messenger = new Messenger(aries)
        this.mediator = new AgentMediator(aries)
        this.mediatorEndpoint = opts.walletMediatorURL
    }

    async sharePeerDID(connection) {
        if (!this.messenger) {
            console.debug('Not sharing router peer DID since blinded routing is disabled !')
            return
        }

        // send message to connection requesting peer DID
        // TODO sending empty message for now, to be implemented
        this.messenger.sendAndWaitForReply(connection, {
            "@id": uuid(),
            "@type": requestPeerDIDMsgType,
            "~l10n": {"locale": "en"},
            "sent_time": new Date().toJSON(),
        })

        // send message to router requesting peerDID
        // TODO response peer DID from previous step to be sent to router
        this.mediator.requestDID()

        // end router peerDID to connection
        // TODO sending empty message for now, to be implemented (sending peer DID from router)
        this.messenger.send(connection, {
            "@id": uuid(),
            "@type": sendPeerDIDMsgType,
            "~l10n": {"locale": "en"},
            "sent_time": new Date().toJSON(),
        })

    }
}
