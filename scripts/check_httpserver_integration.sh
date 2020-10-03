#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e


echo "Running http-server integration tests..."

ROOT=$(pwd)

cd test/bdd
go test -count=1 -v -cover . -p 1 -timeout=20m -race

cd $ROOT
