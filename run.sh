#!/bin/bash -xe

#ps -A | grep lv
trap 'kill $(jobs -p)' SIGINT SIGTERM EXIT

#info|debug|trace
BIN=~/go/bin
export LV_LOGLEVEL="${1:-info}"
export LV_NBE_DEBUG="127.0.0.1:5000"
export LV_NBE_ENDPOINT="127.0.0.1:5001"
export LV_DPM_ENDPOINT="127.0.0.1:5002"
export LV_NUP_ENDPOINT="127.0.0.1:5001"
export LV_CBE_ENDPOINT="127.0.0.1:5003"
export LV_SBE_ENDPOINT="127.0.0.1:5004"
export LV_NSS_LOGS="/tmp"
MOD="github.com/YeicoLabs/laurelview"
mkdir -p cmd/lvnbe/build
mkdir -p cmd/lvcbe/build
mkdir -p cmd/lvsbe/build
touch cmd/lvnbe/build/empty.txt
touch cmd/lvcbe/build/empty.txt
touch cmd/lvsbe/build/empty.txt
go install $MOD/cmd/lvdpm
go install $MOD/cmd/lvnbe
go install $MOD/cmd/lvnup
go install $MOD/cmd/lvcbe
go install $MOD/cmd/lvsbe

run_lv() {
    rm -f /tmp/$1.fifo
    mkfifo /tmp/$1.fifo
    cat /tmp/$1.fifo | $BIN/$1 &
}

run_lv "lvdpm"
run_lv "lvnbe"
run_lv "lvnup"
run_lv "lvcbe"
run_lv "lvsbe"

read -p "Press ENTER to quit..."
