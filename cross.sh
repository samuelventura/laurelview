#!/bin/bash -x

#See memo: Athasha ToughBC - Linux Buildroot BBB
export PATH=~/local/go/bin/:$PATH
export CC="$HOME/buildroot/output/host/bin/arm-buildroot-linux-uclibcgnueabihf-gcc"
export GOOS=linux 
export GOARCH=arm
export GOARM=7 
export CGO_ENABLED=1 

MOD="github.com/YeicoLabs/laurelview"
CMD=$MOD/cmd
DST=~/nfsroot/usr/bin

go build -ldflags="-extld=$CC" -o $DST/lvdpm $CMD/lvdpm
go build -ldflags="-extld=$CC" -o $DST/lvnbe $CMD/lvnbe
go build -ldflags="-extld=$CC" -o $DST/lvnup $CMD/lvnup
go build -ldflags="-extld=$CC" -o $DST/lvnss $CMD/lvnss
