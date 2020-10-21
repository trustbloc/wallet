// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/trustbloc/edge-agent/cmd/agent-js-worker

go 1.15

require (
	github.com/google/uuid v1.1.2
	github.com/hyperledger/aries-framework-go v0.1.5-0.20201020160650-4535370d64e4
	github.com/hyperledger/aries-framework-go/component/storage/jsindexeddb v0.0.0-20201020181044-b37cbb80b085 // indirect
	github.com/mitchellh/mapstructure v1.3.3
	github.com/stretchr/testify v1.6.1
	github.com/trustbloc/edge-agent v0.0.0
	github.com/trustbloc/edge-core v0.1.5-0.20200916124536-c32454a16108
)

replace github.com/trustbloc/edge-agent => ../../
