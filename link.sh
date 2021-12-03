#!/bin/bash -xe

#ps -A | grep go-
trap 'kill $(jobs -p)' SIGINT SIGTERM EXIT

#beaglebone debian
#TEP=10.77.3.155:80
TEP=127.0.0.1:5001 #must include port 
BIN=~/go/bin
(cd cmd/lvnlk && go install)
(cd cmd/lvclk && go install)

run_daemon() {
    rm -f /tmp/$1.fifo
    mkfifo /tmp/$1.fifo
    cat /tmp/$1.fifo | $BIN/$1 &
}

export LV_CLK_ENDPOINT_PROXY=127.0.0.1:31080
export LV_CLK_ENDPOINT_TARGET=$TEP
run_daemon "lvclk"

sleep 1

export LV_NLK_DOCK_POOL=127.0.0.1:31622
run_daemon "lvnlk"

#./pack.sh
#./run.sh
#http://127.0.0.1:31080/
read -p "Press ENTER to quit..."
