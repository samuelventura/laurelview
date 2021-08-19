#!/bin/sh -x

if [[ "$OSTYPE" == "linux"* ]]; then
    SRC=$HOME/go/bin
    DST=/usr/local/bin
    sudo cp $SRC/lvnbe $DST
    sudo cp $SRC/lvnrt $DST
    sudo cp $SRC/lvnss $DST
    sudo $DST/lvnss -service install
fi
