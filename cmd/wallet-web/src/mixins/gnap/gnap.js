/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { GNAPClient } from '@trustbloc/wallet-sdk';

export async function gnapRequestAccess(signer, gnapAccessTokens, gnapAuthServerURL, nonceVal) {
  const exportedJwk = await window.crypto.subtle.exportKey('jwk', signer.SignatureVal.publicKey);
  const gnapReq = {
    access_token: [gnapAccessTokens],
    client: { key: { jwk: exportedJwk } },
    interact: {
      start: ['redirect'],
      finish: {
        method: 'redirect',
        uri: '/',
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
