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

export async function initiateGnapAuth(store, router) {
  const gnapAccessTokens = await store.getters['getGnapAccessTokenConfig'];
  const gnapAuthServerURL = store.getters.hubAuthURL;
  const walletWebUrl = store.getters.walletWebUrl;
  const gnapKeyPair = await getGnapKeyPair();
  console.log('gnapKeyPair', gnapKeyPair);
  const signer = { SignatureVal: gnapKeyPair };
  console.log('signer', signer.SignatureVal);
  const clientNonceVal = (Math.random() + 1).toString(36).substring(7);

  console.log(
    'requestAccess with:',
    JSON.stringify(
      {
        ...signer,
        gnapAccessTokens,
        gnapAuthServerURL,
        walletWebUrl,
        clientNonceVal,
      },
      null,
      2
    )
  );

  const resp = await gnapRequestAccess(
    signer,
    gnapAccessTokens,
    gnapAuthServerURL,
    walletWebUrl,
    clientNonceVal
  );
  console.log('resp', resp);
  // If user have already signed in then just redirect
  if (resp.data.access_token) {
    const accessToken = resp.data.access_token[0].value;
    const subjectId = resp.data.subject.sub_ids[0].id;
    console.log('111 new accessToken', accessToken);
    console.log('111 new subjectId', subjectId);
    store.dispatch('updateSessionToken', accessToken);
    store.dispatch('updateSubjectId', subjectId);
    router.push({ name: 'DashboardLayout' });
    return;
  }
  const respMetaData = {
    uri: resp.data.continue.uri,
    continue_access_token: resp.data.continue.access_token,
    finish: resp.data.interact.finish,
    clientNonceVal: clientNonceVal,
  };
  await store.dispatch('updateGnapReqAccessResp', respMetaData);
  console.log('redirect to:', resp.data.interact.redirect);
  return resp.data.interact.redirect;
}
