#!/bin/sh -x

if [[ "$OSTYPE" == "linux"* ]]; then
    SRC=$HOME/go/bin
    DST=/usr/local/bin
    if [[ -f "$DST/lvnss" ]]; then
        sudo systemctl stop LaurelView
        sudo $DST/lvnss -service uninstall
        sleep 3
    fi
    sudo cp $SRC/lvdpm $DST
    sudo cp $SRC/lvnbe $DST
    sudo cp $SRC/lvnup $DST
    sudo cp $SRC/lvnss $DST
    sudo $DST/lvnss -service install
    sudo systemctl restart LaurelView
fi
