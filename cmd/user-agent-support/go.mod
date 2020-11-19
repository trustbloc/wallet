// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/trustbloc/edge-agent/cmd/user-agent-support

go 1.15

require (
	github.com/cenkalti/backoff/v4 v4.0.2
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/duo-labs/webauthn v0.0.0-20200714211715-1daaee874e43
	github.com/golang/gddo v0.0.0-20200831202555-721e228c7686 // indirect
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.0
	github.com/lpar/gzipped v1.1.0
	github.com/pquerna/cachecontrol v0.0.0-20200819021114-67c6ae64274f // indirect
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v0.0.6
	github.com/stretchr/testify v1.6.1
	github.com/trustbloc/edge-agent v0.0.0-00010101000000-000000000000
	github.com/trustbloc/edge-core v0.1.5-0.20201106164919-76ecfeca954f
)

replace github.com/trustbloc/edge-agent => ../..

// Added redirect as a workaround for https://github.com/duo-labs/webauthn/issues/76
replace (
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
	google.golang.org/grpc/examples => google.golang.org/grpc/examples v1.29.1
)
