/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const config = require('bedrock').config;
const os = require('os');
const path = require('path');

// common paths
config.paths.cache = path.join(__dirname, '..', '.cache');
config.paths.log = path.join(os.tmpdir(), 'authorization.localhost');

// serve contexts/images/etc
config.express.static.push(path.join(__dirname, '..', 'static'));

// do not require strict SSL in dev mode
config.jsonld.strictSSL = false;

if (process.env.TLS_CERT_FILE) {
    config.server.cert = process.env.TLS_CERT_FILE;
}

if (process.env.TLS_KEY_FILE) {
    config.server.key = process.env.TLS_KEY_FILE;
}

if (process.env.REGISTRATION_URL) {
    config.views.vars.registrationURL = process.env.REGISTRATION_URL
}


