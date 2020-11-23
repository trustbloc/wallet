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
cd $ROOT/test/bdd/fixtures/wallet-web
(source .env && docker-compose down && docker-compose up --force-recreate -d)

echo "waiting for containers to start..."
sleep 15s

cd $ROOT/cmd/wallet-web
npm install
npm run serve
