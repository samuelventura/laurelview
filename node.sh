#!/bin/bash -x

APP="${1:-cloud}"

case $APP in
    node)
    [ ! -d web/lvnfe/node_modules ] && (cd web/lvnfe; npm i)
    (cd web/lvnfe; npm start)
    ;;
    cloud)
    [ ! -d web/lvcfe/node_modules ] && (cd web/lvcfe; npm i)
    (cd web/lvcfe; npm start)
    ;;
esac
