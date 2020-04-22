# Edge Agent - Test

## Prerequisites (General)

- Vue.js
- Go 1.13
- Npm
- GitHub packages setup : generate github [personal token](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line#creating-a-token) 
and set it to GITHUB_TOKEN [environment variable](https://en.wikipedia.org/wiki/Environment_variable).

## Prerequisites (for running tests and demos)
- Docker
- Docker-Compose
- Make

## BDD Test Prerequisites

- Run (`make generate-test-keys`) to generate tls keys and import ec-cacert.pem in cert chain

- You need to modify your hosts file (`/etc/hosts` on \*NIX) to add the following lines, to allow few of the bdd test containers to be connected to externally. 

    127.0.0.1 testnet.trustbloc.local
    127.0.0.1 stakeholder.one
    127.0.0.1 sidetree-mock

## Targets
```
# run checks and unit tests
make all

# run linter checks
make checks

# run unit test
make unit-test

# run unit test for wasm components
make unit-test-wasm

# create docker image for user agent wasm
make user-agent-wasm-docker

# generate tls keys
make generate-test-keys
```

## Steps

```bash
make clean user-agent-wasm-docker generate-test-keys
cd test/bdd/fixtures/agent-wasm
(source .env && docker-compose down && docker-compose up --force-recreate)
```

## Agents

- To access user agent wasm open [user home page](https://127.0.0.1:8091/dashboard).
