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


${DOCKER_CMD} run --rm -v $(pwd):/opt/workspace -w /opt/workspace/cmd/issuer-agent/web node:12.13.1-alpine3.10 /bin/sh -c "npm install;npm run lint"
${DOCKER_CMD} run --rm -v $(pwd):/opt/workspace -w /opt/workspace golangci/golangci-lint:v1.21 golangci-lint run
${DOCKER_CMD} run --rm -e GOOS=js -e GOARCH=wasm -v $(pwd):/opt/workspace -w /opt/workspace golangci/golangci-lint:v1.21 golangci-lint run
${DOCKER_CMD} run --rm -e GOOS=js -e GOARCH=wasm -v $(pwd):/opt/workspace -w /opt/workspace/cmd/issuer-agent/wasm golangci/golangci-lint:v1.21 golangci-lint run -c ../../../.golangci.yml
${DOCKER_CMD} run --rm -v $(pwd):/opt/workspace -w /opt/workspace/cmd/http-server golangci/golangci-lint:v1.21 golangci-lint run -c ../../.golangci.yml
