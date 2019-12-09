# Copyright SecureKey Technologies Inc.
#
# SPDX-License-Identifier: Apache-2.0


WASM_EXEC              = $(shell go env GOROOT)/misc/wasm/wasm_exec.js
ISSUER_AGENT_WASM_PATH = cmd/issuer-agent-wasm
WASM_COMPONENTS        = cmd/components
HTTP_SERVER_PATH       = cmd/http-server
GOBIN_PATH             = $(abspath .)/build/bin


.PHONY: all
all: clean checks unit-test unit-test-wasm

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

.PHONY: issuer-agent-wasm
issuer-agent-wasm:
	@echo "Building issuer-agent-wasm"
	@mkdir -p ./build/bin/wasm/issuer
	@cp $(WASM_EXEC)  ./build/bin/wasm/issuer
	@cp ${ISSUER_AGENT_WASM_PATH}/*.html  ./build/bin/wasm/issuer
	@cp -r ${WASM_COMPONENTS}  ./build/bin/wasm/issuer
	@cd ${ISSUER_AGENT_WASM_PATH} && GOOS=js GOARCH=wasm go build -o ../../build/bin/wasm/issuer/issuer-agent.wasm main.go
	@gzip --best build/bin/wasm/issuer/issuer-agent.wasm

.PHONY: http-server
http-server:
	@echo "Building http-server"
	@mkdir -p ./build/bin/wasm
	@cd ${HTTP_SERVER_PATH} && go build -o ../../build/bin/http-server main.go

.PHONY: clean
clean: clean-build

.PHONY: clean-build
clean-build:
	@rm -Rf ./build
