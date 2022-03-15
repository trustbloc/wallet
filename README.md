[![Release](https://img.shields.io/github/release/trustbloc/wallet.svg?style=flat-square)](https://github.com/trustbloc/wallet/releases/latest)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://raw.githubusercontent.com/trustbloc/wallet/main/LICENSE)
[![Godocs](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/trustbloc/wallet)

[![Build Status](https://github.com/trustbloc/wallet/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/trustbloc/wallet/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/trustbloc/wallet/branch/main/graph/badge.svg)](https://codecov.io/gh/trustbloc/wallet)
[![Go Report Card](https://goreportcard.com/badge/github.com/trustbloc/wallet)](https://goreportcard.com/report/github.com/trustbloc/wallet)

# Wallet

The TrustBloc Wallet repo contains a Web Wallet to store/share [W3C Verifiable Credential(VC)](https://w3c.github.io/vc-data-model/) signed/verified with [W3C Decentralized Identifier(DID)](https://w3c.github.io/did-core/). These wallets are useful for the [holder](https://www.w3.org/TR/vc-data-model/#dfn-holders) role defined in [W3C Verifiable Credential](https://www.w3.org/TR/vc-data-model/#ecosystem-overview) Specification.

The wallet uses [TrustBloc Wallet SDK](https://github.com/trustbloc/agent-sdk/tree/main/cmd/wallet-js-sdk) built on top of [Aries Verifiable Credential wallet](https://github.com/hyperledger/aries-framework-go/blob/main/docs/vc_wallet.md) based on [Universal Wallet Specifications](https://w3c-ccg.github.io/universal-wallet-interop-spec/) implementation.

The [Wallet Web](./cmd/wallet-web) works together with [Wallet Server](./cmd/wallet-server) to support user actions.

## Components

- [Wallet Web](docs/components/wallet_web.md)
- [Wallet Server](docs/components/wallet_server.md)

## Build

To build from source see our [build documentation](docs/build.md).

## Test

- [Test](docs/test/test.md)

## Contributing

Thank you for your interest in contributing. Please see our [community contribution guidelines](https://github.com/trustbloc/community/blob/main/CONTRIBUTING.md) for more information.

## License

Apache License, Version 2.0 (Apache-2.0). See the [LICENSE](LICENSE) file.
