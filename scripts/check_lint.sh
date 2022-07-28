#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

echo "Running $0"

echo "Linting wallet-web..."
(cd $(pwd)/cmd/wallet-web && npm install && npm run lint)
