#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

echo "Running $0"

DOCKER_CMD=${DOCKER_CMD:-docker}

if [ ! $(command -v ${DOCKER_CMD}) ]; then
    exit 0
fi

echo "Linting user-agent..."
(cd $(pwd)/cmd/user-agent && npm install && npm run lint)
echo "Linting http-server..."
${DOCKER_CMD} run --rm -v $(pwd):/opt/workspace -w /opt/workspace/cmd/http-server golangci/golangci-lint:v1.31.0 golangci-lint run -c ../../.golangci.yml
echo "Linting BDD test packages..."
${DOCKER_CMD} run --rm -v $(pwd):/opt/workspace -w /opt/workspace/test/bdd golangci/golangci-lint:v1.31.0 golangci-lint run -c ../../.golangci.yml
