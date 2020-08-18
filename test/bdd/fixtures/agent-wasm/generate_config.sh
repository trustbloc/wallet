#!/usr/bin/env bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

echo `pwd`

../../../../build/did-method-cli/cli create-config --sidetree-url https://localhost:48326/sidetree/0.0.1 --tls-cacerts ../keys/tls/ec-cacert.pem \
--sidetree-write-token rw_token --config-file ./config-data/config.json --output-directory ./config 2>&1
rm -rf ./config/stakeholder.one
mv ./config/stakeholder.one:8088 ./config/stakeholder.one
