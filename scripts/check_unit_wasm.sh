
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

# run unit test for agent-js-worker
cd cmd/agent-js-worker
PATH="$GOBIN:$PATH" GOOS=js GOARCH=wasm go test "github.com/trustbloc/edge-agent/cmd/agent-js-worker" -count=1 -exec=wasmbrowsertest -cover -timeout=10m
