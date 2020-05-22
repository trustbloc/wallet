/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const dbName = "did-metadata"
const storeName = "metadata"

/**
 * DIDStore is did metadata store
 * @class
 */
export class DIDStore {
    async getDIDMetadata(id) {
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
}