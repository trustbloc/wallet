#!/bin/sh
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

echo "Creating demo wallet-server client with hub-auth"
# will use --skip-tls-verify because hydra doesn't trust self-signed certificate
# remove it when using real certificate
hydra clients create \
    --endpoint https://demo-hub-auth-hydra.trustbloc.local:5556 \
    --id client-id \
    --secret client-secret \
    --grant-types authorization_code,refresh_token \
    --response-types code,id_token \
    --scope openid,profile,email \
    --skip-tls-verify \
    --callbacks https://user-agent.example.com:8090/oidc/callback,https://second-user-agent.example.com:8070/oidc/callback
echo "Finish creating demo wallet-server clients with hub-auth"
