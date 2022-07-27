#!/bin/sh
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e


echo "Generating wallet Test PKI"

cd /opt/workspace/wallet
mkdir -p test/fixtures/keys/tls
tmp=$(mktemp)
echo "subjectKeyIdentifier=hash
authorityKeyIdentifier = keyid,issuer
extendedKeyUsage = serverAuth
keyUsage = Digital Signature, Key Encipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
DNS.2 = *.trustbloc.local
DNS.3 = testnet.orb.local
DNS.4 = hydra
DNS.5 = *.example.com
DNS.6 = wallet.trustbloc.local
DNS.8 = mediator.trustbloc.local
DNS.9 = uni-resolver-web.trustbloc.local
DNS.10 = file-server.trustbloc.local
DNS.11 = demo-adapter.trustbloc.local" >> "$tmp"

CERT_CA="test/fixtures/keys/tls/ec-cacert.pem"
if [ ! -f "$CERT_CA" ]; then
#create CA
openssl ecparam -name prime256v1 -genkey -noout -out test/fixtures/keys/tls/ec-cakey.pem
openssl req -new -x509 -key test/fixtures/keys/tls/ec-cakey.pem -subj "/C=CA/ST=ON/O=Example Internet CA Inc.:CA Sec/OU=CA Sec" -out $CERT_CA
else
    echo "Skipping CA generation - already exists"
fi

#create TLS creds
openssl ecparam -name prime256v1 -genkey -noout -out test/fixtures/keys/tls/ec-key.pem
openssl req -new -key test/fixtures/keys/tls/ec-key.pem -subj "/C=CA/ST=ON/O=Example Inc.:wallet/OU=wallet/CN=localhost" -out test/fixtures/keys/tls/ec-key.csr
openssl x509 -req -in test/fixtures/keys/tls/ec-key.csr -CA test/fixtures/keys/tls/ec-cacert.pem -CAkey test/fixtures/keys/tls/ec-cakey.pem -CAcreateserial -extfile "$tmp" -out test/fixtures/keys/tls/ec-pubCert.pem -days 365

#create session cookie keys
mkdir -p test/fixtures/keys/session_cookies
openssl rand -out test/fixtures/keys/session_cookies/auth.key 32
openssl rand -out test/fixtures/keys/session_cookies/enc.key 32

#create master key for secret lock
openssl rand 32 | base64 | sed 's/+/-/g; s/\//_/g' > test/fixtures/keys/tls/secret-lock.key

#create private key for GNAP signer (kms)
openssl ecparam -name prime256v1 -genkey -noout -out test/fixtures/keys/gnap-kms-priv-key.pem

#create private key for GNAP signer (edv)
openssl ecparam -name prime256v1 -genkey -noout -out test/fixtures/keys/gnap-edv-priv-key.pem

mkdir -p test/fixtures/keys/device
openssl ecparam -name prime256v1 -genkey -noout -out test/fixtures/keys/device/ec-cakey.pem
openssl req -new -x509 -key test/fixtures/keys/device/ec-cakey.pem -subj "/C=CA/ST=ON/O=Example Auth Device Inc.:CA Sec/OU=CA Sec" -out test/fixtures/keys/device/ec-cacert.pem

echo "done generating wallet PKI"
