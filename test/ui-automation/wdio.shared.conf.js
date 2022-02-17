/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

exports.config = {
  runner: "local",
  // Test files
  specs: [
    "./test/specs/waci-flow.js",
    "./test/specs/credential-interaction-flow.js",
    "./test/specs/vault-interaction-flow.js",
  ],
  // Maximum number of total parallel running workers
  maxInstances: 5,
  capabilities: [
    {
      // Maximum number of total parallel running workers per capability
      maxInstances: 5,
      browserName: "chrome",
      "goog:chromeOptions": {
        // to run chrome headless the following flags are required
        // (see https://developers.google.com/web/updates/2017/04/headless-chrome)
        args: [
          "--headless",
          "--no-sandbox",
          "--disable-gpu",
          "--disable-dev-shm-usage",
          "--window-size=1920,1080",
          "--disable-web-security",
          "--ignore-certificate-errors",
        ],
      },
    },
  ],

  // Level of logging verbosity: trace | debug | info | warn | error | silent
  logLevel: "warn",

  // Default timeout for all waitFor* commands.
  waitforTimeout: 60000,

  // Default timeout in milliseconds for request
  // if browser driver or grid doesn't send response
  connectionRetryTimeout: 120000,

  // Default request retries count
  connectionRetryCount: 3,

  // Test runner services
  services: ["chromedriver"],

  // Framework you want to run your specs with.
  framework: "mocha",

  reporters: ["spec"],

  // Options to be passed to Mocha.
  mochaOpts: {
    ui: "bdd",
    timeout: 120000,
    require: ["@babel/register"],
  },
};
