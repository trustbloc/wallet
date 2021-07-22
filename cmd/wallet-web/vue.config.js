/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
const path = require('path');
const isSnapshotAgent = require('./package.json').dependencies.hasOwnProperty("@trustbloc-cicd/agent-sdk-web")
const isSnapshotSDK = require('./package.json').dependencies.hasOwnProperty("@trustbloc-cicd/wallet-sdk")
const agent_sdk = (isSnapshotAgent) ? "@trustbloc-cicd/agent-sdk-web" : "@trustbloc/agent-sdk-web"
const wallet_sdk = (isSnapshotSDK) ? "@trustbloc-cicd/wallet-sdk" : "@trustbloc/wallet-sdk"

module.exports = {
    chainWebpack: config => config.resolve.symlinks(false),
    runtimeCompiler: true,

    configureWebpack: {
        resolve: {
            alias: {
                "@trustbloc/agent-sdk-web": path.resolve(__dirname, 'node_modules/' + agent_sdk),
                "@trustbloc/wallet-sdk": path.resolve(__dirname, 'node_modules/' + wallet_sdk)
            }
        }
    },

    devServer: {
        https: true
    },

    pluginOptions: {
      i18n: {
        locale: 'en',
        fallbackLocale: 'en',
        localeDir: 'locales',
        enableInSFC: true
      }
    }
}
