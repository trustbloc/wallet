/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import axios from 'axios';
import {POST_STATE, waitForEvent} from "../../../../events";
import {Messenger} from "..";

var uuid = require('uuid/v4')

const routerCreateInvitationPath = `/didcomm/invitation`
const stateCompleted = 'completed'
const topicDidExchangeStates = 'didexchange_states'
const createConnReqType = 'https://trustbloc.dev/blinded-routing/1.0/create-conn-req'
const createConnResTopic = 'create-conn-resp'

/**
 * AgentMediator provides mediator features
 * @param agent instance
 * @class
 */
export class AgentMediator {
    constructor(agent) {
        this.agent = agent
        this.messenger = new Messenger(agent)
    }

    async connect(endpoint) {
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
                if (!e.message.includes("timeout waiting for grant from the router") || i === retries) {
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

    async reconnect() {
        try {
            let res = await this.agent.mediator.getConnections()
            for (const connection of res.connections) {
                await this.agent.mediator.reconnect({connectionID: connection})
            }
        } catch (e) {
            console.error('unable to reconnect to router', e)
        }
    }

    async isAlreadyConnected() {
        let res
        try {
            res = await this.agent.mediator.getConnections()
        } catch (e) {
            throw e
        }

        return res.connections && res.connections.length > 0
    }

    async createInvitation() {
        // creates invitation through the out-of-band protocol
        let response = await this.agent.outofband.createInvitation({
            label: 'agent-label',
            router_connection_id: await getMediatorConnections(this.agent, true)
        })

        return response.invitation
    }

    async requestDID(reqDoc) {
        let connection = await getMediatorConnections(this.agent, true)
        if (!connection) {
            console.error('failed to send connection request to router, no connection found!')
            throw 'could not find connection with router'
        }

        let response = await this.messenger.send(connection, {
            "@id": uuid(),
            "@type": createConnReqType,
            data: {didDoc: reqDoc},
            "sent_time": new Date().toJSON(),
        }, {replyTopic: createConnResTopic})


        // TODO currently getting routerDIDDoc as byte[], to be fixed
        if (response.data.didDoc && response.data.didDoc.length > 0) {
            return JSON.parse(String.fromCharCode.apply(String, response.data.didDoc))
        }

        console.error('failed to request DID from router, failed to get connection response')
        throw 'failed to request DID from router, failed to get connection response'
    }
}

export async function getMediatorConnections(agent, single) {
    let resp = await agent.mediator.getConnections()
    if (!resp.connections || resp.connections.length === 0) {
        return ""
    }

    if (single) {
        return resp.connections[Math.floor(Math.random() * resp.connections.length)]
    }

    return resp.connections.join(",");
}

const createInvitationFromRouter = async (endpoint) => {
    const response = await axios.get(`${endpoint}${routerCreateInvitationPath}`)
    return response.data.invitation
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
