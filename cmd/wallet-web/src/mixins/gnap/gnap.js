/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { GNAPClient } from '@trustbloc/wallet-sdk';

// TODO Issue-364 Add gnap create key pair logic to agent sdk library
export async function createKeyPair() {
  const keyPair = await window.crypto.subtle.generateKey(
    {
      name: 'ECDSA',
      namedCurve: 'P-256',
    },
    true,
    ['sign', 'verify']
  );

  return keyPair;
}

export async function gnapRequestAccess(signer, gnapAuthServerURL) {
  // TODO Issue-1709 pass the valid request to agent sdk request access
  let req;
  const gnapClient = new GNAPClient({
    signer: signer,
    gnapAuthServerURL: gnapAuthServerURL,
  });
  const { resp } = await gnapClient.requestAccess(req);
  return resp;
}
