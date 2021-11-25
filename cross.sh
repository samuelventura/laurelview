#!/bin/bash -x

# bbb|rpi4
BOARD="${1:-bbb}"

case $OSTYPE in
    linux*)
    DEVARCH=linux_x86_64-1.4.3
    ;;
    darwin*)
    DEVARCH=darwin_arm-1.4.3
    ;;
esac

case $BOARD in
    bbb)
    TARGARCH=armv7_nerves_linux_gnueabihf
    CCEXE=armv7-nerves-linux-gnueabihf
    export TOOLCHAIN="$HOME/.nerves/artifacts/nerves_toolchain_$TARGARCH-$DEVARCH"
    export CC="$TOOLCHAIN/bin/$CCEXE-cc"
    export GOOS=linux 
    export GOARCH=arm
    export GOARM=7 
    export CGO_ENABLED=1 
    ;;
    rpi4)
    TARGARCH=aarch64_nerves_linux_gnu
    CCEXE=aarch64-nerves-linux-gnu
    export TOOLCHAIN="$HOME/.nerves/artifacts/nerves_toolchain_$TARGARCH-$DEVARCH"
    export CC="$TOOLCHAIN/bin/$CCEXE-cc"
    export GOOS=linux 
    export GOARCH=arm64
    export CGO_ENABLED=1 
    ;;
esac

MOD="github.com/YeicoLabs/laurelview"
CMD=$MOD/cmd
DST=nfw/rootfs_overlay/lvbin
mkdir -p $DST

[ ! -d $TOOLCHAIN ] && (cd nfw; MIX_TARGET=$BOARD mix deps.get)

go build -ldflags="-extld=$CC" -o $DST/lvdpm $CMD/lvdpm
go build -ldflags="-extld=$CC" -o $DST/lvnbe $CMD/lvnbe
go build -ldflags="-extld=$CC" -o $DST/lvnup $CMD/lvnup
go build -ldflags="-extld=$CC" -o $DST/lvnss $CMD/lvnss

echo "LV_NUP_ENDPOINT=127.0.0.1:80" > $DST/lvnup.env
echo "LV_DPM_ENDPOINT=127.0.0.1:81" > $DST/lvdpm.env
echo "LV_NBE_ENDPOINT=0.0.0.0:80" > $DST/lvnbe.env
echo "LV_NBE_DEBUG=127.0.0.1:82" >> $DST/lvnbe.env
echo "LV_NBE_DATABASE=/data/lvnbe.db3" >> $DST/lvnbe.env
