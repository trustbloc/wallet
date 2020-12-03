/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


/**
 * KeyValueStore is key value based store
 * @class
 */
export class KeyValueStore {
    constructor(agent, dbName, storeName) {
        if (!dbName || !storeName) {
            throw "db name and store name is mandatory"
        }

        this.agent = agent
        this.dbName = dbName
        this.storeName = storeName
    }

    async getAll() {
        let data = await this.agent.store.iterator({
            start_key: this.dbName + this.storeName,
            end_key: this.dbName + this.storeName + "!!",
        })

        let results = [];

        if (!data.results) {
            return results
        }

        for (const val of data.results) {
            results.push(JSON.parse(atob(val)))
        }

        return results
    }

    async get(id) {
        let _this = this;
        return new Promise(function (resolve, reject) {
            _this.agent.store.get({
                key: _this.dbName + _this.storeName + id,
            }).then(data => {
                resolve(JSON.parse(atob(data.result)))
            }).catch(err => {
                if (err.message.includes("data not found")) {
                    resolve({})
                    console.warn(err)

                    return
                }

                reject(err)
            })
        });
    }

    async store(id, value) {
        let resp = this.agent.store.put({
            key: this.dbName + this.storeName + id,
            value: btoa(JSON.stringify(value)),
        })
        // need to flush
        await this.agent.store.flush()

        return resp
    }

    async clear() {
    }
}
