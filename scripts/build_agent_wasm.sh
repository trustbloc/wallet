#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

echo "Building user-agent"

mkdir -p ./build/bin/user-agent

cd cmd/user-agent
npm install
npm run build

cp -R dist/* ../../build/bin/user-agent
