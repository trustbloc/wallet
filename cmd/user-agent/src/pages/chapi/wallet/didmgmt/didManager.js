/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {KeyValueStore} from '../common/keyValStore.js'

const dbName = "did-metadata"
const storeName = "metadata"
const sigTypeIndex = new Map([["Ed25519Signature2018", "Ed25519VerificationKey2018"], ["JsonWebSignature2020", "JwsVerificationKey2020"]]);
const keyTypeIndex = new Map([["Ed25519", "ED25519"], ["P256", "ECDSAP256IEEEP1363"]]);

/**
 * DIDManager is manages DID create/store/query features
 * @class
 */
export class DIDManager extends KeyValueStore {
    constructor(aries, trustblocAgent, opts) {
        super(dbName, storeName)

        // params needed for create DID operation
        this.aries = aries
        this.trustblocAgent = trustblocAgent
        this.trustblocStartupOpts = opts
    }

    async createDID(keyType, signType) {
        if (!this.aries || !this.trustblocAgent) {
            console.error("aries and trustbloc agents are required to created DIDs")
            throw "operation not supported"
        }

        let generateKeyType = keyTypeIndex.get(keyType)

        const keySet = await this.aries.kms.createKeySet({keyType: generateKeyType})
        const recoveryKeySet = await this.aries.kms.createKeySet({keyType: generateKeyType})
        const opsKeySet = await this.aries.kms.createKeySet({keyType: generateKeyType})

        const createDIDRequest = {
            "publicKeys": [{
                "id": opsKeySet.keyID,
                "type": sigTypeIndex.get("JsonWebSignature2020"),
                "value": opsKeySet.publicKey,
                "encoding": "Jwk",
                "keyType": keyType,
                "usage": ["ops"]
            }, {
                "id": keySet.keyID,
                "type": sigTypeIndex.get(signType),
                "value": keySet.publicKey,
                "encoding": "Jwk",
                "keyType": keyType,
                "usage": ["general", "auth"]
            }, {
                "id": recoveryKeySet.keyID,
                "type": sigTypeIndex.get(signType),
                "value": recoveryKeySet.publicKey,
                "encoding": "Jwk",
                "keyType": keyType,
                "recovery": true
            }
            ]
        };

        const t = await new this.trustblocAgent.Framework(JSON.parse(this.trustblocStartupOpts))

        let did
        await t.didclient.createDID(createDIDRequest).then(
            resp => {
                // TODO generate public key from generic wasm
                // TODO pass public key to createDID
                did = resp.DID

            })
            .catch(err => {
                t.destroy()
                console.error("failed to create did", err)
                throw err
            })

        await t.destroy()

        return did
    }

    async saveDID(name, did){
        if (!this.aries) {
            console.error("aries agent required for saving DIDs")
            throw "operation not supported"
        }

        await this.aries.vdri.saveDID({
                name: name,
                did: did
            }
        )
    }

    async getAllDIDMetadata() {
        return this.getAll()
    }

    async getDIDMetadata(did) {
        return this.get(did)
    }

    async storeDIDMetadata(did, metadata) {
        return this.store(did, metadata)
    }

}