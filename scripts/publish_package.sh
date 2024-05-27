#!/usr/bin/env bash

set -e

USAGE_MSG='Usage: BRANCH=[BRANCH] publish_package.sh'
if [ -z "$BRANCH" ]
then
    (>&2 echo 'You should provide branch')
    echo "$USAGE_MSG"
    exit 1
fi
if [ -z "$NODE_AUTH_TOKEN" ]
then
    (>&2 echo 'You should provide a node auth token')
    echo "$USAGE_MSG"
    exit 2
fi

cd "$(dirname "$0")/.."

BRANCH=$(echo $BRANCH | tr [:upper:] [:lower:] | tr -d [:space:])
VERSION=$(BRANCH=$BRANCH ./scripts/calculate_version.sh)

TAG=""
if ! [[ $BRANCH == 'stable' ]]
then
    TAG="--tag $BRANCH"
fi

if [[ "$VERSION" == *-stable.0 ]]
then
    VERSION=${VERSION%-stable.0}
fi

echo "Using $VERSION as a new version"

cp LICENSE contracts
cp README.md contracts

cd contracts

# set version
package="$(jq --arg v "$VERSION" '.version = $v' package.json)"
echo -E "${package}" > package.json

touch yarn.lock
yarn

yarn config set enableImmutableInstalls false
yarn config set npmAuthToken "$NODE_AUTH_TOKEN"
yarn npm publish --access public $TAG
