/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import axios from 'axios';

const format = 'ldp_vc';
const responseType = 'code';
const scope = 'openid';

// TODO hardcoded for now, to be read dynamically or to be configured, Issue #1531
const clientID = 'm1CppYUvt7';

/**
 * Reads Open ID connect server metadata from wellknown URL of given Open ID server.
 */
export async function readOpenIDConfiguration(server) {
  const { status, data } = await axios.get(`${server}/.well-known/openid-configuration`);

  if (status != 200) {
    throw 'failed to read OIDC configuration';
  }

  return data;
}

/**
 * Performs OIDC issuer authorize redirect by reading give OIDC server configuration, initiate issuance params.
 */
export function sendCredentialAuthorizeRequest(configuration, initiateRequest, callback, state) {
  if (!configuration.authorization_endpoint) {
    throw "unable to read 'authorization_endpoint' from OIDC server configuration";
  }

  const claims = [];

  if (initiateRequest.credential_type) {
    const types = Array.isArray(initiateRequest.credential_type)
      ? initiateRequest.credential_type
      : [initiateRequest.credential_type];
    types.forEach((type) => claims.push({ type, format }));
  } else if (initiateRequest.manifest_id) {
    const manifestIDs = Array.isArray(initiateRequest.manifest_id)
      ? initiateRequest.manifest_id
      : [initiateRequest.manifest_id];
    manifestIDs.forEach((manifest_id) => claims.push({ manifest_id, format }));
  } else {
    throw "failed to prepare claim request, 'credential_type' or 'manifest_id' expected in initiate request";
  }

  console.debug('reading issuer OIDC configuration', JSON.stringify(configuration, null, 2));
  console.debug('initiate credenital save request', JSON.stringify(initiateRequest, null, 2));

  const href = encodeURI(`${configuration.authorization_endpoint}?claims=${encodeURIComponent(
    JSON.stringify(claims)
  )}&response_type=${responseType}
  &client_id=${clientID}&scope=${scope}&state=${state}&redirect_uri=${encodeURIComponent(
    callback
  )}`);

  window.location.href = href;
}

export async function requestToken(
  endpoint,
  { redirect_uri, code, grant_type = 'authorization_code' }
) {
  const config = {
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
  };

  var formPayload = new FormData();
  formPayload.append('redirect_uri', redirect_uri);
  formPayload.append('code', code);
  formPayload.append('grant_type', grant_type);

  // TODO authorization header to be added
  const { status, data } = await axios.post(endpoint, formPayload, config);

  if (status != 200) {
    throw 'token request failed';
  }

  return data;
}

export async function requestCredential(
  endpoint,
  { access_token, token_type = 'Bearer', credentialType, format = 'ldp_vc' }
) {
  const config = {
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      Authorization: `${token_type} ${access_token}`,
    },
  };

  var formPayload = new FormData();
  formPayload.append('type', credentialType);
  formPayload.append('format', format);

  // TODO authorization header to be added
  const { status, data } = await axios.post(endpoint, formPayload, config);

  if (status != 200) {
    throw 'credential request failed';
  }

  return data;
}
