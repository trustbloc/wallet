
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#


# Release Parameters
BASE_VERSION=0.1.5
IS_RELEASE=false

SOURCE_REPO=edge-agent
BASE_WALLET_WEB_PKG_NAME=wallet-web
BASE_WALLET_SERVER_PKG_NAME=wallet-server

RELEASE_REPO=docker.pkg.github.com/trustbloc/${SOURCE_REPO}
SNAPSHOT_REPO=docker.pkg.github.com/trustbloc-cicd/snapshot

if [ ${IS_RELEASE} = false ]
then
  EXTRA_VERSION=snapshot-$(git rev-parse --short=7 HEAD)
  PROJECT_VERSION=${BASE_VERSION}-${EXTRA_VERSION}
  PROJECT_PKG_REPO=${SNAPSHOT_REPO}
else
  PROJECT_VERSION=${BASE_VERSION}
  PROJECT_PKG_REPO=${RELEASE_REPO}
fi

export WALLET_WEB_TAG=$PROJECT_VERSION
export WALLET_WEB_PKG=${PROJECT_PKG_REPO}/${BASE_WALLET_WEB_PKG_NAME}

export WALLET_SERVER_TAG=$PROJECT_VERSION
export WALLET_SERVER_PKG=${PROJECT_PKG_REPO}/${BASE_WALLET_SERVER_PKG_NAME}
