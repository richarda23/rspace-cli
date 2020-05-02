#!/bin/bash

##
# Builds binaries of rspace-cli for 3 amd64 architectures for win,max,linux
# Requires single argument, a version number
# Requires a pre-created distribution folder  'dist' with 3 subfolders win.mac,linux
#
##

VERSION=$1
if [[ -z "VERSION" ]]; then
    echo "Version number required"
    exit 1
fi
DIST_DIR=dist
BUILD_LOG=$DIST_DIR/build.log
WIN_EXE=rspace-$VERSION.exe
MAC_EXE=rspace-$VERSION
LINUX_EXE=rspace-$VERSION

### remove any old builds of same name
rm $DIST_DIR/win/$WIN_EXE
rm $DIST_DIR/mac/$MAC_EXE
rm $DIST_DIR/linux/$LINUX_EXE


echo "Building Windows amd64"
env GOOS=windows GOARCH=amd64 go build -o $DIST_DIR/win/$WIN_EXE
echo "Building Mac amd64"

env GOOS=darwin GOARCH=amd64 go build -o $DIST_DIR/mac/$MAC_EXE
echo "Building Linux amd64"
env GOOS=linux GOARCH=amd64 go build -o $DIST_DIR/linux/$LINUX_EXE 
echo "Done, appending to build log $BUILD_LOG"
echo "$VERSION built on $(date)" >> $BUILD_LOG


