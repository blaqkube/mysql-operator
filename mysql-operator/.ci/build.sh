#!/bin/sh

set -e

cd $GALLY_PROJECT_ROOT
export AGENT_VERSION=$(gally list -p agent | grep GALLY_PROJECT_VERSION | cut -d'"' -f4)
export IMG=quay.io/blaqkube/mysql-controller:$VERSION
export BUNDLE_IMG=quay.io/blaqkube/mysql-operator:$VERSION
make docker-build
make docker-push
make bundle
git status -s
git diff
make bundle-build
