/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import axios from 'axios';

const routerCreateInvitationPath = `/connections/create-invitation`

/**
 * AgentMediator provides aries mediator features
 * @param aries agent instance
 * @class
 */
export class AgentMediator {
    constructor(aries) {
        this.aries = aries
    }

    async connect(endpoint) {
        let invitation = await createInvitationFromRouter(endpoint)
        let aries = this.aries
        let conn = await aries.didexchange.receiveInvitation(invitation)

        await this.waitFor(conn.connection_id, 'invited', function () {
            return aries.didexchange.acceptInvitation({
                id: conn.connection_id
            })
        })

        await this.waitFor(conn.connection_id, 'completed', function () {
            return aries.mediator.register({connectionID: conn.connection_id})
        })

        let res = await aries.mediator.getConnection().catch(err => {
            if (!err.message.includes("router not registered")) {
                throw err
            }
        })

        console.log("router registered successfully..!!", res.connectionID)

        // return handle for disconnect
        return () => aries.mediator.unregister()
    }

    async isAlreadyConnected() {
        let res
        try {
            res = await this.aries.mediator.getConnection()
        } catch (e) {
            if (e.toString().includes("router not registered")){
                return false
            }

            throw e
        }

        return res.connectionID != ""
    }

    async createInvitation() {
        let response = await this.aries.didexchange.createInvitation()
        return response.invitation
    }

    async waitFor(connectionID, state, callback) {
        return new Promise((resolve, reject) => {
            const stop = this.aries.startNotifier(notice => {
                if (connectionID !== notice.payload.connection_id) {
                    return
                }

                if (state && notice.payload.state !== state) {
                    return
                }

                stop()

                if (callback) {
                    try {
                        callback().then(() => {
                            resolve()
                        })
                    } catch (err) {
                        reject(err)
                    }
                } else {
                    resolve()
                }

            }, ["all"])

            setTimeout(() => {
                stop()
                reject(new Error("time out while waiting for connection"))
            }, 5000)
        })
    }
}

const createInvitationFromRouter = async (endpoint) => {
    const response = await axios.post(`${endpoint}${routerCreateInvitationPath}`)
    return response.data.invitation
}