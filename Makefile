# Copyright SecureKey Technologies Inc.
#
# SPDX-License-Identifier: Apache-2.0


WALLET_SERVER_PATH	= cmd/wallet-server

# GO version
ALPINE_VER ?= 3.12
GO_TAGS    ?=
GO_VER     ?= 1.16

# open API configuration
OPENAPI_SPEC_PATH=build/rest/openapi/spec
OPENAPI_DOCKER_IMG=quay.io/goswagger/swagger
OPENAPI_DOCKER_IMG_VERSION=v0.23.0

# Namespace for the images
DOCKER_OUTPUT_NS         ?= ghcr.io
REPO_IMAGE_NAME          ?= trustbloc

# shell/term
export SHELL	:= /bin/bash
export TERM		:= xterm-256color

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

.PHONY: wallet-web-test
wallet-web-test:
	@set -e
	@cd cmd/wallet-web && npm install && npm run test

.PHONY: wallet-web-prettier-check
wallet-web-prettier-check:
	@set -e
	@cd cmd/wallet-web && npm install && npm run prettier-check

.PHONY: wallet-web-eslint-check
wallet-web-eslint-check:
	@set -e
	@cd cmd/wallet-web && npm install && npm run lint

.PHONY: wallet-web
wallet-web:
	@scripts/build_wallet_web.sh

.PHONY: wallet-web-docker
wallet-web-docker: wallet-web
	@echo "Building wallet-web docker image"
	@docker build -f ./images/wallet-web/Dockerfile --no-cache -t $(DOCKER_OUTPUT_NS)/$(REPO_IMAGE_NAME)/wallet-web:latest .

.PHONY: wallet-server
wallet-server:
	@echo "Building wallet-server"
	@cd ${WALLET_SERVER_PATH} && go build -o ../../build/bin/wallet-server main.go

.PHONY: wallet-server-docker
wallet-server-docker:
	@echo "Building wallet-server docker image"
	@docker build -f ./images/wallet-server/Dockerfile --no-cache -t $(DOCKER_OUTPUT_NS)/$(REPO_IMAGE_NAME)/wallet-server:latest \
	--build-arg GO_VER=$(GO_VER) \
	--build-arg ALPINE_VER=$(ALPINE_VER) \
	--build-arg GO_TAGS=$(GO_TAGS) .

.PHONY: generate-test-keys
generate-test-keys:
	@mkdir -p -p test/fixtures/keys/tls
	@docker run -i --rm \
		-v $(abspath .):/opt/workspace/edge-agent \
		--entrypoint "/opt/workspace/edge-agent/scripts/generate_test_keys.sh" \
		frapsoft/openssl

.PHONY: generate-openapi-spec
generate-openapi-spec:
	@echo "Generating and validating controller API specifications using Open API"
	@mkdir -p build/rest/openapi/spec
	@SPEC_LOC=${OPENAPI_SPEC_PATH}  \
	DOCKER_IMAGE=$(OPENAPI_DOCKER_IMG) DOCKER_IMAGE_VERSION=$(OPENAPI_DOCKER_IMG_VERSION)  \
	scripts/generate-openapi-spec.sh

.PHONY: generate-openapi-demo-specs
generate-openapi-demo-specs: generate-openapi-spec
	@echo "Generate demo wallet server rest controller API specifications using Open API"
	@SPEC_PATH=${OPENAPI_SPEC_PATH} OPENAPI_DEMO_PATH=test/fixtures/wallet-web \
    	DOCKER_IMAGE=$(OPENAPI_DOCKER_IMG) DOCKER_IMAGE_VERSION=$(OPENAPI_DOCKER_IMG_VERSION)  \
    	scripts/generate-openapi-demo-specs.sh

.PHONY: run-openapi-demo
run-openapi-demo: generate-openapi-demo-specs wallet-server-docker generate-test-keys mock-images
	@echo "Starting demo wallet server rest containers ..."
	@DEMO_COMPOSE_PATH=test/fixtures/wallet-web scripts/run-openapi-demo.sh

.PHONY: wallet-web-start
wallet-web-start: clean wallet-web-docker wallet-server-docker generate-test-keys mock-images
	@scripts/wallet_web_start.sh

# starting wallet-web in dev mode for hot deployment
.PHONY: wallet-web-dev-start
wallet-web-dev-start:
	@scripts/wallet_web_dev_start.sh

.PHONY: mock-demo-login-consent-docker
mock-demo-login-consent-docker:
	@echo "Building login consent server for demo..."
	@cd test/mock/demo-login-consent-server && docker build -f image/Dockerfile --build-arg GO_VER=$(GO_VER) --build-arg ALPINE_VER=$(ALPINE_VER) -t edgeagent/demologinconsent:latest .

.PHONY: mock-images
mock-images: mock-demo-login-consent-docker

.PHONY: automation-test
automation-test: clean wallet-web-docker wallet-server-docker generate-test-keys mock-images
	@scripts/run_ui_automation.sh

.PHONY: clean
clean:
	@rm -Rf ./build
	@rm -Rf ./cmd/wallet-web/dist
	@rm -Rf ./cmd/wallet-web/node_modules
	@rm -Rf ./test/fixtures/wallet-web/config
