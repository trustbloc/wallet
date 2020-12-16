/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


import {POST_STATE, waitForEvent} from "../../../../events";
import {getMediatorConnections} from "../../../../pages/chapi/wallet/didcomm/mediator.js"

const stateCompleted = 'completed'
const topicDidExchangeStates = 'didexchange_states'
const stateCompleteMessageTopic = 'didexchange-state-complete'
const stateCompleteMessageType = 'https://trustbloc.dev/didexchange/1.0/state-complete'

/**
 * DIDExchange provides exchange features
 * @param agent instance
 * @class
 */
export class DIDExchange {
    constructor(agent) {
        this.agent = agent
    }

    async connect(invitation, {waitForCompletion = ''} = {}) {
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

        const record = await this.agent.didexchange.queryConnectionByID({id: connID})

        if (waitForCompletion) {
            this.agent.messaging.registerService({
                "name": `${stateCompleteMessageTopic}`,
                "type": `${stateCompleteMessageType}`,
            })

            try {
                await new Promise((resolve, reject) => {
                    setTimeout(() => reject(new Error("time out waiting for state complete message")), 15000)
                    const stop = this.agent.startNotifier(msg => {
                        if (record.result.MyDID == msg.payload.mydid && record.result.TheirDID == msg.payload.theirdid) {
                            stop()
                            console.debug('received state complete msg received !!!!')
                            resolve(msg.payload.message)
                        }
                    }, [stateCompleteMessageTopic])
                })

            } catch (e) {
                console.warn('error while waiting for state complete msg !!', e)
            }
        }

        return record
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
