/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

import {config} from "./wdio.shared.conf";

exports.config = {
  ...config,
  walletName: "TrustBloc Wallet",
  walletURL: "https://wallet.trustbloc.local:8091",
  walletURLFrench: "https://wallet.trustbloc.local:8091/fr/",
  demoVerifierURL: "https://demo-adapter.trustbloc.local:8094/verifier",
  demoIssuerURL: "https://demo-adapter.trustbloc.local:8094/issuer",
  webWalletURL: "https://demo-adapter.trustbloc.local:8094/web-wallet",
};
