/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const path = require('path');
const { dependencies } = require('./package.json');

const vueSrc = './src';

const isSnapshotAgent = Object.prototype.hasOwnProperty.call(
  dependencies,
  '@trustbloc-cicd/agent-sdk-web'
);
const isSnapshotSDK = Object.prototype.hasOwnProperty.call(
  dependencies,
  '@trustbloc-cicd/wallet-sdk'
);

const AGENT_SDK = isSnapshotAgent ? '@trustbloc-cicd/agent-sdk-web' : '@trustbloc/agent-sdk-web';
const WALLET_SDK = isSnapshotSDK ? '@trustbloc-cicd/wallet-sdk' : '@trustbloc/wallet-sdk';

module.exports = {
  alias: {
    '@trustbloc/agent-sdk-web': path.resolve(__dirname, `node_modules/${AGENT_SDK}`),
    '@trustbloc/wallet-sdk': path.resolve(__dirname, `node_modules/${WALLET_SDK}`),
    '@': path.resolve(__dirname, vueSrc),
  },
};
