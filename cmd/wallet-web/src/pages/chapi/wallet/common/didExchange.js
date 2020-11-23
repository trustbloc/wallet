/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


import {POST_STATE, waitForEvent} from "../../../../events";
import {getMediatorConnections} from "../../../../pages/chapi/wallet/didcomm/mediator.js"

const stateCompleted = 'completed'
const topicDidExchangeStates = 'didexchange_states'

/**
 * DIDExchange provides exchange features
 * @param agent instance
 * @class
 */
export class DIDExchange {
    constructor(agent) {
        this.agent = agent
    }

    async connect(invitation) {
        let conn = await this.agent.outofband.acceptInvitation({
            my_label: 'agent-default-label',
            invitation: invitation,
            router_connections: await getMediatorConnections(this.agent, true),
        })

        let connID = conn['connection_id']
        await waitForEvent(this.agent, {
            type: POST_STATE,
            stateID: stateCompleted,
            connectionID: connID,
            topic: topicDidExchangeStates,
        })

        return await this.agent.didexchange.queryConnectionByID({id: connID})
    }

    cancel() {
        this.sendResponse("Response", "permission denied")
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
