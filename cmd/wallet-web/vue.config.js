/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
const path = require('path');
const isSnapshot = require('./package.json').dependencies.hasOwnProperty("@trustbloc-cicd/agent-sdk")
const agent_sdk = (isSnapshot) ? "@trustbloc-cicd/agent-sdk" : "@trustbloc/agent-sdk"

module.exports = {
    chainWebpack: config => config.resolve.symlinks(false),
    runtimeCompiler: true,
    configureWebpack: {
        resolve: {
            alias: {
                "@trustbloc/agent-sdk": path.resolve(__dirname, 'node_modules/' + agent_sdk)
            }
        }
    },
    devServer: {
        https: true
    }
}
