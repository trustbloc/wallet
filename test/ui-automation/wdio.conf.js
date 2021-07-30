/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

'use strict';

const {config} = require('./wdio.shared.conf');

const domain = ".dev.trustbloc.dev"

exports.config = {
    ...config,
    "walletName": "wallet.trustbloc.local",
    "walletURL": "https://wallet.trustbloc.local:8091",
};
