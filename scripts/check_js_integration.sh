#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

echo "Running $0"

ROOT=`pwd`

npm -v
echo "starting containers..."
cd $ROOT/test/bdd/fixtures/agent-wasm
(source .env && docker-compose down && docker-compose up --force-recreate -d)

echo "waiting for containers to start..."
sleep 15s

cd $ROOT/cmd/user-agent
npm install
npm run build

cd $ROOT/test/user-agent
rm -rf package-lock.json
rm -rf node_modules
npm install

echo "running tests..."
npm run test

echo "stopping containers..."
cd $ROOT/test/bdd/fixtures/agent-wasm
(source .env && docker-compose down --remove-orphans)
