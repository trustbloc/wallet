/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


/**
 * KeyValueStore is key value based store
 * @class
 */
export class KeyValueStore {
    constructor(dbName, storeName) {
        if (!dbName || !storeName) {
            throw "db name and store name is mandatory"
        }

        this.dbName = dbName
        this.storeName = storeName
    }

    async getAll() {
        const dbName = this.dbName
        const storeName = this.storeName

        return new Promise(function (resolve) {
            let openDB = indexedDB.open(dbName, 1);

            openDB.onupgradeneeded = function () {
                var db = {}
                db.result = openDB.result;
                db.store = db.result.createObjectStore(storeName, {keyPath: "id"});
            };

            openDB.onsuccess = function () {
                let db = {};
                db.result = openDB.result;
                db.tx = db.result.transaction(storeName, "readonly");
                db.store = db.tx.objectStore(storeName);
                let request = db.store.getAll();
                request.onsuccess = function () {
                    resolve(request.result)
                };
            }
        });
    }

    async get(id) {
        const dbName = this.dbName
        const storeName = this.storeName

        return new Promise(function (resolve) {
            var openDB = indexedDB.open(dbName, 1);

            openDB.onupgradeneeded = function () {
                var db = {}
                db.result = openDB.result;
                db.store = db.result.createObjectStore(storeName, {keyPath: "id"});
            };

            openDB.onsuccess = function () {
                var db = {};
                db.result = openDB.result;
                db.tx = db.result.transaction(storeName, "readonly");
                db.store = db.tx.objectStore(storeName);
                let getData = db.store.get(id);
                getData.onsuccess = function () {
                    resolve(getData.result);
                };

                db.tx.oncomplete = function () {
                    db.result.close();
                };
            }
        });
    }

    async store(id, value) {
        const dbName = this.dbName
        const storeName = this.storeName

        return new Promise(function (resolve){
            var openDB = indexedDB.open(dbName, 1);

            openDB.onupgradeneeded = function () {
                var db = {}
                db.result = openDB.result;
                db.store = db.result.createObjectStore(storeName, {keyPath: "id"});
            };

            if (!value.id) {
                value.id = id
            }

            openDB.onsuccess = function () {
                console.log("save on success")
                var db = {};
                db.result = openDB.result;
                db.tx = db.result.transaction(storeName, "readwrite");
                db.store = db.tx.objectStore(storeName);
                db.store.put(value);
                resolve()
            }
        });
    }

    async clear() {
        const dbName = this.dbName
        const storeName = this.storeName

        return new Promise(function (resolve) {
            let openDB = indexedDB.open(dbName, 1);

            openDB.onupgradeneeded = function () {
                var db = {}
                db.result = openDB.result;
                db.store = db.result.createObjectStore(storeName, {keyPath: "id"});
            };

            openDB.onsuccess = function () {
                let db = {};
                db.result = openDB.result;
                db.tx = db.result.transaction(storeName, "readwrite");
                db.store = db.tx.objectStore(storeName);
                let request = db.store.clear();
                request.onsuccess = function () {
                    resolve()
                };
            }
        });
    }


}