#!/bin/bash -xe

MOD="github.com/YeicoLabs/laurelview"
BDST=build/msys

#for elixir testing
IDST=~/go/bin

if [[ "$OSTYPE" == "msys"* ]]; then
    [[ $(type -P rsrc) ]] || go get github.com/akavel/rsrc
    mkdir -p $BDST
    rsrc -ico icon.ico -o $BDST/rsrc.syso
    cp $BDST/rsrc.syso cmd/lvnss/
    cp $BDST/rsrc.syso cmd/lvnbe/
    cp $BDST/rsrc.syso cmd/lvnup/
    cp $BDST/rsrc.syso cmd/lvdpm/
    cp $BDST/rsrc.syso cmd/lvsbe/
    cp $BDST/rsrc.syso cmd/lvsss/
	IDST=/c/users/samuel/go/bin
fi

#remove lvtry.exe
#rm /c/Users/samuel/go/bin/lv*.exe
mkdir -p cmd/lvnbe/build/
mkdir -p cmd/lvcbe/build/
mkdir -p cmd/lvsbe/build/
touch cmd/lvnbe/build/empty.txt
touch cmd/lvcbe/build/empty.txt
touch cmd/lvsbe/build/empty.txt

go install $MOD/cmd/lvdpm
go install $MOD/cmd/lvnbe
go install $MOD/cmd/lvnup
go install $MOD/cmd/lvnss
#fails on windows
#go install $MOD/cmd/lvcbe
go install $MOD/cmd/lvsbe
go install $MOD/cmd/lvsss

echo "LV_NUP_ENDPOINT=127.0.0.1:8800" > $IDST/lvnup.env
echo "LV_DPM_ENDPOINT=127.0.0.1:8801" > $IDST/lvdpm.env
echo "LV_NBE_ENDPOINT=0.0.0.0:8800" > $IDST/lvnbe.env
echo "LV_NBE_DEBUG=127.0.0.1:8802" >> $IDST/lvnbe.env
echo "LV_SBE_ENDPOINT=127.0.0.1:8803" > $IDST/lvsbe.env

