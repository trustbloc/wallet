#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

rm -rf public/agent-js-worker/
mkdir -p public/agent-js-worker/assets
npm install
cp -Rp node_modules/@trustbloc/agent-js-worker/dist/assets/* public/agent-js-worker/assets
cp public/agent-js-worker/assets/agent-js-worker.wasm.gz public/agent-js-worker/assets/agent-js-worker.wasm.gz.bak
gunzip public/agent-js-worker/assets/agent-js-worker.wasm.gz
mv public/agent-js-worker/assets/agent-js-worker.wasm.gz.bak public/agent-js-worker/assets/agent-js-worker.wasm.gz

rm -rf public/trustbloc-agent/
mkdir -p public/trustbloc-agent/assets
npm install
cp -Rp node_modules/@trustbloc/trustbloc-agent/dist/assets/* public/trustbloc-agent/assets
cp public/trustbloc-agent/assets/trustbloc-agent-js-worker.wasm.gz public/trustbloc-agent/assets/trustbloc-agent-js-worker.wasm.gz.bak
gunzip public/trustbloc-agent/assets/trustbloc-agent-js-worker.wasm.gz
mv public/trustbloc-agent/assets/trustbloc-agent-js-worker.wasm.gz.bak public/trustbloc-agent/assets/trustbloc-agent-js-worker.wasm.gz
