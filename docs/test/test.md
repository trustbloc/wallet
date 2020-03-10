# Edge Agent - Test

## Prerequisites (General)

- Vue.js
- Go 1.13
- Npm

## Prerequisites (for running tests and demos)
- Docker
- Docker-Compose
- Make

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

- To access user agent wasm open [user home page](https://127.0.0.1:8091/index.html). 
