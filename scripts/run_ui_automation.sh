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
cd $ROOT/test/bdd/fixtures/wallet-web
(source .env && docker-compose -f docker-compose-server.yml -f docker-compose-web.yml down && docker-compose -f docker-compose-server.yml -f docker-compose-web.yml up  --force-recreate -d)

sleep 60

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
healthCheck wallet-web https://user-ui-agent.example.com:8091/healthcheck 200
healthCheck wallet-web-2 https://second-ui-user-agent.example.com:8071/healthcheck 200
healthCheck wallet-server https://localhost:8077/healthcheck 200

echo "running tests..."
cd $ROOT/test/ui-automation
npm run test && npm run report

echo "stopping containers..."
cd $ROOT/test/bdd/fixtures/wallet-web

(source .env && docker-compose -f docker-compose-server.yml -f docker-compose-web.yml down --remove-orphans)
