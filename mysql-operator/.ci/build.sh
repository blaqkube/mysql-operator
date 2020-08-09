#!/bin/sh

cd $GALLY_PROJECT_ROOT
export AGENT_VERSION=$(gally list -p agent | grep GALLY_PROJECT_VERSION | cut -d'"' -f4)
docker build --build-arg agent_version=$AGENT_VERSION -t $TAG:$GALLY_PROJECT_VERSION .
docker push $TAG:$GALLY_PROJECT_VERSION
cd $GALLY_ROOT
