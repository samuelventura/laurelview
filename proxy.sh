#!/bin/bash -x

#ps -A | grep go-
trap 'kill $(jobs -p)' SIGINT SIGTERM EXIT

#beaglebone debian
#TEP=10.77.3.155:80
#./pack.sh
#./run.sh
TEP=127.0.0.1:5001 
BIN=~/go/bin
SRC=~/github
cd $SRC/go-ship-ms && go install
cd $SRC/go-dock-ms && go install
cd $SRC/go-proxy-ms && go install

run_goms() {
    rm -f $BIN/go-$1-ms.db3
    rm -f /tmp/go-$1-ms.fifo
    mkfifo /tmp/go-$1-ms.fifo
    cat /tmp/go-$1-ms.fifo | $BIN/go-$1-ms &
}

export PROXY_SERVER_CRT=$SRC/go-proxy-ms/server.crt
export PROXY_SERVER_KEY=$SRC/go-proxy-ms/server.key
export PROXY_HTTPS_EP=127.0.0.1:31080
export PROXY_HTTP_EP=127.0.0.1:31081
export PROXY_DOCK_EP=127.0.0.1:31023
export PROXY_API_EP=127.0.0.1:31088
export PROXY_MAIN_URL=https://www.google.com
export PROXY_HOSTNAME=demo
run_goms "proxy"

export DOCK_ENDPOINT_SSH=127.0.0.1:31022
export DOCK_ENDPOINT_API=127.0.0.1:31023
export DOCK_HOSTKEY=$SRC/go-dock-ms/id_rsa.key
run_goms "dock"

sleep 1
#register to proxy
curl -X POST http://127.0.0.1:31088/api/ship/add/demo?prefix=http://$TEP
curl -X POST http://127.0.0.1:31088/api/ship/enable/demo
#register to dock
curl -X POST http://127.0.0.1:31023/api/key/add/default -F "file=@$SRC/go-dock-ms/id_rsa.pub"
curl -X POST http://127.0.0.1:31023/api/key/enable/default
curl -X POST http://127.0.0.1:31023/api/ship/add/demo
curl -X POST http://127.0.0.1:31023/api/ship/enable/demo

export SHIP_NAME=demo
export SHIP_DOCK_KEYPATH=$SRC/go-ship-ms/id_rsa.key
export SHIP_DOCK_POOL=127.0.0.1:31022
run_goms "ship"

#https://127.0.0.1:31080/proxy/demo
read -p "Press any key..."
