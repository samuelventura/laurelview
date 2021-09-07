#!/bin/bash -x

[ ! -d "web/lvnfe/node_modules" ] && (cd web/lvnfe; npm i)
(cd web/lvnfe; npm run build)
rm -fr cmd/lvnbe/build
mkdir -p cmd/lvnbe/build
touch cmd/lvnbe/build/.empty
cp -fr web/lvnfe/build cmd/lvnbe/

