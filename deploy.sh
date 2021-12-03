#!/bin/bash -xe

# cloud|node
TARGET="${1:-cloud}"

#monitor upgrade in target machine
#tail -f /usr/local/bin/lvclk.err.log
#tail -f /usr/local/bin/lvnlk.err.log

deploy() {
    DAEMON=$1
    HOST=$2
    BIN=~/go/bin
    export GOARCH=amd64
    export GOOS=linux
    (cd cmd/$DAEMON && go build -o /tmp/$DAEMON)
    rsync /tmp/$DAEMON $HOST:/tmp/$DAEMON
    DST=/usr/local/bin
    ssh $HOST "cat > /tmp/$DAEMON.sh; chmod a+x /tmp/$DAEMON.sh" << EOF
curl -X POST http://127.0.0.1:31600/api/daemon/stop/$DAEMON
curl -X POST http://127.0.0.1:31600/api/daemon/uninstall/$DAEMON
sudo cp /tmp/$DAEMON $DST
curl -X POST http://127.0.0.1:31600/api/daemon/install/$DAEMON?path=$DST/$DAEMON
curl -X POST http://127.0.0.1:31600/api/daemon/enable/$DAEMON
curl -X POST http://127.0.0.1:31600/api/daemon/start/$DAEMON
curl -X GET http://127.0.0.1:31600/api/daemon/info/$DAEMON
curl -X GET http://127.0.0.1:31600/api/daemon/env/$DAEMON    
EOF
    ssh $HOST eval "/tmp/$DAEMON.sh > /dev/null 2>&1"
}

case $TARGET in
    cloud)
    deploy lvclk ssh.laurelview.io
    ;;
    node)
    deploy lvnlk 10.77.0.49
    ;;
esac
