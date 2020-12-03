/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {KeyValueStore} from '../common/keyValStore.js'

var uuid = require('uuid/v4')

const dbName = "wallet"
const metadataStore = "metadata"
const manifestStore = "manifest"

/**
 * WalletManager manages create/store/query features for wallet metadata
 * @class
 */
// TODO multiuser support, no need to clear wallet metadata data after logout
export class WalletManager extends KeyValueStore {
    constructor(agent) {
        super(agent, `${dbName}-${metadataStore}`, metadataStore)

        // TODO EDV will be used in future for these stores #268
        this.manifestStore = new KeyValueStore(agent, `${dbName}-${manifestStore}`, manifestStore)
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

    async clear() {
        await super.clear()
        await this.manifestStore.clear()
    }
}
