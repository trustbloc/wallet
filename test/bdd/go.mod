// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/trustbloc/edge-agent/test/bdd

go 1.16

require (
	github.com/cucumber/godog v0.9.0
	github.com/duo-labs/webauthn v0.0.0-20200714211715-1daaee874e43
	github.com/fsouza/go-dockerclient v1.7.0
	github.com/fxamacker/cbor/v2 v2.2.0
	github.com/google/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.6.7
	github.com/trustbloc/edge-agent v0.0.0
	github.com/trustbloc/edge-core v0.1.7-0.20210527163745-994ae929f957
)

replace (
	github.com/trustbloc/edge-agent => ../..
	// https://github.com/ory/dockertest/issues/208#issuecomment-686820414
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200909081042-eff7692f9009
	// Added redirect as a workaround for https://github.com/duo-labs/webauthn/issues/76
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
	google.golang.org/grpc/examples => google.golang.org/grpc/examples v1.29.1
)
