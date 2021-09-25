#!/bin/bash -x

MOD="github.com/samuelventura/laurelview"

if [[ "$OSTYPE" == "msys"* ]]; then
    #go get github.com/akavel/rsrc
    rsrc -ico icon.ico -o build/rsrc.syso
    cp build/rsrc.syso cmd/lvnss/
    cp build/rsrc.syso cmd/lvnbe/
fi
#remove lvtry.exe
#rm /c/Users/samuel/go/bin/lv*.exe
go install $MOD/cmd/lvdpm
go install $MOD/cmd/lvnbe
go install $MOD/cmd/lvnup
go install $MOD/cmd/lvnss
