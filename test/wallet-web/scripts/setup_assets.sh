#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

rm -rf public
mkdir -p public/agent-js-worker/assets

if [[ $(grep "@trustbloc-cicd/agent-sdk-web" "package.json") ]] ; then
  cp -Rp node_modules/@trustbloc-cicd/agent-sdk-web/dist/assets/* public/agent-js-worker/assets
else
  cp -Rp node_modules/@trustbloc/agent-sdk-web/dist/assets/* public/agent-js-worker/assets
fi

gunzip public/agent-js-worker/assets/agent-js-worker.wasm.gz
