#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

echo "Running $0"

cd test/bdd/fixtures/wallet-web
(source .env && docker-compose down && docker-compose up --force-recreate)