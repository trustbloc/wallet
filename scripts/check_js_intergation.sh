#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

echo "Running $0"

ROOT=`pwd`

#TODO container start/stop to be enabled once tests depending on containers will be added
#echo "starting containers..."
#cd $ROOT/test/bdd/fixtures/agent-wasm
#(source .env && docker-compose down && docker-compose up --force-recreate -d)

cd $ROOT/test/user-agent
npm install
echo "running tests..."
npm run test

#echo "stopping containers..."
#cd $ROOT/test/bdd/fixtures/agent-wasm
#(source .env && docker-compose down --remove-orphans)