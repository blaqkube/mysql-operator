#!/bin/sh

set -e

echo "no bundle for $GALLY_PROJECT_NAME"
echo "  * CIRCLE_TAG: $CIRCLE_TAG"
echo "  * CIRCLE_SHA1: $CIRCLE_SHA1"
exit 0

cd $GALLY_PROJECT_ROOT
export AGENT_VERSION=$(gally list -p agent | grep GALLY_PROJECT_VERSION | cut -d'"' -f4)
export IMG=quay.io/blaqkube/mysql-controller:$VERSION
export BUNDLE_IMG=quay.io/blaqkube/mysql-operator:$VERSION
make docker-build
make docker-push
make bundle
git status -s | cat
make bundle-build
docker push $BUNDLE_IMG
