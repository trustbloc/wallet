/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

module.exports = {
    publicPath: '/wallet',
    chainWebpack: config => config.resolve.symlinks(false),
    runtimeCompiler: true,
    devServer: {
        https: true
    }
}
