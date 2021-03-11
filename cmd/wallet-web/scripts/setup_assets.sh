#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

rm -rf public/agent-js-worker/
mkdir -p public/agent-js-worker/assets
npm install

if [[ $(grep "@trustbloc-cicd/agent-sdk-web" "package.json") ]] ; then
   cp -Rp node_modules/@trustbloc-cicd/agent-sdk-web/dist/assets/* public/agent-js-worker/assets
else
   cp -Rp node_modules/@trustbloc/agent-sdk-web/dist/assets/* public/agent-js-worker/assets
fi

cp public/agent-js-worker/assets/agent-js-worker.wasm.gz public/agent-js-worker/assets/agent-js-worker.wasm.gz.bak
gunzip public/agent-js-worker/assets/agent-js-worker.wasm.gz
mv public/agent-js-worker/assets/agent-js-worker.wasm.gz.bak public/agent-js-worker/assets/agent-js-worker.wasm.gz
