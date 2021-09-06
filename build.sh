#!/bin/sh -x

export TARGET="${1:-local}"

MOD="github.com/samuelventura/laurelview"

case $TARGET in
    local)
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
    ;;
    pi)
    OUT=/tmp/lvpi
    mkdir -p $OUT
    FLAGS="GOOS=linux GOARCH=arm GOARM=7"
    env $FLAGS go build -o $OUT/lvdpm $MOD/cmd/lvdpm
    env $FLAGS go build -o $OUT/lvnbe $MOD/cmd/lvnbe
    env $FLAGS go build -o $OUT/lvnup $MOD/cmd/lvnup
    env $FLAGS go build -o $OUT/lvnss $MOD/cmd/lvnss
    ;;
    bb)
    OUT=/tmp/lvbb
    mkdir -p $OUT
    CC="/C/SysGCC/beaglebone/bin/arm-linux-gnueabihf-gcc.exe"
    LD="/C/SysGCC/beaglebone/bin/arm-linux-gnueabihf-ld.exe"
    SR="/C/SysGCC/beaglebone/arm-linux-gnueabihf/sysroot"
    CGO_CFLAGS="--sysroot=$SR"
    CGO_LDFLAGS="--sysroot=$SR -m=armelf_linux_eabi"
    FLAGS0="GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 CC=$CC CGO_CFLAGS=$CGO_CFLAGS CGO_LDFLAGS=$CGO_LDFLAGS"
    FLAGS1="GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=$CC CGO_CFLAGS=$CGO_CFLAGS CGO_LDFLAGS=$CGO_LDFLAGS"
    env $FLAGS1 go build -x -o $OUT/lvnbe $MOD/cmd/lvnbe
    # env $FLAGS0 go build -o $OUT/lvdpm $MOD/cmd/lvdpm
    # env $FLAGS0 go build -o $OUT/lvnup $MOD/cmd/lvnup
    # env $FLAGS0 go build -o $OUT/lvnss $MOD/cmd/lvnss
    ;;
esac
