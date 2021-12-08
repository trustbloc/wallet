// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/trustbloc/edge-agent

go 1.16

require (
	github.com/bluele/gcache v0.0.2 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/duo-labs/webauthn v0.0.0-20200714211715-1daaee874e43
	github.com/duo-labs/webauthn.io v0.0.0-20200929144140-c031a3e0f95d
	github.com/fxamacker/cbor/v2 v2.2.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/sessions v1.2.1
	github.com/hyperledger/aries-framework-go v0.1.8-0.20211203093644-b7d189cc06f4
	github.com/hyperledger/aries-framework-go/component/storageutil v0.0.0-20210916154931-0196c3a2d102
	github.com/hyperledger/aries-framework-go/spi v0.0.0-20210916154931-0196c3a2d102
	github.com/hyperledger/aries-framework-go/test/component v0.0.0-20210820175050-dcc7a225178d // indirect
	github.com/igor-pavlenko/httpsignatures-go v0.0.23
	github.com/piprate/json-gold v0.4.1-0.20210813112359-33b90c4ca86c
	github.com/stretchr/testify v1.7.0
	github.com/trustbloc/edge-core v0.1.7
	github.com/trustbloc/edv v0.1.7
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
)

// Added redirect as a workaround for https://github.com/duo-labs/webauthn/issues/76
replace (
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
	google.golang.org/grpc/examples => google.golang.org/grpc/examples v1.29.1
)
