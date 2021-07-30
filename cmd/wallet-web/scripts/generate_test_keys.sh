#!/bin/sh
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
    #
# SPDX-License-Identifier: Apache-2.0
#

set -e

echo "Generating Wallet Web Test PKI"
mkdir -p test/fixtures/keys/tls
tmp=$(mktemp)
echo "subjectKeyIdentifier=hash
authorityKeyIdentifier = keyid,issuer
extendedKeyUsage = serverAuth
keyUsage = Digital Signature, Key Encipherment
subjectAltName = @alt_names
    [alt_names]
DNS.1 = localhost
DNS.2 = testnet.orb.local
DNS.3 = hub-router.trustbloc.local" >> "$tmp"

#create CA
openssl ecparam -name prime256v1 -genkey -noout -out test/fixtures/keys/tls/ec-cakey.pem
openssl req -new -x509 -key test/fixtures/keys/tls/ec-cakey.pem -subj "/C=CA/ST=ON/O=Example Internet CA Inc.:CA Sec/OU=CA Sec" -out test/fixtures/keys/tls/ec-cacert.pem

#create TLS creds
openssl ecparam -name prime256v1 -genkey -noout -out test/fixtures/keys/tls/ec-key.pem
openssl req -new -key test/fixtures/keys/tls/ec-key.pem -subj "/C=CA/ST=ON/O=Example Inc.:Aries-Framework-Go/OU=Aries-Framework-Go/CN=*.example.com" -out test/fixtures/keys/tls/ec-key.csr
openssl x509 -req -in test/fixtures/keys/tls/ec-key.csr -CA test/fixtures/keys/tls/ec-cacert.pem -CAkey test/fixtures/keys/tls/ec-cakey.pem -CAcreateserial -extfile "$tmp" -out test/fixtures/keys/tls/ec-pubCert.pem -days 365


echo "done generating wallet web Test PKI"
