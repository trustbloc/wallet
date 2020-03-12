#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

rm -rf public/aries-framework-go/
mkdir -p public/aries-framework-go/assets
npm install
cp -Rp node_modules/@hyperledger/aries-framework-go/dist/assets/* public/aries-framework-go/assets
gunzip public/aries-framework-go/assets/aries-js-worker.wasm.gz
