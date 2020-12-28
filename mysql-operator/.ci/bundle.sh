#!/bin/bash

set -e

echo "no bundle for $GALLY_PROJECT_NAME"
echo "  * CIRCLE_TAG:  $CIRCLE_TAG"
echo "  * CIRCLE_SHA1: $CIRCLE_SHA1"
echo "  * VERSION:     $VERSION"
echo "  * PREVIOUS:    $PREV_VERSION"

if [[ "v${VERSION}" != "$CIRCLE_TAG" ]]; then
  echo "Error, only tag v${VERSION} is allowed"
  exit 1
fi

echo "Building version v${VERSION}..."
cd $GALLY_PROJECT_ROOT
export AGENT_CODE=$(grep "DefaultAgentVersion =" main.go | cut -d '"' -f2)
export AGENT_VERSION=$(gally list -p agent | grep GALLY_PROJECT_VERSION | cut -d'"' -f4)
if [[ "${AGENT_CODE}" != "${AGENT_VERSION}" ]]; then
  echo "Agent version does not match code(${AGENT_CODE}) != main(${AGENT_VERSION})..."
  exit 1
fi

# export EXISTS=$(grep "  replaces: $PREV_VERSION" config/manifests/bases/mysql-operator.clusterserviceversion.yaml -c)
# if [[ "${EXISTS}" != "1" ]]; then
#   echo "Previous bundle version is not in CSV..."
#   exit 1
# fi

export AGENT_IMG=quay.io/blaqkube/mysql-agent:${AGENT_CODE}
docker pull $AGENT_IMG
docker tag $AGENT_IMG quay.io/blaqkube/mysql-agent:$VERSION
docker push quay.io/blaqkube/mysql-agent:$VERSION
export IMG=quay.io/blaqkube/mysql-controller:$VERSION
export BUNDLE_IMG=quay.io/blaqkube/mysql-operator:$VERSION
make docker-build
make docker-push
make bundle
git status -s | cat
make bundle-build
docker push $BUNDLE_IMG

docker pull quay.io/blaqkube/mysql-operator:$VERSION
# opm index add --container-tool docker \
#   --bundles quay.io/blaqkube/mysql-operator:$VERSION \
#   --from-index quay.io/blaqkube/operators-index:$PREV_VERSION
#   --tag quay.io/blaqkube/operators-index:$VERSION
opm index add --container-tool docker \
  --bundles quay.io/blaqkube/mysql-operator:$VERSION \
  --tag quay.io/blaqkube/operators-index:$VERSION
docker push quay.io/blaqkube/operators-index:$VERSION
