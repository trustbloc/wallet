#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

echo "Running $0"

ROOT=`pwd`

echo "starting containers..."
cd $ROOT/test/bdd/fixtures/agent-wasm
(source .env && docker-compose down && docker-compose up --force-recreate -d)

echo "waiting for containers to start..."
sleep 15s

cd $ROOT/cmd/user-agent
npm install
npm run serve
