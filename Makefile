# Copyright SecureKey Technologies Inc.
#
# SPDX-License-Identifier: Apache-2.0


HTTP_SERVER_PATH       = cmd/http-server
GOBIN_PATH             = $(abspath .)/build/bin

# GO version
ALPINE_VER ?= 3.11
GO_TAGS    ?=
GO_VER     ?= 1.13.9
NODE_VER     ?= 14.8.0

# Namespace for the images
DOCKER_OUTPUT_NS         ?= docker.pkg.github.com
REPO_IMAGE_NAME          ?= trustbloc/edge-agent


.PHONY: all
all: clean checks unit-test

.PHONY: depend
depend:
	@mkdir -p ./build/bin
	GO111MODULE=off GOBIN=$(GOBIN_PATH) go get github.com/agnivade/wasmbrowsertest

.PHONY: checks
checks: license lint

.PHONY: lint
lint:
	@scripts/check_lint.sh

.PHONY: license
license:
	@scripts/check_license.sh

.PHONY: unit-test
unit-test:
	@scripts/check_unit.sh

.PHONY: unit-test-wasm
unit-test-wasm: export GOBIN=$(GOBIN_PATH)
unit-test-wasm: depend
	@scripts/check_unit_wasm.sh

.PHONY: agent-wasm
agent-wasm:
	@scripts/build_agent_wasm.sh ${AGENT_NAME}

.PHONY: trustbloc-agent-wasm
trustbloc-agent-wasm:
	@scripts/build_trustbloc_agent_wasm.sh

.PHONY: user-agent-wasm
user-agent-wasm:
	AGENT_NAME="user" make agent-wasm

.PHONY: http-server
http-server:
	@echo "Building http-server"
	@mkdir -p ./build/bin/wasm
	@cd ${HTTP_SERVER_PATH} && go build -o ../../build/bin/http-server main.go

.PHONY: user-agent-wasm-docker
user-agent-wasm-docker:
	AGENT_NAME="user" make agent-wasm-docker

.PHONY: agent-wasm-docker
agent-wasm-docker: clean
	@echo "Building agent wasm docker image"
	@docker build -f ./images/agent-wasm/Dockerfile --no-cache -t $(DOCKER_OUTPUT_NS)/$(REPO_IMAGE_NAME)/${AGENT_NAME}-agent-wasm:latest \
	--build-arg GO_VER=$(GO_VER) \
	--build-arg ALPINE_VER=$(ALPINE_VER) \
	--build-arg GO_TAGS=$(GO_TAGS) \
	--build-arg GITHUB_TOKEN=$(GITHUB_TOKEN) \
	--build-arg NAME=${AGENT_NAME} .

.PHONY: credential-mediator-docker
credential-mediator-docker: clean
	@echo "Building credential mediator docker image"
	@docker build -f ./images/credential-mediator/Dockerfile --no-cache -t $(DOCKER_OUTPUT_NS)/$(REPO_IMAGE_NAME)/credential-mediator:latest \
	--build-arg NODE_VER=$(NODE_VER) \
	--build-arg ALPINE_VER=$(ALPINE_VER) \
	--build-arg GITHUB_TOKEN=$(GITHUB_TOKEN) .

.PHONY: generate-test-keys
generate-test-keys:
	@mkdir -p -p test/bdd/fixtures/keys/tls
	@docker run -i --rm \
		-v $(abspath .):/opt/workspace/edge-agent \
		--entrypoint "/opt/workspace/edge-agent/scripts/generate_test_keys.sh" \
		frapsoft/openssl

.PHONY: user-agent-start
user-agent-start: clean user-agent-wasm-docker credential-mediator-docker generate-test-config generate-test-keys
	@scripts/user_agent_start.sh

.PHONY: generate-test-config
generate-test-config:
	@/bin/bash scripts/generate_test_config.sh

.PHONY: bdd-test-js
bdd-test-js: clean user-agent-wasm-docker credential-mediator-docker generate-test-config generate-test-keys
	@scripts/check_js_intergation.sh

.PHONY: clean
clean: clean-build

.PHONY: clean-build
clean-build:
	@rm -Rf ./build
	@rm -Rf ./cmd/user-agent/dist
	@rm -Rf ./cmd/user-agent/node_modules
	@rm -Rf ./cmd/trustbloc-agent-js-worker/node_modules
	@rm -Rf ./cmd/trustbloc-agent-js-worker/dist
	@rm -Rf ./test/bdd/fixtures/agent-wasm/config
