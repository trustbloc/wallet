/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


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
        console.log("performing did exchange for inviation", invitation)
        // perform did exchange
        let res = await this.aries.didexchange.receiveInvitation(invitation)

        let aries = this.aries
        await this.waitFor(res.connection_id, ['invited'], function () {
            return aries.didexchange.acceptInvitation({
                id: res.connection_id
            })
        })

        // wait for status to be completed
        let connectionID
        try {
            let completed = await this.waitFor(res.connection_id, ['completed'], null, 10000)
            connectionID = completed.connection_id
        } catch (e) {
            // do not fail if connection is not yet completed, return current state in response
            if (!e.toString().includes("time out while waiting for connection")) {
                throw e
            }
            console.error("time out while waiting for connection to be completed")
            connectionID = res.connection_id
        }

        let connection = await this.aries.didexchange.queryConnectionByID({id: connectionID})
        return connection
    }

    async waitFor(connectionID, states, callback, timeout) {
        return new Promise((resolve, reject) => {
            const stop = this.aries.startNotifier(notice => {
                const event = notice.payload
                if (connectionID !== event.Properties.connectionID) {
                    return
                }

                if (states && !states.includes(event.StateID)) {
                    return
                }

                if (event.Type !== "post_state") {
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
                    resolve(notice.payload)
                }

            }, ["all"])

            setTimeout(() => {
                stop()
                reject(new Error("time out while waiting for connection"))
            }, timeout ? timeout : 5000)
        })
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