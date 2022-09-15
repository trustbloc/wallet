/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

'use strict';

import { config } from './wdio.shared.conf';

const mockDemoDomain = 'https://demo-adapter.trustbloc.local:8094';

exports.config = {
  ...config,

  // Test files
  specs: [
    './test/specs/oidc-flow.js',
    './test/specs/waci-flow.js',
    './test/specs/credential-interaction-flow.js',
    './test/specs/vault-interaction-flow.js',
  ],

  walletName: 'wallet.trustbloc.local:8091',
  walletURL: 'https://wallet.trustbloc.local:8091',
  walletURLFrench: 'https://wallet.trustbloc.local:8091/fr/',
  authURL: 'https://auth.trustbloc.local:8044',

  // oidc
  oidcDemoVerifierURL: mockDemoDomain + '/verifier/oidc',
  oidcDemoIssuerURL: mockDemoDomain + '/issuer/oidc',

  // waci
  waciDemoVerifierURL: mockDemoDomain + '/verifier/waci',
  waciDemoIssuerURL: mockDemoDomain + '/issuer/waci',

  // chapi
  isCHAPIEnabled: true,
  chapiDemoURL: mockDemoDomain + '/web-wallet',
};
