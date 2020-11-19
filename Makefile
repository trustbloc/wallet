# Copyright SecureKey Technologies Inc.
#
# SPDX-License-Identifier: Apache-2.0


USER_AGENT_SUPPORT_PATH	= cmd/user-agent-support

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

.PHONY: user-agent
user-agent:
	@scripts/build_agent_wasm.sh

.PHONY: user-agent-docker
user-agent-docker: clean user-agent
	@echo "Building user agent (UI) docker image"
	@docker build -f ./images/user-agent/Dockerfile --no-cache -t $(DOCKER_OUTPUT_NS)/$(REPO_IMAGE_NAME)/user-agent:latest .

.PHONY: user-agent-support
user-agent-support:
	@echo "Building user-agent-support"
	@cd ${USER_AGENT_SUPPORT_PATH} && go build -o ../../build/bin/user-agent-support main.go

.PHONY: user-agent-support-docker
user-agent-support-docker: clean
	@echo "Building user agent support (backend) docker image"
	@docker build -f ./images/user-agent-support/Dockerfile --no-cache -t $(DOCKER_OUTPUT_NS)/$(REPO_IMAGE_NAME)/user-agent-support:latest \
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

.PHONY: user-agent-start
user-agent-start: clean user-agent-docker user-agent-support-docker generate-test-config generate-test-keys mock-images
	@scripts/user_agent_start.sh

# starting user agent in dev mode for hot deployment
.PHONY: user-agent-dev-start
user-agent-dev-start:
	@scripts/user_agent_dev_start.sh

.PHONY: generate-test-config
generate-test-config:
	@/bin/bash scripts/generate_test_config.sh

.PHONY: bdd-test
bdd-test: bdd-test-user-agent bdd-test-user-agent-support

.PHONY: bdd-test-user-agent
bdd-test-user-agent: clean user-agent-docker user-agent-support-docker generate-test-config generate-test-keys mock-images
	@scripts/check_js_integration.sh

.PHONY: bdd-test-user-agent-support
bdd-test-user-agent-support: clean user-agent-docker user-agent-support-docker generate-test-config generate-test-keys mock-images
	@scripts/check_user_agent_support_integration.sh

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
	@rm -Rf ./cmd/user-agent/dist
	@rm -Rf ./cmd/user-agent/node_modules
	@rm -Rf ./test/bdd/fixtures/user-agent/config
