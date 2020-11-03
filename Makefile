# Copyright SecureKey Technologies Inc.
#
# SPDX-License-Identifier: Apache-2.0


HTTP_SERVER_PATH       = cmd/http-server

# GO version
ALPINE_VER ?= 3.12
GO_TAGS    ?=
GO_VER     ?= 1.15.2

# Namespace for the images
DOCKER_OUTPUT_NS         ?= docker.pkg.github.com
REPO_IMAGE_NAME          ?= trustbloc/edge-agent


.PHONY: all
all: clean checks unit-test

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

.PHONY: agent-wasm
agent-wasm:
	@scripts/build_agent_wasm.sh ${AGENT_NAME}

.PHONY: user-agent-wasm
user-agent-wasm:
	AGENT_NAME="user" make agent-wasm

.PHONY: http-server
http-server:
	@echo "Building http-server"
	@mkdir -p ./build/bin/wasm
	@cd ${HTTP_SERVER_PATH} && go build -o ../../build/bin/http-server main.go

.PHONY: user-agent-wasm-docker
user-agent-wasm-docker: clean user-agent-wasm
	AGENT_NAME="user" make agent-wasm-docker

.PHONY: agent-wasm-docker
agent-wasm-docker:
	@echo "Building agent wasm docker image"
	@docker build -f ./images/agent-wasm/Dockerfile --no-cache -t $(DOCKER_OUTPUT_NS)/$(REPO_IMAGE_NAME)/${AGENT_NAME}-agent-wasm:latest \
	--build-arg GO_VER=$(GO_VER) \
	--build-arg ALPINE_VER=$(ALPINE_VER) \
	--build-arg GO_TAGS=$(GO_TAGS) \
	--build-arg NAME=${AGENT_NAME} .

.PHONY: generate-test-keys
generate-test-keys:
	@mkdir -p -p test/bdd/fixtures/keys/tls
	@docker run -i --rm \
		-v $(abspath .):/opt/workspace/edge-agent \
		--entrypoint "/opt/workspace/edge-agent/scripts/generate_test_keys.sh" \
		frapsoft/openssl

.PHONY: user-agent-start
user-agent-start: clean user-agent-wasm-docker generate-test-config generate-test-keys mock-login-consent-docker
	@scripts/user_agent_start.sh

# starting user agent in dev mode for hot deployment
.PHONY: user-agent-dev-start
user-agent-dev-start:
	@scripts/user_agent_dev_start.sh

.PHONY: generate-test-config
generate-test-config:
	@/bin/bash scripts/generate_test_config.sh

.PHONY: bdd-test
bdd-test: bdd-test-js bdd-test-http-server

.PHONY: bdd-test-js
bdd-test-js: clean user-agent-wasm-docker generate-test-config generate-test-keys mock-login-consent-docker
	@scripts/check_js_integration.sh

.PHONY: bdd-test-http-server
bdd-test-http-server: clean user-agent-wasm-docker generate-test-config generate-test-keys mock-login-consent-docker
	@scripts/check_httpserver_integration.sh

.PHONY: mock-login-consent-docker
mock-login-consent-docker:
	@echo "Building mock login consent server for BDD tests..."
	@cd test/bdd/mock/loginconsent && docker build -f image/Dockerfile --build-arg GO_VER=$(GO_VER) --build-arg ALPINE_VER=$(ALPINE_VER) -t edgeagent/mockloginconsent:latest .

.PHONY: clean
clean: clean-build

.PHONY: clean-build
clean-build:
	@rm -Rf ./build
	@rm -Rf ./cmd/user-agent/dist
	@rm -Rf ./cmd/user-agent/node_modules
	@rm -Rf ./test/bdd/fixtures/agent-wasm/config
