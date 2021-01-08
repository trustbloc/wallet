#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

ARG GO_VER
ARG ALPINE_VER

FROM golang:${GO_VER}-alpine${ALPINE_VER} as golang
RUN apk add --no-cache \
	git \
	libtool \
	make;
ADD . /opt/workspace/edge-agent
WORKDIR /opt/workspace/edge-agent
ENV EXECUTABLES go git

FROM golang as golang_build
ARG GO_TAGS
RUN GO_TAGS=${GO_TAGS} make wallet-server

FROM alpine:${ALPINE_VER} as base
LABEL org.opencontainers.image.source https://github.com/trustbloc/edge-agent

COPY --from=golang_build /opt/workspace/edge-agent/build/bin/wallet-server /usr/local/bin/wallet-server
ENV PATH=/usr/local/bin/wallet-server:$PATH

ENTRYPOINT ["wallet-server"]
