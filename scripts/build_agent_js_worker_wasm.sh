#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e


echo "Building agent js worker wasm"

cd cmd/agent-js-worker
npm install
npm run build
