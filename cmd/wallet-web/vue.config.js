/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const { alias } = require('./alias.config');
const path = require('path');

module.exports = {
  chainWebpack: (config) => {
    config.resolve.symlinks(false);
  },
  runtimeCompiler: true,

  configureWebpack: {
    resolve: {
      alias,
    },
  },

  devServer: {
    https: true,
  },

  pluginOptions: {
    i18n: {
      locale: process.env.VUE_APP_I18N_FALLBACK_LOCALE || 'en',
      fallbackLocale: process.env.VUE_APP_I18N_FALLBACK_LOCALE || 'en',
      localeDir: 'translations',
      enableInSFC: true,
    },
  },
};
