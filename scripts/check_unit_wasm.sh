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
PKGS="github.com/trustbloc/edge-agent/pkg/didexchange/invitation github.com/trustbloc/edge-agent/pkg/store/vc"
PATH="$GOBIN:$PATH" GOOS=js GOARCH=wasm go test $PKGS -count=1 -exec=wasmbrowsertest -cover -timeout=10m
