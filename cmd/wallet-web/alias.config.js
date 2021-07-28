/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
const path = require('path');

let isSnapshotAgent = require('./package.json').dependencies.hasOwnProperty(
    '@trustbloc-cicd/agent-sdk-web'
);
let isSnapshotSDK = require('./package.json').dependencies.hasOwnProperty(
    '@trustbloc-cicd/wallet-sdk'
);

let AGENT_SDK = isSnapshotAgent ? '@trustbloc-cicd/agent-sdk-web' : '@trustbloc/agent-sdk-web';
let WALLET_SDK = isSnapshotSDK ? '@trustbloc-cicd/wallet-sdk' : '@trustbloc/wallet-sdk';

module.exports = {
    alias: {
        "@trustbloc/agent-sdk-web": path.resolve(__dirname, `node_modules/${AGENT_SDK}`),
        '@trustbloc/wallet-sdk': path.resolve(__dirname, `node_modules/${WALLET_SDK}`)
    }
}

