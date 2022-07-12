/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { GNAPClient } from '@trustbloc/wallet-sdk';
import { getGnapKeyPair } from '@/mixins';

export async function gnapRequestAccess(
  signer,
  gnapAccessTokens,
  gnapAuthServerURL,
  walletWebURL,
  nonceVal
) {
  const exportedJwk = await window.crypto.subtle.exportKey('jwk', signer.SignatureVal.publicKey);
  const gnapReq = {
    access_token: [gnapAccessTokens],
    client: { key: { jwk: exportedJwk } },
    interact: {
      start: ['redirect'],
      finish: {
        method: 'redirect',
        uri: walletWebURL + '/gnap/redirect',
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

export async function gnapContinue(signer, gnapAuthServerURL, interactRef, accessToken) {
  const gnapClient = new GNAPClient({
    signer: signer,
    gnapAuthServerURL: gnapAuthServerURL,
  });
  const gnapContinueReq = {
    interact_ref: interactRef,
  };
  const resp = await gnapClient.continue(gnapContinueReq, accessToken);
  return resp;
}
