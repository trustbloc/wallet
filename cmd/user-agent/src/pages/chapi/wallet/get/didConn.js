/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {WalletManager} from '../register/walletManager'

/**
 * DIDConn provides CHAPI did connection/exchange features
 * @param aries instance & credential event
 * @class
 */
export class DIDConn {
    constructor(aries, credEvent, walletUser) {
        this.aries = aries
        this.walletUser = walletUser
        this.walletManager = new WalletManager()
        this.credEvent = credEvent

        const {invitation} = getInvitation(credEvent);
        this.invitation = invitation
    }

    async connect() {
        // perform did exchange
        let res = await this.aries.didexchange.receiveInvitation(this.invitation)

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

        // save wallet metadata
        if (this.walletUser.connections) {
            this.walletUser.connections.push(connection.result)
        } else {
            this.walletUser.connections = [connection.result]
        }
        await this.walletManager.storeWalletMetadata(this.walletUser.id, this.walletUser)


        // TODO response to be signed
        this.sendResponse("VerifiablePresentation", didConnResponse(this.walletUser.did, connection.result))
    }

    async waitFor(connectionID, states, callback, timeout) {
        return new Promise((resolve, reject) => {
            const stop = this.aries.startNotifier(notice => {
                if (connectionID !== notice.payload.connection_id) {
                    return
                }

                if (states && !states.includes(notice.payload.state)) {
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

function getInvitation(credEvent) {
    if (!credEvent.credentialRequestOptions.web.VerifiablePresentation) {
        throw "invitation not found in did connect credential event"
    }

    return credEvent.credentialRequestOptions.web.VerifiablePresentation
}

function didConnResponse(holder, connection) {
    return {
        "@context": "https://www.w3.org/2018/credentials/v1",
        type: "VerifiablePresentation",
        holder: holder,
        verifiableCredential: {
            "@context": [
                "https://www.w3.org/2018/credentials/v1"
            ],
            type: [
                "VerifiableCredential",
                "DIDConnection"
            ],
            credentialSubject: {
                id: connection.ConnectionID,
                threadID: connection.ThreadID,
                inviteeDID: connection.MyDID,
                inviterDID: connection.TheirDID,
                inviterLabel: connection.TheirLabel,
                state: connection.State,
            }
        }
    }
}