/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {KeyValueStore} from '../common/keyValStore.js'
import {getMediatorConnections} from '../didcomm/mediator'

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

    async createTrustBlocDID(keyType, signType) {
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
                "purposes": ["authentication"]
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

        // TODO generate public key from generic wasm and pass public key to createDID
        let resp = await this.agent.didclient.createTrustBlocDID(createDIDRequest)

        return resp.DID
    }

    async createPeerDID() {
        let routerConnectionID = await getMediatorConnections(this.agent, true)

        return await this.agent.didclient.createPeerDID({routerConnectionID})
    }

    async saveDID(user, name, signType, did){
        if (!this.agent) {
            console.error("agent is required for saving DIDs")
            throw "operation not supported"
        }

        // Save DID to Aries agent storage
        await this.agent.vdr.saveDID({
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
