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
RUN GO_TAGS=${GO_TAGS} make user-agent-support

FROM alpine:${ALPINE_VER} as base

COPY --from=golang_build /opt/workspace/edge-agent/build/bin/user-agent-support /usr/local/bin/user-agent-support
ENV PATH=/usr/local/bin/user-agent-support:$PATH

ENTRYPOINT ["user-agent-support"]
