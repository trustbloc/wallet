# Copyright SecureKey Technologies Inc.
#
# SPDX-License-Identifier: Apache-2.0


# GO version
ALPINE_VER ?= 3.15
GO_TAGS    ?=
GO_VER     ?= 1.17

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
all: clean checks automation-test

.PHONY: checks
checks: license lint

.PHONY: lint
lint:
	@scripts/check_lint.sh

.PHONY: license
license:
	@scripts/check_license.sh

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

.PHONY: generate-test-keys
generate-test-keys:
	@mkdir -p -p test/fixtures/keys/tls
	@docker run -i --rm \
		-v $(abspath .):/opt/workspace/wallet \
		--entrypoint "/opt/workspace/wallet/scripts/generate_test_keys.sh" \
		frapsoft/openssl

.PHONY: wallet-web-start
wallet-web-start: clean wallet-web-docker mock-images generate-test-keys
	@scripts/wallet_web_start.sh

# starting wallet-web in dev mode for hot deployment
.PHONY: wallet-web-dev-start
wallet-web-dev-start:
	@scripts/wallet_web_dev_start.sh

.PHONY: mock-demo-login-consent-docker
mock-demo-login-consent-docker:
	@echo "Building login consent server for demo..."
	@cd test/mock/demo-login-consent-server && docker build -f image/Dockerfile --build-arg GO_VER=$(GO_VER) --build-arg ALPINE_VER=$(ALPINE_VER) -t edgeagent/demologinconsent:latest .

.PHONY: mock-adapter
mock-adapter:
	@echo "Building mock adapter for demo..."
	@cd test/mock/adapter && docker build -f image/Dockerfile --build-arg GO_VER=$(GO_VER) --build-arg ALPINE_VER=$(ALPINE_VER) -t edgeagent/mockadapter:latest .

.PHONY: mock-images
mock-images: mock-adapter mock-demo-login-consent-docker

.PHONY: automation-test-vcwallet
automation-test-vcwallet: clean wallet-web-docker mock-images generate-test-keys
	@scripts/run_vcwallet_ui_automation.sh

.PHONY: automation-test-wallet
automation-test-wallet: clean wallet-web-docker mock-images generate-test-keys
	@scripts/run_wallet_ui_automation.sh

.PHONY: clean
clean:
	@rm -Rf ./build
	@rm -Rf ./cmd/wallet-web/dist
	@rm -Rf ./cmd/wallet-web/node_modules
	@rm -Rf ./test/fixtures/wallet-web/config
