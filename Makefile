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
all: clean checks test-unit test-e2e-vcwallet test-e2e-wallet

.PHONY: checks
checks: check-license check-lint check-prettier

.PHONY: check-license
check-license:
	@scripts/check_license.sh

.PHONY: check-lint
check-lint:
	@scripts/check_lint.sh

.PHONY: check-prettier
check-prettier:
	@set -e
	@cd npm install && npm run prettier-check

.PHONY: build-wallet-web
build-wallet-web:
	@scripts/build_wallet_web.sh

.PHONY: build-wallet-web-docker
build-wallet-web-docker: build-wallet-web
	@echo "Building wallet-web docker image"
	@docker build -f ./images/wallet-web/Dockerfile --no-cache -t $(DOCKER_OUTPUT_NS)/$(REPO_IMAGE_NAME)/wallet-web:latest .

.PHONY: generate-test-keys
generate-test-keys:
	@mkdir -p -p test/fixtures/keys/tls
	@docker run -i --rm \
		-v $(abspath .):/opt/workspace/wallet \
		--entrypoint "/opt/workspace/wallet/scripts/generate_test_keys.sh" \
		frapsoft/openssl

.PHONY: start-wallet-web
start-wallet-web: clean build-wallet-web-docker build-mock-images generate-test-keys
	@scripts/wallet_web_start.sh

# starting wallet-web in dev mode for hot deployment
.PHONY: start-wallet-web-dev
start-wallet-web-dev:
	@scripts/wallet_web_dev_start.sh

.PHONY: build-mock-demo-login-consent-docker
build-mock-demo-login-consent-docker:
	@echo "Building login consent server for demo..."
	@cd test/mock/demo-login-consent-server && docker build -f image/Dockerfile --build-arg GO_VER=$(GO_VER) --build-arg ALPINE_VER=$(ALPINE_VER) -t edgeagent/demologinconsent:latest .

.PHONY: build-mock-adapter
build-mock-adapter:
	@echo "Building mock adapter for demo..."
	@cd test/mock/adapter && docker build -f image/Dockerfile --build-arg GO_VER=$(GO_VER) --build-arg ALPINE_VER=$(ALPINE_VER) -t edgeagent/mockadapter:latest .

.PHONY: build-mock-images
build-mock-images: build-mock-adapter build-mock-demo-login-consent-docker

.PHONY: test-unit
test-unit:
	@set -e
	@cd cmd/wallet-web && npm install && npm run test

.PHONY: test-e2e-vcwallet
test-e2e-vcwallet: clean build-wallet-web-docker build-mock-images generate-test-keys
	@scripts/run_vcwallet_ui_automation.sh

.PHONY: test-e2e-wallet
test-e2e-wallet: clean build-wallet-web-docker build-mock-images generate-test-keys
	@scripts/run_wallet_ui_automation.sh

.PHONY: clean
clean:
	@rm -Rf ./build
	@rm -Rf ./node_modules
	@rm -Rf ./cmd/wallet-web/dist
	@rm -Rf ./cmd/wallet-web/node_modules
	@rm -Rf ./test/fixtures/wallet-web/config
	@rm -Rf ./test/fixtures/ui-automation/node_modules
