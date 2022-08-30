#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

echo "Running $0"

# Use corresponding image based on architecture type (arm64/amd64)
arch=$(uname -p)
if  [[ $arch == arm64 ]]; then
    export MYSQL_IMAGE=arm64v8/mysql
fi

cd test/fixtures/wallet-web
(source .env && docker-compose -f docker-compose-demo.yml -f docker-compose-server.yml -f docker-compose-web.yml down && docker-compose -f docker-compose-demo.yml -f docker-compose-server.yml -f docker-compose-web.yml up --force-recreate)
