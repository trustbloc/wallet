/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

const chapi = require("./chapi");
const wallet = require("./wallet");

const api = {};
module.exports = api;

api.chapi = chapi;
api.wallet = wallet;
