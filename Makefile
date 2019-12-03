# Copyright SecureKey Technologies Inc.
#
# SPDX-License-Identifier: Apache-2.0



.PHONY: all
all: checks

.PHONY: checks
checks: license lint

.PHONY: lint
lint:
	@scripts/check_lint.sh

.PHONY: license
license:
	@scripts/check_license.sh
