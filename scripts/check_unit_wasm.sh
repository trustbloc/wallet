#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

echo "Running $0"

# Running wasm unit test
# TODO Support collecting code coverage  https://github.com/agnivade/wasmbrowsertest/issues/5
# TODO remove wasm related files
PKGS=""
PATH="$GOBIN:$PATH" GOOS=js GOARCH=wasm go test $PKGS -count=1 -exec=wasmbrowsertest -cover -timeout=10m
