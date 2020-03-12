#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e


echo "Building $1-agent-wasm"
AGENT_WEB_PATH=cmd/$1-agent

mkdir -p ./build/bin/wasm/$1

cd $AGENT_WEB_PATH
npm install
npm run build

cp -R dist/* ../../build/bin/wasm/$1
