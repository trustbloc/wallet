/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { GNAPClient } from '@trustbloc/wallet-sdk';
import { computed } from 'vue';
import store from '@/store';

// TODO Issue-364 Add gnap create key pair logic to agent sdk library
export function createKeyPair() {
  return window.crypto.subtle.generateKey(
    {
      name: 'ECDSA',
      namedCurve: 'P-256',
    },
    false,
    ['sign', 'verify']
  );
}

export async function gnapRequestAccess(signer) {
  const gnapAccessTokenConfig = computed(() => store.getters['getGnapAccessTokenConfig']);
  const gnapAccessTokens = await gnapAccessTokenConfig.value;
  const gnapAuthServerURL = computed(() => store.getters['hubAuthURL']).value;
  const gnapWalletCallBackURL = computed(() => store.getters['getStaticAssetsUrl']).value;

  const exportedJwk = await window.crypto.subtle.exportKey('jwk', signer.SignatureVal.publicKey);
  //TODO Issue-1699 persist nonce value and add valid uri in the request
  const nonceVal = 'wallet-nonce';
  const gnapReq = {
    access_token: [gnapAccessTokens],
    client: { key: { jwk: exportedJwk } },
    interact: {
      start: ['redirect'],
      finish: {
        method: 'redirect',
        uri: gnapWalletCallBackURL,
        nonce: nonceVal,
      },
    },
  };

  const gnapClient = new GNAPClient({
    signer: signer,
    gnapAuthServerURL: gnapAuthServerURL,
  });
  const resp = await gnapClient.requestAccess(gnapReq);
  return resp;
}
