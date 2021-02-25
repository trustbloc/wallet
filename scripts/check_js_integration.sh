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
cd $ROOT/test/bdd/fixtures/wallet-web
(source .env && docker-compose -f docker-compose-server.yml -f docker-compose-web.yml down && docker-compose -f docker-compose-server.yml -f docker-compose-web.yml up  --force-recreate -d)

echo "waiting for containers to start..."
sleep 15s

cd $ROOT/cmd/wallet-web
rm -rf node_modules
npm install
npm run build

cd $ROOT/test/wallet-web
rm -rf node_modules
npm install

echo "running tests..."
npm run test

echo "stopping containers..."
cd $ROOT/test/bdd/fixtures/wallet-web

(source .env && docker-compose -f docker-compose-server.yml -f docker-compose-web.yml down --remove-orphans)
