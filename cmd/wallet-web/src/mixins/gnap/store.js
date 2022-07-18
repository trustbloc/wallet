/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
import { createKeyPair } from '@trustbloc/wallet-sdk';

const dbName = 'gnap-store';
const dbObjectStore = 'gnap-keypair-store';
const storekey = 'gnap';

export async function getGnapKeyPair() {
  const gnapKeyPair = await fetchStoredGnapKeyPair();
  if (gnapKeyPair === undefined) {
    return createAndStoreGnapKeyPair();
  }
  return gnapKeyPair;
}

async function createAndStoreGnapKeyPair() {
  const gnapKeyPair = await createKeyPair('ECDSA', 'P-256', true, ['sign', 'verify']);
  callOnStore(function (store) {
    store.add(gnapKeyPair, storekey);
  });
  return gnapKeyPair;
}

function fetchStoredGnapKeyPair() {
  return new Promise(function (resolve) {
    callOnStore(function (store) {
      store.get(storekey).onsuccess = function (event) {
        return resolve(event.target.result);
      };
    });
  });
}

export function clearGnapStoreData() {
  return new Promise(function (resolve) {
    callOnStore(function (store) {
      store.delete(storekey);
    });
  });
}

function callOnStore(fn_) {
  // This works on all devices/browsers, and uses IndexedDBShim as a final fallback
  const indexedDB =
    window.indexedDB ||
    window.mozIndexedDB ||
    window.webkitIndexedDB ||
    window.msIndexedDB ||
    window.shimIndexedDB;

  const open = indexedDB.open(dbName, 3);
  // Create the schema
  open.onupgradeneeded = function () {
    const db = open.result;
    const store = db.createObjectStore(dbObjectStore);
  };

  open.onsuccess = function () {
    // Start a new transaction
    const db = open.result;
    const tx = db.transaction(dbObjectStore, 'readwrite');
    const store = tx.objectStore(dbObjectStore);
    fn_(store);
    // Close the db when the transaction is done
    tx.oncomplete = function () {
      db.close();
    };
  };
}
