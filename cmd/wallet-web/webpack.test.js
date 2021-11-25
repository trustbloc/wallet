/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const { VueLoaderPlugin } = require('vue-loader');
const webpack = require('webpack');

const { alias } = require('./alias.config');

module.exports = {
  mode: 'development',
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: 'vue-loader',
      },
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: 'babel-loader',
      },
      {
        test: /\.css$/,
        use: ['vue-style-loader', 'css-loader'],
      },
      {
        test: /\.(svg|png)/,
        type: 'asset/inline',
      },
    ],
  },
  resolve: {
    alias,
  },
  plugins: [
    new VueLoaderPlugin(),
    new webpack.DefinePlugin({
      process: { env: {} },
    }),
  ],
  devtool: 'inline-cheap-module-source-map',
};
