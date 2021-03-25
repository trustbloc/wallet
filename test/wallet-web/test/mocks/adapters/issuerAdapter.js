/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {getMediatorConnections, createInvitationFromRouter} from "../../../../../cmd/wallet-web/src/pages/chapi/wallet/didcomm/mediator";
import {POST_STATE, waitForEvent} from "../../../../../cmd/wallet-web/src/events";

const msgServices = [
    {name: 'request-for-diddoc', type: 'https://trustbloc.dev/blinded-routing/1.0/diddoc-req'},
    {name: 'register-route-req', type: 'https://trustbloc.dev/blinded-routing/1.0/register-route-req'},
    {name: 'diddoc-req', type: 'https://trustbloc.dev/adapter/1.0/diddoc-req'}
]

var uuid = require('uuid/v4')

const routerCreateInvitationPath = `/didcomm/invitation`
const stateCompleted = 'completed'
const topicDidExchangeStates = 'didexchange_states'


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

    async connectToMediator(endpoint) {
        let invitation = await createInvitationFromRouter(endpoint)
        let conn = await this.agent.outofband.acceptInvitation({
            my_label: 'agent-default-label',
            invitation: invitation,
        })

        let connID = conn['connection_id']

        await waitForEvent(this.agent, {
            type: POST_STATE,
            stateID: stateCompleted,
            connectionID: connID,
            topic: topicDidExchangeStates,
        })


        const retries = 10;
        for (let i = 1; i <= retries; i++) {
            try {
                await this.agent.mediator.register({connectionID: connID})
            } catch (e) {
                if (!e.message.includes("router registration : get grant: store: data not found") || i === retries) {
                    throw e
                }
                await sleep(500);
                continue
            }
            break
        }

        let res = await this.agent.mediator.getConnections()

        if (res.connections.includes(connID)) {
            console.log("router registered successfully!", connID)
        } else {
            console.log("router was not registered!", connID)
        }

        // return handle for disconnect
        return () => this.agent.mediator.unregister({connectionID: connID})
    }

    async createInvitation() {
        let response = await this.agent.mediatorclient.createInvitation({
            label: 'mockadapter-label'
        })

        return response.invitation
    }


    /**
     * acceptExchangeRequest waits for did-exchange event and approves did connection request from adapter
     */
    async acceptExchangeRequest() {
        let res = await waitForEvent(this.agent, {topic: 'didexchange_actions', id: '{issuer/rp adapter}'})

        await this.agent.didexchange.acceptExchangeRequest({
            id: res.Properties.connectionID,
            router_connections: await getMediatorConnections(this.agent, true),
        })

        await waitForEvent(this.agent, {
            type: POST_STATE,
            stateID: stateCompleted,
            connectionID: res.Properties.connectionID,
            topic: topicDidExchangeStates,
        })

        this.agent.messaging.send({"connection_ID": `${res.Properties.connectionID}`, "message_body": {
                "@id": uuid(),
                "@type": 'https://trustbloc.dev/didexchange/1.0/state-complete',
            }})
    }

    /**
     * sharePeerDID exchanges peer DID with other party over messages for blinded routing
     *
     * @returns {Promise<Object>} containing DID shared by other party
     */
    async sharePeerDID() {
        // wait for request from wallet for peer DID
        let reqForDID = await waitForEvent(this.agent, {topic: 'request-for-diddoc',
            timeoutError: 'timeout waiting for peer DID request'})

        console.debug(`adapter.sharePeerDID: reqForDIDL ${JSON.stringify(reqForDID, null, 2)}`)

        // send any sample peer DID to wallet
        let sampleRes = await this.agent.vdr.resolveDID({id: reqForDID.mydid})

        console.debug(`adapter.sharePeerDID: sampleRes: ${JSON.stringify(sampleRes, null, 2)}`)

        this.agent.messaging.reply({
            "message_ID": reqForDID.message['@id'],
            "message_body": {
                "@id": uuid(),
                "@type": 'https://trustbloc.dev/blinded-routing/1.0/diddoc-resp',
                data: {didDoc: sampleRes.did},
            }
        })

        // wait for did shared by other party
        let sharedDID = await waitForEvent(this.agent, {topic: 'register-route-req',
            timeoutError: 'timeout waiting for register route request'})

        console.debug(`adapter.sharePeerDID: sharedDID: ${JSON.stringify(sharedDID, null, 2)}`)

        // send acknowledgement to wallet
        this.agent.messaging.reply({
            "message_ID": sharedDID.message['@id'],
            "message_body": {
                "@id": uuid(),
                "@type": 'https://trustbloc.dev/blinded-routing/1.0/register-route-resp'
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
        let res = await waitForEvent(this.agent, {topic: 'issue-credential_actions',
            timeoutError: 'timeout waiting for issue credential request'})
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

    /**
     * shareNewPeerDID shares new peer DID with wallet for DID comm
     */
    async shareNewPeerDID() {
        // wait for DID doc request from wallet
        let request = await waitForEvent(this.agent, {topic: 'diddoc-req',
            timeoutError: 'timeout waiting request for new peer DID'})

        console.debug(`rpAdapter.shareNewPeerDID: request: ${JSON.stringify(request, null, 2)}`)

        let sampleRes = await this.agent.vdr.resolveDID({id: request.mydid})

        console.debug(`rpAdapter.shareNewPeerDID: sampleRes: ${JSON.stringify(sampleRes, null, 2)}`)

        // send new peer DID to wallet
        this.agent.messaging.reply({
            "message_ID": request.message['@id'],
            "message_body": {
                "@id": uuid(),
                "@type": 'https://trustbloc.dev/adapter/1.0/diddoc-resp',
                data: {didDoc: sampleRes.did}
            }
        })
    }
}
