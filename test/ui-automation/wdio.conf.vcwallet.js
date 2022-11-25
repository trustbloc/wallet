/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

'use strict';

import { config } from './wdio.shared.conf';

const mockDemoDomain = 'https://demo-adapter.trustbloc.local:8094';

exports.config = {
  ...config,

  // Test files
  specs: [
    './test/specs/openid4vc-flow.js',
    './test/specs/oidc-flow.js',
    './test/specs/vault-interaction-flow.js',
  ],

  walletName: 'vcwallet.trustbloc.local:8071',
  walletURL: 'https://vcwallet.trustbloc.local:8071',
  walletURLFrench: 'https://vcwallet.trustbloc.local:8071/fr/',
  authURL: 'https://auth.trustbloc.local:8044',

  // openid4vc
  openid4vpInitiateRequestURL: `openid-vc://?request_uri=${mockDemoDomain}/verifier/openid4vc/share`,

  // oidc
  oidcDemoVerifierURL: mockDemoDomain + '/verifier/oidc',
  oidcDemoIssuerURL: mockDemoDomain + '/issuer/oidc',
};
