#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

ARG GO_VER
ARG ALPINE_VER

FROM golang:${GO_VER}-alpine${ALPINE_VER} as golang
RUN apk add --no-cache \
	gcc \
	musl-dev \
	git \
	libtool \
	bash \
	npm \
	make;
ADD . /opt/workspace/edge-agent
WORKDIR /opt/workspace/edge-agent
ENV EXECUTABLES go git

FROM golang as golang_build
ARG GO_TAGS
ARG NAME
ARG GITHUB_TOKEN
ENV GITHUB_TOKEN=${GITHUB_TOKEN}
RUN GO_TAGS=${GO_TAGS} make ${NAME}-agent-wasm http-server

FROM alpine:${ALPINE_VER} as base
ARG NAME
COPY --from=golang_build /opt/workspace/edge-agent/build/bin/wasm/${NAME} /usr/local/bin/${NAME}
COPY --from=golang_build /opt/workspace/edge-agent/build/bin/http-server /usr/local/bin/${NAME}
ENV PATH=/usr/local/bin/${NAME}:$PATH
ENV HTTP_SERVER_WASM_PATH=/usr/local/bin/${NAME}
ENTRYPOINT ["http-server"]
