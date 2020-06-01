/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {KeyValueStore} from '../common/keyValStore.js'

const dbName = "wallet-metadata"
const storeName = "metadata"

/**
 * WalletManager manages create/store/query features for wallet metadata
 * @class
 */
export class WalletManager extends KeyValueStore {
    constructor() {
        super(dbName, storeName)
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

    async getRegisteredUser() {
        let result = await this.getAllWalletMetadata()
        if (!result || result.length == 0) {
            return
        }

        return result[0]
    }

}