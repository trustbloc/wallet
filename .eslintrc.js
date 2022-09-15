/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

module.exports = {
  root: true,
  env: { browser: true, es2022: true, node: true, mocha: true },
  extends: [
    'eslint:recommended',
    'plugin:vue/vue3-recommended',
    'plugin:tailwindcss/recommended',
    'plugin:wdio/recommended',
    'plugin:mocha/recommended',
    'plugin:eslint-comments/recommended',
    'plugin:prettier/recommended',
    'prettier',
  ],
  plugins: ['vue', 'tailwindcss', 'wdio', 'mocha', 'eslint-comments', 'prettier'],
  // Default rules for any file we lint
  rules: {
    'vue/multi-word-component-names': [
      'warn',
      {
        ignores: [],
      },
    ],
    /**
     * Force prettier formatting
     */
    'prettier/prettier': 'error',
    /**
     * Disallow the use of console
     * https://eslint.org/docs/rules/no-console
     */
    'no-console': 'warn',

    /**
     * Disallow Reassignment of Function Parameters
     * https://eslint.org/docs/rules/no-param-reassign
     */
    'no-param-reassign': ['error', { props: false }],

    /**
     * Disallow using an async function as a Promise executor
     * https://eslint.org/docs/rules/no-async-promise-executor
     */
    'no-async-promise-executor': 'error',

    /**
     * Disallow await inside of loops
     * https://eslint.org/docs/rules/no-await-in-loop
     */
    // TODO: https://github.com/trustbloc/wallet/issues/1885 fix related code and enable this rule
    'no-await-in-loop': 'warn',

    /**
     * Disallow assignments that can lead to race conditions due to
     * usage of await or yield
     * https://eslint.org/docs/rules/require-atomic-updates
     */
    'require-atomic-updates': 'error',

    /**
     * Disallow async functions which have no await expression
     * https://eslint.org/docs/rules/require-await
     */
    'require-await': 'error',

    /**
     * Require or disallow named function expressions
     * https://eslint.org/docs/rules/func-names
     */
    'func-names': 'off',

    /**
     * Disallow enforcement of consistent linebreak style
     * https://eslint.org/docs/rules/func-names
     */
    'linebreak-style': 'off',

    /**
     * Allow ES6 classes to override methods without using this
     * https://eslint.org/docs/rules/class-methods-use-this
     */
    'class-methods-use-this': 'off',

    'mocha/max-top-level-suites': ['warn', { limit: 2 }],

    /**
     * Added temporary to disable error for defineProps in setup scripts in Vue components
     */
    'no-unused-vars': ['error', { varsIgnorePattern: 'props' }],

    /**
     * TODO: https://github.com/trustbloc/wallet/issues/1886 remove once corresponding issues are fixed
     */
    'mocha/no-setup-in-describe': 'off',
  },
  globals: {
    __webpack_public_path__: 'writable',
  },
  settings: {
    tailwindcss: {
      config: 'cmd/wallet-web/tailwind.js',
    },
  },
};
