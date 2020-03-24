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

${DOCKER_CMD} run --rm -e GITHUB_TOKEN=${GITHUB_TOKEN} -v $(pwd):/opt/workspace -w /opt/workspace/cmd/user-agent/ node:12.13.1-alpine3.10 /bin/sh -c "npm install;npm run lint"
${DOCKER_CMD} run --rm -v $(pwd):/opt/workspace -w /opt/workspace/cmd/http-server golangci/golangci-lint:v1.21 golangci-lint run -c ../../.golangci.yml
