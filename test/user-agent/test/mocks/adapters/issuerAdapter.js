/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {getMediatorConnections} from "../../../../../cmd/user-agent/src/pages/chapi/wallet/didcomm/mediator";
import {waitForEvent} from "../../../../../cmd/user-agent/src/events";

const msgServices = [
    {name: 'request-for-diddoc', type: 'https://trustbloc.github.io/blinded-routing/1.0/diddoc-req'},
    {name: 'share-diddoc-req', type: 'https://trustbloc.github.io/blinded-routing/1.0/share-diddoc-req'},
]

var uuid = require('uuid/v4')

/**
 * Adapter mocks common adapter features
 * @param agent instance
 * @class
 */
class Adapter {
    constructor(agent) {
        this.agent = agent
    }

    /**
     *  init performs initialization operations for adapter
     */
    async init() {
        try {
            for (const msgSvc of msgServices) {
                await this.agent.messaging.registerService(msgSvc)
            }
        } catch (e) {
            console.warn('failed to register message service', e)
        }
    }

    /**
     * acceptExchangeRequest waits for did-exchange event and approves did connection request from adapter
     */
    async acceptExchangeRequest() {
        let res = await waitForEvent(this.agent, {topic: 'didexchange_actions'})

        await this.agent.didexchange.acceptExchangeRequest({
            id: res.Properties.connectionID,
            router_connections: await getMediatorConnections(this.agent, true),
        })
    }

    /**
     * sharePeerDID exchanges peer DID with other party over messages for blinded routing
     *
     * @returns {Promise<Object>} containing DID shared by other party
     */
    async sharePeerDID() {
        // wait for request from wallet for peer DID
        let reqForDID = await waitForEvent(this.agent, {topic: 'request-for-diddoc'})

        // send any sample peer DID to wallet
        let sampleRes = await this.agent.vdri.resolveDID({id: reqForDID.mydid})
        this.agent.messaging.reply({
            "message_ID": reqForDID.message['@id'],
            "message_body": {
                "@id": uuid(),
                "@type": 'https://trustbloc.github.io/blinded-routing/1.0/diddoc-res',
                data: {didDoc: sampleRes.did},
            }
        })

        // wait for did shared by other party
        let sharedDID = await waitForEvent(this.agent, {topic: 'share-diddoc-req'})
        // send acknowledgement to wallet
        this.agent.messaging.reply({
            "message_ID": sharedDID.message['@id'],
            "message_body": {
                "@id": uuid(),
                "@type": 'https://trustbloc.github.io/blinded-routing/1.0/share-diddoc-res'
            }
        })

        return sharedDID.message.data ? sharedDID.message.data.didDoc : null
    }
}

/**
 *
 * IssuerAdapter is mock issuer adapter
 *
 * @param agent instance
 * @class
 */
export class IssuerAdapter extends Adapter {

    /**
     * issueCredential issues given credential from issuer adapter
     */
    async issueCredential(cred) {
        let res = await waitForEvent(this.agent, {topic: 'issue-credential_actions'})
        await this.agent.issuecredential.acceptRequest({
            piid: res.Properties.piid,
            issue_credential: cred
        })
    }

}

/**
 *
 * RPAdapter is mock RP adapter
 *
 * @param agent instance
 * @class
 */
export class RPAdapter extends Adapter {

}
