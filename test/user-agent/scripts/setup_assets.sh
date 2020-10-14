#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

rm -rf public
mkdir -p public/agent-js-worker/assets
cp -Rp node_modules/@trustbloc/agent-js-worker/dist/assets/* public/agent-js-worker/assets
gunzip public/agent-js-worker/assets/agent-js-worker.wasm.gz
