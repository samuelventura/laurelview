#!/bin/sh -x

[ ! -d "web/lvnfe/node_modules" ] && (cd web/lvnfe; npm i)
(cd web/lvnfe; npm run build)
rm -fr pkg/lvnbe/build
mkdir -p pkg/lvnbe/build
touch pkg/lvnbe/build/.empty
cp -fr web/lvnfe/build pkg/lvnbe/
MOD="github.com/samuelventura/laurelview"
if [[ "$OSTYPE" == "msys"* ]]; then
    #go get github.com/akavel/rsrc
    rsrc -ico icon.ico -o build/rsrc.syso
    cp build/rsrc.syso cmd/lvnss/
    cp build/rsrc.syso cmd/lvnbe/
    cp build/rsrc.syso cmd/lvnrt/
fi
go install $MOD/cmd/lvdpm
go install $MOD/cmd/lvnbe
go install $MOD/cmd/lvnrt
go install $MOD/cmd/lvnss
