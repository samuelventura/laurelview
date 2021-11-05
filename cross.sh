#!/bin/bash -x

export TOOLCHAIN="$HOME/.nerves/artifacts/nerves_toolchain_armv7_nerves_linux_gnueabihf-linux_x86_64-1.4.3"
export CC="$TOOLCHAIN/bin/armv7-nerves-linux-gnueabihf-cc"
export GOOS=linux 
export GOARCH=arm
export GOARM=7 
export CGO_ENABLED=1 

MOD="github.com/YeicoLabs/laurelview"
CMD=$MOD/cmd
DST=build
mkdir -p $DST

[ ! -d $TOOLCHAIN ] && (cd nfw; MIX_TARGET=bbb mix deps.get)

go build -ldflags="-extld=$CC" -o $DST/lvdpm $CMD/lvdpm
go build -ldflags="-extld=$CC" -o $DST/lvnbe $CMD/lvnbe
go build -ldflags="-extld=$CC" -o $DST/lvnup $CMD/lvnup
go build -ldflags="-extld=$CC" -o $DST/lvnss $CMD/lvnss

zip nss/priv/lvbin.zip $DST/lvdpm $DST/lvnbe $DST/lvnup
