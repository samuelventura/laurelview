#!/bin/sh -x

export TARGET="${1:-local}"

case $TARGET in
    local)
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
    ;;
    pi)
    PIHOST=10.77.3.143
    scp -r /tmp/lvpi pi@$PIHOST: && ssh pi@$PIHOST "
    [[ -f /etc/systemd/system/LaurelView.service ]] || sudo /home/pi/lvpi/lvnss -service install
    sudo systemctl restart LaurelView
    "
    ;;
esac
