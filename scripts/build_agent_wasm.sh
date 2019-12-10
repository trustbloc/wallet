#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e


echo "Building $1-agent-wasm"
WASM_EXEC=$(go env GOROOT)/misc/wasm/wasm_exec.js
ISSUER_AGENT_WASM_PATH=cmd/$1-agent-wasm
ISSUER_AGENT_WASM_PATH=cmd/$1-agent-wasm
WASM_COMPONENTS=cmd/components

mkdir -p ./build/bin/wasm/$1
cp $WASM_EXEC  ./build/bin/wasm/$1
cp $ISSUER_AGENT_WASM_PATH/*.html  ./build/bin/wasm/$1
cp -r $WASM_COMPONENTS  ./build/bin/wasm/$1
cd $ISSUER_AGENT_WASM_PATH
GOOS=js GOARCH=wasm go build -o ../../build/bin/wasm/$1/$1-agent.wasm main.go
gzip -f ../../build/bin/wasm/$1/$1-agent.wasm
