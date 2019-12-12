#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e


echo "Building $1-agent-wasm"
pwd=`pwd`
WASM_EXEC=$(go env GOROOT)/misc/wasm/wasm_exec.js
AGENT_WASM_PATH=cmd/$1-agent/wasm
AGENT_WEB_PATH=cmd/$1-agent/web


mkdir -p ./build/bin/wasm/$1
cp $WASM_EXEC  ./build/bin/wasm/$1
cd $AGENT_WASM_PATH
GOOS=js GOARCH=wasm go build -o ../../../build/bin/wasm/$1/$1-agent.wasm main.go
gzip -f ../../../build/bin/wasm/$1/$1-agent.wasm

cd $pwd
cd $AGENT_WEB_PATH
npm install
npm run build
cp -R dist/* ../../../build/bin/wasm/$1
