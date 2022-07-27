/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import axios from 'axios';
import { GNAPClient } from '@trustbloc/wallet-sdk';
import { exportJWKGnapPublicKey } from '@/mixins/gnap/store';

export async function gnapRequestAccess(
  signer,
  gnapAccessTokens,
  gnapAuthServerURL,
  walletWebURL,
  nonceVal
) {
  const signingKey = await exportJWKGnapPublicKey();
  const gnapReq = {
    access_token: [gnapAccessTokens],
    client: { key: { jwk: signingKey, proof: signer.proofType() } },
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

export async function getBootstrapData(agentOpts, hubAuthURL, dispatch, accessToken) {
  // Prop Validation
  if (!agentOpts) throw new Error(`Error getting bootstrap data: agentOpts can't be empty`);
  if (!hubAuthURL) throw new Error(`Error getting bootstrap data: hubAuthURL can't be empty`);
  if (!dispatch) throw new Error(`Error getting bootstrap data: dispatch is undefined`);
  if (!accessToken) throw new Error(`Error getting bootstrap data: accessToken can't be empty`);

  const newAgentOpts = {};

  const newOpts = await axios
    .get(hubAuthURL + '/gnap/bootstrap', {
      headers: { Authorization: `GNAP ${accessToken}` },
    })
    .then((resp) => {
      const { data } = resp;

      // TODO to be removed after universal wallet migration
      if (agentOpts.storageType === 'edv') {
        Object.assign(newAgentOpts, { edvVaultID: data?.data?.userEDVVaultID });
        // TODO this property is not returned from the bootstrap data - remove if not needed
        Object.assign(newAgentOpts, { edvCapability: data?.data?.edvCapability });
      }

      // TODO to be removed after universal wallet migration
      if (agentOpts.kmsType === 'webkms') {
        Object.assign(newAgentOpts, { opsKeyStoreURL: data?.data?.opsKeyStoreURL });
        Object.assign(newAgentOpts, { edvOpsKIDURL: data?.data?.edvOpsKIDURL });
        Object.assign(newAgentOpts, { edvHMACKIDURL: data?.data?.edvHMACKIDURL });
      }

      // TODO to be removed after universal wallet migration
      // TODO this property is not returned from the bootstrap data - remove if not needed
      Object.assign(newAgentOpts, { authzKeyStoreURL: data?.data?.authzKeyStoreURL });
      // TODO this property is not returned from the bootstrap data - remove if not needed
      Object.assign(newAgentOpts, { opsKMSCapability: data?.data?.opsKMSCapability });

      return { newAgentOpts, newProfileOpts: { bootstrap: data } };
    })
    .catch((e) => {
      // http 403 denotes that user is not authenticated
      if (e.response && e.response.status === 403) {
        console.debug(e);
      }
      // http 400 denotes expired cookie at server - logout the user and prompt user to authenticate
      // TODO: also logout user if access token has expired - initiate auth flow from the beginning
      else if (e.response && e.response.status === 400) {
        dispatch('agent/logout');
      } else {
        console.error('Error getting bootstrap data:', e);
      }
    });
  return newOpts;
}
