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
    OUT=/tmp/pibuild
    mkdir -p $OUT
    FLAGS="GOOS=linux GOARCH=arm GOARM=7"
    env $FLAGS go build -o $OUT/lvdpm $MOD/cmd/lvdpm
    env $FLAGS go build -o $OUT/lvnbe $MOD/cmd/lvnbe
    env $FLAGS go build -o $OUT/lvnup $MOD/cmd/lvnup
    env $FLAGS go build -o $OUT/lvnss $MOD/cmd/lvnss
    ;;
esac
