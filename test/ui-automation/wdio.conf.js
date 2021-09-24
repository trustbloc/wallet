/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

import { config } from "./wdio.shared.conf";

exports.config = {
  ...config,
  walletName: "TrustBloc Wallet",
  walletURL: "https://wallet.trustbloc.local:8091",
};
