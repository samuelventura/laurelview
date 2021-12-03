#!/bin/bash -xe

MOD="github.com/YeicoLabs/laurelview"
BDST=build/msys

if [[ "$OSTYPE" == "msys"* ]]; then
    [[ $(type -P rsrc) ]] || go install github.com/akavel/rsrc
    mkdir -p $BDST
    rsrc -ico icon.ico -o $BDST
    cp $BDST/rsrc.syso cmd/lvnss/
    cp $BDST/rsrc.syso cmd/lvnbe/
    cp $BDST/rsrc.syso cmd/lvnup/
    cp $BDST/rsrc.syso cmd/lvdpm/
fi

#remove lvtry.exe
#rm /c/Users/samuel/go/bin/lv*.exe
go install $MOD/cmd/lvdpm
go install $MOD/cmd/lvnbe
go install $MOD/cmd/lvnup
go install $MOD/cmd/lvnss
go install $MOD/cmd/lvcbe

#for elixir testing
IDST=~/go/bin

echo "LV_NUP_ENDPOINT=127.0.0.1:8800" > $IDST/lvnup.env
echo "LV_DPM_ENDPOINT=127.0.0.1:8801" > $IDST/lvdpm.env
echo "LV_NBE_ENDPOINT=0.0.0.0:8800" > $IDST/lvnbe.env
echo "LV_NBE_DEBUG=127.0.0.1:8802" >> $IDST/lvnbe.env
