#!/bin/bash -xe

#node|cloud
TARGET="${1:-node}"
MOD="github.com/YeicoLabs/laurelview"

case $TARGET in
    node)
    [ ! -d web/lvnfe/node_modules ] && (cd web/lvnfe; yarn install)
    (cd web/lvnfe; yarn run build)
    rm -fr cmd/lvnbe/build
    mkdir -p cmd/lvnbe/build
    touch cmd/lvnbe/build/.empty
    cp -fr web/lvnfe/build cmd/lvnbe/
    ;;
    cloud)
    [ ! -d web/lvcfe/node_modules ] && (cd web/lvcfe; yarn install)
    (cd web/lvcfe; yarn run build)
    rm -fr cmd/lvcbe/build
    mkdir -p cmd/lvcbe/build
    touch cmd/lvcbe/build/.empty
    cp -fr web/lvcfe/build cmd/lvcbe/
    ;;
esac
