# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

FROM nginx:latest
LABEL org.opencontainers.image.source https://github.com/trustbloc/wallet

COPY build/bin/wallet-web/ /usr/share/nginx/www/
COPY images/wallet-web/templates/ /etc/nginx/templates/

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
ENV STORAGE_TYPE=edv
ENV EDV_SERVER_URL=
ENV KMS_TYPE=local
ENV USE_EDV_CACHE=false
ENV KMS_TYPE=local
ENV LOCAL_KMS_PASSPHRASE=demo
ENV EDV_CLEAR_CACHE=
ENV EDGE_AGENT_SERVER=
ENV USE_EDV_BATCH=false
ENV EDV_BATCH_SIZE=0
ENV CACHE_SIZE=100
ENV DID_ANCHOR_ORIGIN=
ENV SIDETREE_TOKEN=
ENV STATIC_ASSETS=/etc/static
ENV STATIC_ASSETS_URL=
ENV WALLET_WEB_URL=
ENV CONTEXT_PROVIDER_URL=
ENV DIDCOMM_MEDIA_TYPE_PROFILES=
ENV DIDCOMM_KEY_TYPE=
ENV DIDCOMM_KEY_AGREEMENT_TYPE=
ENV WEB_SOCKET_READ_LIMIT=
ENV KMS_SERVER_URL=
