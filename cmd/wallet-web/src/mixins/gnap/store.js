/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
import { createKeyPair } from '@trustbloc/wallet-sdk';

const dbName = 'gnap-store';
const dbObjectStore = 'gnap-keypair-store';
const storekey = 'gnap';

export async function getGnapKeyPair(kid, alg) {
  const gnapKeyPair = await fetchStoredGnapKeyPair();
  if (gnapKeyPair === undefined) {
    if (!kid)
      throw new Error('Error getting GNAP keypair: keypair does not exist and kid is missing');
    if (!alg)
      throw new Error('Error getting GNAP keypair: keypair does not exist and alg is missing');
    return createAndStoreGnapKeyPair(kid, alg);
  }
  return gnapKeyPair;
}

async function createAndStoreGnapKeyPair(kid, alg) {
  if (!kid) throw new Error('Error getting GNAP keypair: kid is missing');
  if (!alg) throw new Error('Error getting GNAP keypair: alg is missing');
  const gnapKeyPair = await createKeyPair('ECDSA', 'P-256', true, ['sign', 'verify']);
  Object.assign(gnapKeyPair, { kid, alg });
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
