/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {KeyValueStore} from '../common/keyValStore.js'

var uuid = require('uuid/v4')

const dbName = "wallet"
const metadataStore = "metadata"
const manifestStore = "manifest"
const connectionsStore = "connections"

/**
 * WalletManager manages create/store/query features for wallet metadata
 * @class
 */
// TODO multiuser support, no need to clear wallet metadata data after logout
export class WalletManager extends KeyValueStore {
    constructor() {
        super(`${dbName}-${metadataStore}`, metadataStore)

        // TODO EDV will be used in future for these stores #268
        this.manifestStore = new KeyValueStore(`${dbName}-${manifestStore}`, manifestStore)
        this.connectionStore = new KeyValueStore(`${dbName}-${connectionsStore}`, connectionsStore)
    }

    async getAllWalletMetadata() {
        return this.getAll()
    }

    async getWalletMetadata(user) {
        return this.get(user)
    }

    async storeWalletMetadata(user, metadata) {
        return this.store(user, metadata)
    }

    async storeManifest(connectionID, manifest) {
        let id = (manifest.id) ? manifest.id : uuid()
        manifest.connection = connectionID
        await this.manifestStore.store(id, manifest)
    }

    async getAllManifests() {
        return this.manifestStore.getAll()
    }

    async storeConnection(connectionID, connection) {
        return this.connectionStore.store(connectionID, connection)
    }

    async getAllConnections() {
        return this.connectionStore.getAll()
    }

    async getConnectionByID(id) {
        return this.connectionStore.get(id)
    }

    async clear() {
        await super.clear()
        await this.connectionStore.clear()
        await this.manifestStore.clear()
    }
}
