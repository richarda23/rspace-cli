#!/bin/bash

##
# Builds binaries of rspace-cli for 3 amd64 architectures for win,max,linux
# Requires single argument, a version numbepr
# 2nd argument can be 'b' (build only, default); 'p' (publish only) or 'bp' (build and publish)
# 
# Publishing pushes to bintray creating a new version; will fail if this version already exists
# The files uploaded need manual publishing on bintray website
# 
# 'build' rrequires a pre-created distribution folder  'dist' with 3 subfolders win.mac,linux
#
##

VERSION=$1
CMD=$2 ##'p', 'b', 'pb'
if [[ -z "$CMD" ]]; then
    CMD="b" ## build only is default
fi
if [[ -z "VERSION" ]]; then
    echo "Version number required"
    exit 1
fi
DIST_DIR=dist
BUILD_LOG=$DIST_DIR/build.log
WIN_EXE=rspace-$VERSION.exe
MAC_EXE=rspace-$VERSION
LINUX_EXE=rspace-$VERSION


function build {

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
}
## publishes to bin tray if 2nd arg is 'p' or 'bp'
function publish {
    ##  this file should contain 2 lines, out side version control
    ## APIKEY=yourbintrayapiky
    ## BINTRAY_USERNAME=yourbintray_username
    source .bintray-api.sh

    PKG="rspace-cli"
    BINTRAY_URL="https://api.bintray.com/content"
    FILE_UPLOAD_BASE=$BINTRAY_URL/$BINTRAY_USERNAME/rspace-cli/$PKG/$VERSION
    ## WIN
    echo "Uploading  $WIN_EXE"
    curl -T $DIST_DIR/win/$WIN_EXE -u$BINTRAY_USERNAME:$APIKEY "$FILE_UPLOAD_BASE/amd-64-windows/rspace.exe"
    ## MAC
    echo "Uploading $MAC_EXE"
    curl -T $DIST_DIR/mac/$MAC_EXE -u$BINTRAY_USERNAME:$APIKEY "$FILE_UPLOAD_BASE/amd-64-macosx/rspace"
    ## LINUX
    echo "Uploading $LINUX_EXE"
    curl -T $DIST_DIR/linux/$LINUX_EXE -u$BINTRAY_USERNAME:$APIKEY "$FILE_UPLOAD_BASE/amd-64-linux/rspace"
}
#### main script####
if [[ $CMD == "b" || $CMD == "bp" ]]; then
    echo "building..."
    build;
fi
if [[ $CMD == "p" || $CMD == "bp" ]]; then
    echo "publishing"
    publish;
fi


