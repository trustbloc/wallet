/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {DIDManager} from '../didmgmt/didManager.js';
import {getDomainAndChallenge} from '../common/util.js';

/**
 * DIDAuth provides CHAPI did auth features
 * @param aries instance & credential event
 * @class
 */
export class DIDAuth {
    constructor(aries, credEvent) {
        this.aries = aries
        this.credEvent = credEvent
        this.didManager = new DIDManager()

        const {domain, challenge} = getDomainAndChallenge(credEvent);
        this.domain = domain
        this.challenge = challenge
    }

    async getDIDRecords() {
        let issuers = []

        await this.aries.vdri.getDIDRecords().then(
            resp => {
                if (resp.result) {
                    resp.result.forEach((item, id) => {
                        issuers.push({id: id, name: item.name, key: item.id})
                    })
                }
            })
            .catch(err => {
                throw err
            })

        return issuers
    }

    async authorize(did) {
        let didMetadata = await this.didManager.getDIDMetadata(did)

        let data
        await this.aries.verifiable.generatePresentation({
            presentation: {
                "@context": "https://www.w3.org/2018/credentials/v1",
                "type": "VerifiablePresentation"
            },
            domain: this.domain,
            challenge: this.challenge,
            did: did,
            signatureType: didMetadata.signatureType,
            verificationMethod: didMetadata.keyID,
        }).then(resp => {
                if (!resp.verifiablePresentation) {
                    data = "failed to create did auth presentation"
                    return
                }

                data = resp.verifiablePresentation
                //TODO bug in aries to show '"verifiableCredential": null,' in empty presentations
                if (data.hasOwnProperty('verifiableCredential')) {
                    delete data.verifiableCredential
                }
            }
        ).catch(err => {
            data = err
            console.log('failed to create presentation, errMsg:', err)
        })

        console.log("Response presentation:", data)

        // Call Credential Handler callback
        this.sendResponse("VerifiablePresentation", data)
    }

    cancel() {
        this.sendResponse("Response", "DID Auth cancelled")
    }

    sendResponse(type, data){
        this.credEvent.respondWith(new Promise(function (resolve) {
            return resolve({
                dataType: type,
                data: data
            });
        }))
    }
}