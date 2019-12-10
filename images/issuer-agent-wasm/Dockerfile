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
	make;
ADD . /opt/workspace/edge-agent
WORKDIR /opt/workspace/edge-agent
ENV EXECUTABLES go git

FROM golang as build
ARG GO_TAGS
RUN GO_TAGS=${GO_TAGS} make issuer-agent-wasm http-server


FROM alpine:${ALPINE_VER} as base
COPY --from=build /opt/workspace/edge-agent/build/bin/wasm/issuer /usr/local/bin/issuer
COPY --from=build /opt/workspace/edge-agent/build/bin/http-server /usr/local/bin/issuer
ENV PATH=/usr/local/bin/issuer:$PATH
ENV HTTP_SERVER_WASM_PATH=/usr/local/bin/issuer
ENTRYPOINT ["http-server"]
