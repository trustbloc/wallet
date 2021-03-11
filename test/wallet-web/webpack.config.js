/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const VueLoaderPlugin = require('vue-loader/lib/plugin')
const path = require('path');
const isSnapshot = require('./package.json').dependencies.hasOwnProperty("@trustbloc-cicd/agent-sdk-web")
const agent_sdk = (isSnapshot) ? "@trustbloc-cicd/agent-sdk-web" : "@trustbloc/agent-sdk-web"

module.exports = {
    mode: 'development',
    module: {
        rules: [
            {
                test: /\.vue$/,
                loader: 'vue-loader'
            },
            {
                test: /\.js$/,
                loader: 'babel-loader'
            },
            {
                test: /\.css$/,
                use: [
                    'vue-style-loader',
                    'css-loader'
                ]
            }
        ]
    },
    resolve: {
        alias: {
            "@trustbloc/agent-sdk-web": path.resolve(__dirname, 'node_modules/' + agent_sdk)
        }
    },
    plugins: [
        new VueLoaderPlugin()
    ],
    devtool: 'inline-cheap-module-source-map'
}

