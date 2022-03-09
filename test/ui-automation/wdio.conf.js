/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

import {config} from "./wdio.shared.conf";

const mockDemoDomain = "https://demo-adapter.trustbloc.local:8094"


exports.config = {
  ...config,
  // TODO: changed this to url - fix wallet name issue
  walletName: "wallet.trustbloc.local",
  walletURL: "https://wallet.trustbloc.local:8091",
  walletURLFrench: "https://wallet.trustbloc.local:8091/fr/",
  waciDemoVerifierURL:  mockDemoDomain + "/verifier/waci",
  waciDemoIssuerURL: mockDemoDomain + "/issuer/waci",
  chapiDemoURL: mockDemoDomain + "/web-wallet",
};
