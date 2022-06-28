# Installation and Building

## Prerequisites

In order to build Wallet from the source, first make sure to install the following:

- [Go 1.17](https://go.dev/doc/install)
- [Npm 8](https://docs.npmjs.com/cli/v8/configuring-npm/install)
- [Node 16](https://nodejs.org/)
- GitHub packages setup: you will need to authenticate to GitHub packages with your [personal token](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line#creating-a-token).
- Configuring `npm` for use with GitHub Packages `echo "//npm.pkg.github.com/:_authToken=${PERSONAL_TOKEN}" > ~/.npmrc`
- [Docker](https://docs.docker.com/get-docker/) (**Note**: make sure to set your `Docker` to use `Docker-Compose V1`)
- [Docker-Compose V1](https://docs.docker.com/compose/install/)
- Make

## Test Keys

In addition, run (`make generate-test-keys`) to generate tls keys and import `test/fixtures/keys/tls/ec-cacert.pem` in your machine's cert chain (e.g. Keychain on Mac).

## Connections

You also need to add the following lines to your `hosts` file (`/etc/hosts` on \*NIX). This file is part of Docker networking mechanism which is primarily used to establish communication between Docker containers and the outside world via the host machine where the Docker daemon is running. See [Docker's documentation](https://docs.docker.com/config/containers/container-networking/) for more details.

```
127.0.0.1 testnet.orb.local
127.0.0.1 wallet.trustbloc.local
127.0.0.1 wallet-server.trustbloc.local
127.0.0.1 mediator.trustbloc.local
127.0.0.1 uni-resolver-web.trustbloc.local
127.0.0.1 auth.trustbloc.local
127.0.0.1 auth-hydra.trustbloc.local
127.0.0.1 demo-hydra.trustbloc.local
127.0.0.1 edv.trustbloc.local
127.0.0.1 kms.trustbloc.local
127.0.0.1 file-server.trustbloc.local
127.0.0.1 demo-adapter.trustbloc.local
```

## Running Wallet Locally

There are 2 ways to run Wallet locally:

1. Static Build
2. Dev Mode (hot deployment for JS & Vue components, useful for developing UI)

- For the Static Build, run the following command:

  ```bash
  make wallet-web-start
  ```

- For the Dev Mode, run the following command instead:
  ```bash
  make wallet-web-dev-start
  ```
