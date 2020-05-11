#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

rm -rf public/aries-framework-go/
mkdir -p public/aries-framework-go/assets
npm install
cp -Rp node_modules/@trustbloc-cicd/aries-framework-go/dist/assets/* public/aries-framework-go/assets
cp public/aries-framework-go/assets/aries-js-worker.wasm.gz public/aries-framework-go/assets/aries-js-worker.wasm.gz.bak
gunzip public/aries-framework-go/assets/aries-js-worker.wasm.gz
mv public/aries-framework-go/assets/aries-js-worker.wasm.gz.bak public/aries-framework-go/assets/aries-js-worker.wasm.gz

rm -rf public/trustbloc-agent/
mkdir -p public/trustbloc-agent/assets
npm install
cp -Rp node_modules/@trustbloc/trustbloc-agent/dist/assets/* public/trustbloc-agent/assets
cp public/trustbloc-agent/assets/trustbloc-agent-js-worker.wasm.gz public/trustbloc-agent/assets/trustbloc-agent-js-worker.wasm.gz.bak
gunzip public/trustbloc-agent/assets/trustbloc-agent-js-worker.wasm.gz
mv public/trustbloc-agent/assets/trustbloc-agent-js-worker.wasm.gz.bak public/trustbloc-agent/assets/trustbloc-agent-js-worker.wasm.gz
