/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


import {POST_STATE, waitForEvent} from "../../../../events";

const stateCompleted = 'completed'
const topicDidExchangeStates = 'didexchange_states'

/**
 * DIDExchange provides exchange features
 * @param aries instance
 * @class
 */
export class DIDExchange {
    constructor(aries) {
        this.aries = aries
    }

    async connect(invitation) {
        let conn = await this.aries.outofband.acceptInvitation({
            my_label: 'agent-default-label',
            invitation: invitation,
        })

        let connID = conn['connection_id']
        await waitForEvent(this.aries, {
            type: POST_STATE,
            stateID: stateCompleted,
            connectionID: connID,
            topic: topicDidExchangeStates,
        })

        return await this.aries.didexchange.queryConnectionByID({id: connID})
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