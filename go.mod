// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/trustbloc/edge-agent

go 1.16

require (
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/duo-labs/webauthn v0.0.0-20200714211715-1daaee874e43
	github.com/duo-labs/webauthn.io v0.0.0-20200929144140-c031a3e0f95d
	github.com/fxamacker/cbor/v2 v2.2.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/sessions v1.2.1
	github.com/hyperledger/aries-framework-go v0.1.7-0.20210517160459-a72f856f36b8
	github.com/hyperledger/aries-framework-go/component/storageutil v0.0.0-20210427144858-06fb8b7d2d30
	github.com/hyperledger/aries-framework-go/spi v0.0.0-20210422133815-2ef2d99cb692
	github.com/igor-pavlenko/httpsignatures-go v0.0.23
	github.com/piprate/json-gold v0.4.0
	github.com/stretchr/testify v1.7.0
	github.com/trustbloc/edge-core v0.1.7-0.20210429084532-c385778b4d8b
	github.com/trustbloc/edv v0.1.7-0.20210429171638-8c7322f78e77
	github.com/trustbloc/kms v0.1.7-0.20210429121724-8debc3627519
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
)

// Added redirect as a workaround for https://github.com/duo-labs/webauthn/issues/76
replace (
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
	google.golang.org/grpc/examples => google.golang.org/grpc/examples v1.29.1
)
