#!/bin/bash

set -ex

TIMESTAMP=$(date +"%Y%m%d%H%M%S")
VERSION=$(cat package.json | jq -r .version)
FOLDER_NAME=release-$VERSION-$TIMESTAMP
TARGET_PATH=/home/pi/leds/$FOLDER_NAME
SYMLINK_PATH=/home/pi/leds/data

echo "Building project"
yarn build

echo "Uploading contents to $TARGET_PATH on the host"
rsync -a --chown=:leds --chmod=F664,D775 ./dist/ leds.nwl:$TARGET_PATH/

echo "Symlinking..."
ssh leds.nwl "rm -f $SYMLINK_PATH && ln -sf $TARGET_PATH $SYMLINK_PATH"

ssh leds.nwl "bash /home/pi/leds/cleanup-releases.sh"

