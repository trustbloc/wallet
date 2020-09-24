#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e


echo "Building trustbloc agent wasm"

cd cmd/trustbloc-agent-js-worker
npm install
npm run build

cd ../agent-js-worker
npm install
npm run build
