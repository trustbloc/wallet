#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

echo "Running $0"

cd test/bdd/fixtures/agent-wasm

(source .env && docker-compose down --remove-orphans && docker-compose -f docker-compose.yml up --force-recreate -d)

bash generate_config.sh

(source .env && docker-compose -f discovery-compose.yml up --force-recreate -d)
