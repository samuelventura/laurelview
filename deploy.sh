#!/bin/bash -x

# cloud|node
TARGET="${1:-cloud}"

#monitor upgrade in target machine
#tail -f /usr/local/bin/lvclk.err.log

case $TARGET in
    cloud)
    BIN=~/go/bin
    export GOARCH=amd64
    export GOOS=linux
    (cd cmd/lvclk && go build -o /tmp/lvclk)
    rsync /tmp/lvclk ssh.laurelview.io:/tmp/lvclk
    DAEMON=lvclk
    DST=/usr/local/bin
    ssh ssh.laurelview.io "cat > /tmp/lvclk.sh; chmod a+x /tmp/lvclk.sh" << EOF
curl -X POST http://127.0.0.1:31600/api/daemon/stop/$DAEMON
curl -X POST http://127.0.0.1:31600/api/daemon/uninstall/$DAEMON
sudo cp /tmp/$DAEMON $DST
curl -X POST http://127.0.0.1:31600/api/daemon/install/$DAEMON?path=$DST/$DAEMON
curl -X POST http://127.0.0.1:31600/api/daemon/enable/$DAEMON
curl -X POST http://127.0.0.1:31600/api/daemon/start/$DAEMON
curl -X GET http://127.0.0.1:31600/api/daemon/info/$DAEMON
curl -X GET http://127.0.0.1:31600/api/daemon/env/$DAEMON    
EOF
    ssh ssh.laurelview.io eval "/tmp/lvclk.sh > /dev/null 2>&1"
    ;;
    node)
    ;;
esac
