#!/bin/sh

set -e

cd $GALLY_PROJECT_ROOT
make test
export AGENT_CODE=$(grep "DefaultAgentVersion =" main.go | cut -d '"' -f2)
docker build --build-arg agent_version=$AGENT_CODE -t $TAG:$GALLY_PROJECT_VERSION .
docker push $TAG:$GALLY_PROJECT_VERSION
cd $GALLY_ROOT
