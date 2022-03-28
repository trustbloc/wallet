/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

import { config } from "./wdio.shared.conf";

const mockDemoDomain = "https://demo-adapter.trustbloc.local:8094";

exports.config = {
  ...config,
  walletName: "wallet.trustbloc.local:8091",
  walletURL: "https://wallet.trustbloc.local:8091",
  walletURLFrench: "https://wallet.trustbloc.local:8091/fr/",

  // oidc
  oidcDemoVerifierURL: mockDemoDomain + "/verifier/oidc",
  oidcDemoVerifierCallbackURL: mockDemoDomain + "/verifier/oidc/share/cb",

  // waci
  waciDemoVerifierURL: mockDemoDomain + "/verifier/waci",
  waciDemoIssuerURL: mockDemoDomain + "/issuer/waci",

  // chapi
  chapiDemoURL: mockDemoDomain + "/web-wallet",
};
