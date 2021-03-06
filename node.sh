#!/bin/bash -xe

#node|cloud
TARGET="${1:-node}"
MOD="github.com/YeicoLabs/laurelview"

case $TARGET in
    node)
    [ ! -d web/lvnfe/node_modules ] && (cd web/lvnfe; yarn install)
    (cd web/lvnfe; yarn start)
    ;;
    cloud)
    [ ! -d web/lvcfe/node_modules ] && (cd web/lvcfe; yarn install)
    (cd web/lvcfe; yarn start)
    ;;
    setup)
    [ ! -d web/lvsfe/node_modules ] && (cd web/lvsfe; yarn install)
    (cd web/lvsfe; yarn start)
    ;;
esac
