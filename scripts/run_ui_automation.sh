#!/usr/bin/env bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

echo "Running $0"

ROOT=`pwd`

npm -v
echo "starting containers..."
cd $ROOT/test/fixtures/wallet-web
(source .env && docker-compose -f docker-compose-demo.yml -f docker-compose-server.yml -f docker-compose-web.yml down && docker-compose -f docker-compose-demo.yml -f docker-compose-server.yml -f docker-compose-web.yml up  --force-recreate -d)

sleep 45

echo "running healthcheck..."

# healthCheck function
RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
AQUA=$(tput setaf 6)
NONE=$(tput sgr0)
healthCheck() {
	sleep 2
    n=0
    maxAttempts=200
    if [ "" != "$4" ]
    then
	   maxAttempts=$4
    fi

	echo "running health check : app=$1 url=$2 timeout=$maxAttempts seconds"

	until [ $n -ge $maxAttempts ]
	do
	  response=$(curl -H 'Cache-Control: no-cache' -o /dev/null -s -w "%{http_code}" --insecure "$2")
	  if [ "$response" == "$3" ]
	  then
	    echo "${GREEN}$1 $2 is up ${NONE}"
		break
	   fi
	   n=$((n+1))
	   if [ $n -eq $maxAttempts ]
	   then
		 echo "${RED}failed health check : app=$1 url=$2 responseCode=$response ${NONE}"
	   fi
	   sleep 1
	done
}

# healthcheck
healthCheck wallet-web https://wallet.trustbloc.local:8091/healthcheck 200
healthCheck wallet-web-2 https://wallet-2.trustbloc.local:8071/healthcheck 200
healthCheck wallet-server https://localhost:8090/healthcheck 200

echo "running tests..."
cd $ROOT/test/ui-automation
npm run test && npm run report
if [ $? -ne 0 ]
then
	cd $ROOT/test/fixtures/wallet-web
	docker-compose -f docker-compose-demo.yml -f docker-compose-server.yml -f docker-compose-web.yml logs  --no-color >& docker-compose.log
    
	exit 1
fi

echo "stopping containers..."
cd $ROOT/test/fixtures/wallet-web

(source .env && docker-compose -f docker-compose-demo.yml -f docker-compose-server.yml -f docker-compose-web.yml down --remove-orphans)
