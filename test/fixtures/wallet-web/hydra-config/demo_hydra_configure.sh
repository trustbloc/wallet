#!/bin/sh
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

echo "Registering demo-hub-auth with third party provider"
# will use --skip-tls-verify because hydra doesn't trust self-signed certificate
# remove it when using real certificate
hydra clients create \
    --endpoint https://demo-hydra.trustbloc.local:7778 \
    --id client-id \
    --secret client-secret \
    --grant-types authorization_code,refresh_token \
    --response-types code,id_token \
    --scope openid,profile,email \
    --skip-tls-verify \
    --callbacks https://hub-auth.trustbloc.local:8044/oauth2/callback
echo "Finished registering demo-hub-auth with third party provider"
