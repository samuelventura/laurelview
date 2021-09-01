#!/bin/sh -x

[ ! -d "web/lvnfe/node_modules" ] && (cd web/lvnfe; npm i)
(cd web/lvnfe; npm run build)
rm -fr cmd/lvnbe/build
mkdir -p cmd/lvnbe/build
touch cmd/lvnbe/build/.empty
cp -fr web/lvnfe/build cmd/lvnbe/
MOD="github.com/samuelventura/laurelview"
if [[ "$OSTYPE" == "msys"* ]]; then
    #go get github.com/akavel/rsrc
    rsrc -ico icon.ico -o build/rsrc.syso
    cp build/rsrc.syso cmd/lvnss/
    cp build/rsrc.syso cmd/lvnbe/
fi
go install $MOD/cmd/lvdpm
go install $MOD/cmd/lvnbe
go install $MOD/cmd/lvnup
go install $MOD/cmd/lvnss
