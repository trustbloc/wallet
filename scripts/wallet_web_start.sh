#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

echo "Running $0"

cd test/fixtures/wallet-web
(
  source .env &&
  docker-compose -f docker-compose-deps.yml -f docker-compose-demo.yml -f docker-compose-server.yml -f docker-compose-web.yml down &&
  docker-compose -f docker-compose-deps.yml -f docker-compose-demo.yml -f docker-compose-server.yml -f docker-compose-web.yml up --force-recreate
)
