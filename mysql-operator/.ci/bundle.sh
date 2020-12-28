#!/bin/bash

set -e

echo "no bundle for $GALLY_PROJECT_NAME"
echo "  * CIRCLE_TAG:  $CIRCLE_TAG"
echo "  * CIRCLE_SHA1: $CIRCLE_SHA1"
echo "  * VERSION:     $VERSION"

if [[ "v${VERSION}" != "$CIRCLE_TAG" ]]; then
  echo "Error, only tag v${VERSION} is allowed"
  exit 1
fi

echo "Building version v${VERSION}..."
cd $GALLY_PROJECT_ROOT
export AGENT_CODE=$(grep "DefaultAgentVersion =" main.go | cut -d '"' -f2)
export AGENT_VERSION=$(gally list -p agent | grep GALLY_PROJECT_VERSION | cut -d'"' -f4)
if [[ "v${VERSION}" == "CIRCLE_TAG" ]]; then
export IMG=quay.io/blaqkube/mysql-controller:$VERSION
export BUNDLE_IMG=quay.io/blaqkube/mysql-operator:$VERSION
make docker-build
make docker-push
make bundle
git status -s | cat
make bundle-build
docker push $BUNDLE_IMG 
