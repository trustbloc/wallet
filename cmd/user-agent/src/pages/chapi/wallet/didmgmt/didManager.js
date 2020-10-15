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
 * DIDManager manages DID create/store/query features
 * @class
 */
export class DIDManager extends KeyValueStore {
    constructor(agent, opts) {
        super(dbName, storeName)

        // params needed for create DID operation
        this.agent = agent
        this.startupOpts = opts
    }

    async createDID(keyType, signType) {
        if (!this.agent) {
            console.error("agent is required to create DIDs")
            throw "operation not supported"
        }

        let generateKeyType = keyTypeIndex.get(keyType)

        const keySet = await this.agent.kms.createKeySet({keyType: generateKeyType})
        const recoveryKeySet = await this.agent.kms.createKeySet({keyType: generateKeyType})
        const updateKeySet = await this.agent.kms.createKeySet({keyType: generateKeyType})

        const createDIDRequest = {
            "publicKeys": [{
                "id": keySet.keyID,
                "type": sigTypeIndex.get(signType),
                "value": keySet.publicKey,
                "encoding": "Jwk",
                "keyType": keyType,
                "purpose": ["general", "auth"]
            }, {
                "id": recoveryKeySet.keyID,
                "type": sigTypeIndex.get(signType),
                "value": recoveryKeySet.publicKey,
                "encoding": "Jwk",
                "keyType": keyType,
                "recovery": true
            }, {
                "id": updateKeySet.keyID,
                "type": sigTypeIndex.get(signType),
                "value": updateKeySet.publicKey,
                "encoding": "Jwk",
                "keyType": keyType,
                "update": true
            }
            ]
        };

        let did
        await this.agent.didclient.createDID(createDIDRequest).then(
            resp => {
                // TODO generate public key from generic wasm
                // TODO pass public key to createDID
                did = resp.DID
            })
            .catch(err => {
                console.error("failed to create did", err)
                throw err
            })

        return did
    }

    async saveDID(user, name, signType, did){
        if (!this.agent) {
            console.error("agent is required for saving DIDs")
            throw "operation not supported"
        }

        // Save DID to local browser storage
        await this.agent.vdri.saveDID({
                name: name,
                did: did
            }
        )

        if (this.startupOpts.sdsServerURL) {
            // Save DID to persistent storage
            // TODO: Deal with SDS sync failures better #328
            await this.agent.didclient.saveDID({
                name: name,
                signType: signType,
                did: did,
                userID: user
            })
        } else {
            console.log("Skipping DID storage to SDS since no SDS server URL was configured.")
        }
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
