#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

echo "Running $0"

echo "Checking linter..."
(cd $(pwd) && npm install && npm run lint-check)
