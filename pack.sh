#!/bin/bash -x

[ ! -d "web/lvnfe/node_modules" ] && (cd web/lvnfe; npm i)
(cd web/lvnfe; npm run build)
rm -fr cmd/lvnbe/build
mkdir -p cmd/lvnbe/build
touch cmd/lvnbe/build/.empty
cp -fr web/lvnfe/build cmd/lvnbe/

#rsync -r cmd/lvnbe/build pi@10.77.3.143:bin/laurelview/cmd/lvnbe/
