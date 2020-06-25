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

        const {domain, challenge, invitation, manifest} = getRequestParams(credEvent);
        this.domain = domain
        this.challenge = challenge
        this.invitation = invitation
        this.manifest = manifest
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

        // TODO verify proof and save
        await this.aries.verifiable.saveCredential({
            name: this.invitation["@id"],
            verifiableCredential: JSON.stringify(this.manifest)
        }).then(() => {
                console.log('successfully saved manifest VC:', name)
            }
        ).catch(err => {
            console.log(`failed to save ${name} : errMsg=${err}`)
            throw err
        })


        let responseData = await this._didConnResponse(connection.result)
        this.sendResponse("VerifiablePresentation", responseData)
    }


    async _didConnResponse(connection) {

        let credential = {
            "@context": [
                "https://www.w3.org/2018/credentials/v1",
                "https://trustbloc.github.io/context/vc/examples-ext-v1.jsonld"
            ],
            issuer: this.walletUser.did,
            issuanceDate: new Date(),
            type: ["VerifiableCredential", "DIDConnection"],
            credentialSubject: {
                id: connection.ConnectionID,
                threadID: connection.ThreadID,
                inviteeDID: connection.MyDID,
                inviterDID: connection.TheirDID,
                inviterLabel: connection.TheirLabel,
                connectionState: connection.State,
            }
        }

        // create did connection VC
        let vc, failure
        await this.aries.verifiable.signCredential({
            credential: credential,
            did: this.walletUser.did,
            signatureType: this.walletUser.signatureType
        }).then(resp => {
                if (!resp.verifiableCredential) {
                    failure = "failed to create did connection credential"
                    return
                }

                vc = resp.verifiableCredential
            }
        ).catch(err => {
            failure = err
        })

        if (failure) {
            console.error("failed to create didconnection credential", failure)
            return failure
        }

        // create did connection response VP
        let presentation = {
            "@context": [
                "https://www.w3.org/2018/credentials/v1"
            ],
            type: "VerifiablePresentation",
            holder: this.walletUser.did,
            verifiableCredential: [vc]
        }

        let data
        await this.aries.verifiable.generatePresentation({
            presentation: presentation,
            domain: this.domain,
            challenge: this.challenge,
            did: this.walletUser.did,
            signatureType: this.walletUser.signatureType,
            skipVerify: true,
        }).then(resp => {
                if (!resp.verifiablePresentation) {
                    data = "failed to create did connection presentation"
                    return
                }

                data = resp.verifiablePresentation
            }
        ).catch(err => {
            data = err
            console.log('failed to create presentation, errMsg:', err)
        })

        return data
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

function getRequestParams(credEvent) {
    if (!credEvent.credentialRequestOptions.web.VerifiablePresentation) {
        throw "invitation not found in did connect credential event"
    }

    const verifiable = credEvent.credentialRequestOptions.web.VerifiablePresentation

    let {challenge, domain, query, invitation, manifest} = verifiable;

    if (query && query.challenge) {
        challenge = query.challenge;
    }

    if (query && query.domain) {
        domain = query.domain;
    }

    if (!domain && credEvent.credentialRequestOrigin) {
        domain = credEvent.credentialRequestOrigin.split('//').pop()
    }

    return {domain, challenge, invitation, manifest}
}
