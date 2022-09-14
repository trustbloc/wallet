#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e


if [ "$1" == "setup" ]; then
  echo "generating test keys"
  docker run -i --rm \
  -v $(readlink -f .):/opt/tmp \
  --entrypoint "/opt/tmp/scripts/generate_test_keys.sh" \
  frapsoft/openssl
  echo "setting up agent assets"
  bash ./scripts/setup_assets.sh
  echo "starting containers..."
  cd test/fixtures
  (source .env && docker-compose down --remove-orphans && docker-compose up --force-recreate -d)
  echo "waiting for containers to start..."
  sleep 10
else
   echo "stopping containers..."
   cd test/fixtures
   (source .env && docker-compose down --remove-orphans)
fi

