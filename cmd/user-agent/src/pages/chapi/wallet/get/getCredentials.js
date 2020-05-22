/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {DIDAuth} from "./didAuth";
import {getCredentialType, isCredentialType} from '../common/util.js';

const responseType = "VerifiablePresentation"

/**
 * WalletGet provides CHAPI get credential features
 * @param aries instance & credential event
 * @class
 */
export class WalletGet extends DIDAuth {
    constructor(aries, credEvent) {
        super(aries, credEvent);
        this.setReasonAndSearchKey()
    }

    setReasonAndSearchKey() {
        const verifiable = this.credEvent.credentialRequestOptions.web.VerifiablePresentation

        let query = {}
        if (verifiable.query) {
            // supporting only one query for now
            query = Array.isArray(verifiable.query) ? verifiable.query[0] : verifiable.query;
        }

        if (query.credentialQuery && query.credentialQuery.reason) {
            this.reason = query.credentialQuery.reason
        }

        this.search = ""
        if (query.credentialQuery && query.credentialQuery.example && query.credentialQuery.example.type) {
            let t = query.credentialQuery.example.type
            let key = Array.isArray(t) ? t[0] : t
            if (!isCredentialType(key)) {
                this.search = key
            }
        }
    }

    async getCredentialRecords() {
        let vcs = []

        await this.aries.verifiable.getCredentials().then(
            resp => {
                if (resp.result) {
                    resp.result.forEach((item, id) => {
                        vcs.push({
                            id: id,
                            name: item.name,
                            key: item.id,
                            type: getCredentialType(item.type),
                            holder: item.subjectId,
                        })
                    })
                }
            })
            .catch(err => {
                throw err
            })

        return vcs
    }

    async createAndSendPresentation(did, selections) {
        try {
            let vcs = []
            for (let selectedVC of selections) {
                const resp = await this.aries.verifiable.getCredential({
                    id: selectedVC.key
                })
                vcs.push(JSON.parse(resp.verifiableCredential))
            }

            let didMetadata = await this.didStore.getDIDMetadata(did)

            let data
            await this.aries.verifiable.generatePresentation({
                verifiableCredential: vcs,
                did: did,
                domain: this.domain,
                challenge: this.challenge,
                skipVerify: true,
                signatureType: didMetadata.signatureType,
                verificationMethod: didMetadata.keyID
            }).then(resp => {
                    data = resp.verifiablePresentation
                }
            )

            this.sendResponse(responseType, data)
        } catch (e) {
            this.sendResponse("error", e)
        }

    }
}