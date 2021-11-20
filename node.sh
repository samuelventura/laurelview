#!/bin/bash -x

#node|cloud
TARGET="${1:-node}"
MOD="github.com/YeicoLabs/laurelview"

case $TARGET in
    node)
    [ ! -d web/lvnfe/node_modules ] && (cd web/lvnfe; npm i)
    (cd web/lvnfe; npm start)
    ;;
    cloud)
    [ ! -d web/lvcfe/node_modules ] && (cd web/lvcfe; npm i)
    (cd web/lvcfe; npm start)
    ;;
esac
