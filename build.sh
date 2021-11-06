#!/bin/bash -x

MOD="github.com/YeicoLabs/laurelview"

if [[ "$OSTYPE" == "msys"* ]]; then
    [[ $(type -P rsrc) ]] || go install github.com/akavel/rsrc
    mkdir -p build
    rsrc -ico icon.ico -o build/rsrc.syso
    cp build/rsrc.syso cmd/lvnss/
    cp build/rsrc.syso cmd/lvnbe/
    cp build/rsrc.syso cmd/lvnup/
    cp build/rsrc.syso cmd/lvdpm/
fi
#remove lvtry.exe
#rm /c/Users/samuel/go/bin/lv*.exe
go install $MOD/cmd/lvdpm
go install $MOD/cmd/lvnbe
go install $MOD/cmd/lvnup
go install $MOD/cmd/lvnss

#for elixir testing
DST=~/go/bin

echo "LV_NUP_ENDPOINT=127.0.0.1:8800" > $DST/lvnup.env
echo "LV_DPM_ENDPOINT=127.0.0.1:8801" > $DST/lvdpm.env
echo "LV_NBE_ENDPOINT=0.0.0.0:8800" > $DST/lvnbe.env
echo "LV_NBE_DEBUG=127.0.0.1:8802" >> $DST/lvnbe.env
