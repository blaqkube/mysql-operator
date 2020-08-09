#!/usr/bin/sh

cd $GALLY_PROJECT_ROOT
docker build -t $TAG:$GALLY_PROJECT_VERSION .
docker push $TAG:$GALLY_PROJECT_VERSION
cd $GALLY_ROOT
