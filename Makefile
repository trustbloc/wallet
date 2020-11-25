# Copyright SecureKey Technologies Inc.
#
# SPDX-License-Identifier: Apache-2.0


WALLET_SERVER_PATH	= cmd/wallet-server

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
	@mkdir -p -p test/bdd/fixtures/keys/tls
	@docker run -i --rm \
		-v $(abspath .):/opt/workspace/edge-agent \
		--entrypoint "/opt/workspace/edge-agent/scripts/generate_test_keys.sh" \
		frapsoft/openssl

.PHONY: wallet-web-start
wallet-web-start: clean wallet-web-docker wallet-server-docker generate-test-config generate-test-keys mock-images
	@scripts/wallet_web_start.sh

# starting wallet-web in dev mode for hot deployment
.PHONY: wallet-web-dev-start
wallet-web-dev-start:
	@scripts/wallet_web_dev_start.sh

.PHONY: generate-test-config
generate-test-config:
	@/bin/bash scripts/generate_test_config.sh

.PHONY: bdd-test
bdd-test: bdd-test-wallet-web bdd-test-wallet-server

.PHONY: bdd-test-wallet-web
bdd-test-wallet-web: clean wallet-web-docker wallet-server-docker generate-test-config generate-test-keys mock-images
	@scripts/check_js_integration.sh

.PHONY: bdd-test-wallet-server
bdd-test-wallet-server: clean wallet-web-docker wallet-server-docker generate-test-config generate-test-keys mock-images
	@scripts/check_wallet_server_integration.sh

.PHONY: mock-bddtest-login-consent-docker
mock-bddtest-login-consent-docker:
	@echo "Building mock login consent server for BDD tests..."
	@cd test/bdd/mock/bddtest-login-consent-server && docker build -f image/Dockerfile --build-arg GO_VER=$(GO_VER) --build-arg ALPINE_VER=$(ALPINE_VER) -t edgeagent/mockbddtestloginconsent:latest .

.PHONY: mock-demo-login-consent-docker
mock-demo-login-consent-docker:
	@echo "Building login consent server for demo..."
	@cd test/bdd/mock/demo-login-consent-server && docker build -f image/Dockerfile --build-arg GO_VER=$(GO_VER) --build-arg ALPINE_VER=$(ALPINE_VER) -t edgeagent/demologinconsent:latest .

.PHONY: mock-images
mock-images: mock-bddtest-login-consent-docker mock-demo-login-consent-docker

.PHONY: clean
clean:
	@rm -Rf ./build
	@rm -Rf ./cmd/wallet-web/dist
	@rm -Rf ./cmd/wallet-web/node_modules
	@rm -Rf ./test/bdd/fixtures/wallet-web/config
