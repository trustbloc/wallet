#!/bin/sh
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

echo "Creating hub-auth client"
# will use --skip-tls-verify because hydra doesn't trust self-signed certificate
# remove it when using real certificate
hydra clients create \
    --endpoint https://bdd-hydra.trustbloc.local:4445 \
    --id client-id \
    --secret client-secret \
    --grant-types authorization_code,refresh_token \
    --response-types code,id_token \
    --scope openid,profile,email \
    --skip-tls-verify \
    --callbacks https://localhost:8045/oauth2/callback
echo "Finished creating hub-auth client"
