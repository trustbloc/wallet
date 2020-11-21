# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

FROM nginx:latest

COPY build/bin/user-agent/ /usr/share/nginx/www/
COPY images/user-agent/templates/ /etc/nginx/templates/

# defines environment variables
# NOTE: nginx will not start without them since we are using them in templates/default.conf.template
ENV HTTP_RESOLVER_URL=
ENV AGENT_DEFAULT_LABEL=
ENV AUTO_ACCEPT=false
ENV LOG_LEVEL=
ENV INDEXEDDB_NAMESPACE=
ENV BLOC_DOMAIN=
ENV WALLET_MEDIATOR_URL=
ENV CREDENTIAL_MEDIATOR_URL=
ENV BLINDED_ROUTING=
ENV STORAGE_TYPE=indexeddb
ENV EDV_SERVER_URL=
ENV EDGE_AGENT_SERVER=