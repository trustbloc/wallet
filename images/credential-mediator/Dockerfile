#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

ARG NODE_VER
ARG ALPINE_VER
ARG GITHUB_TOKEN

FROM node:${NODE_VER}-alpine${ALPINE_VER} as node
RUN apk add --no-cache \
	git \
	bash ;

WORKDIR /opt/workspace
ENV GITHUB_TOKEN=${GITHUB_TOKEN}
RUN git clone -b master https://${GITHUB_TOKEN}:x-oauth-basic@github.com/digitalbazaar/authorization.io.git

WORKDIR /opt/workspace/authorization.io
RUN npm install --unsafe-perm

EXPOSE 33443 33080
# TODO to be run in production settings [#279]
CMD [ "node", "authorization.dev.js" ]
