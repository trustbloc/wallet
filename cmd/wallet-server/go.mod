// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/trustbloc/edge-agent/cmd/wallet-server

go 1.15

require (
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/duo-labs/webauthn v0.0.0-20200714211715-1daaee874e43
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.0
	github.com/hyperledger/aries-framework-go v0.1.5
	github.com/hyperledger/aries-framework-go-ext/component/storage/mysql v0.0.0-20210105175634-a0737dbb3570
	github.com/pquerna/cachecontrol v0.0.0-20200819021114-67c6ae64274f // indirect
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v0.0.6
	github.com/stretchr/testify v1.6.1
	github.com/trustbloc/edge-agent v0.0.0-00010101000000-000000000000
	github.com/trustbloc/edge-core v0.1.5
)

replace github.com/trustbloc/edge-agent => ../..

// Added redirect as a workaround for https://github.com/duo-labs/webauthn/issues/76
replace (
	github.com/kilic/bls12-381 => github.com/trustbloc/bls12-381 v0.0.0-20201104214312-31de2a204df8
	github.com/trustbloc/edge-core => github.com/trustbloc/edge-core v0.1.5-0.20201126210935-53388acb41fc
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
	google.golang.org/grpc/examples => google.golang.org/grpc/examples v1.29.1

)
