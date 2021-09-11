#!/bin/bash

MAC=$1

echo $MAC >> /tmp/macs

case $MAC in
    "90:59:af:55:b0:aa")
    #bbb
    HTTPP=31602
    SSHP=31002
    ;;
    "80:30:dc:50:71:ea")
    #bbbw
    HTTPP=31601
    SSHP=31001
    ;;
    "b8:27:eb:03:61:49")
    #pi
    HTTPP=31603
    SSHP=31003
    ;;
esac

function killp() {
    #nc -l makes ssh tunnel listen only on ipv6 and netstat lists 2 different pids
    PORT=$1
    PIDS=`sudo netstat -lnp | grep ":$PORT" | gawk 'match($0,/LISTEN\s+([0-9]+)\//,a){print a[1]}' | uniq | tr '\n' ' '`
    [ -z "$PIDS" ] || sudo kill -9 $PIDS
}

killp $HTTPP >/dev/null 2>&1
killp $SSHP >/dev/null 2>&1

echo "LVLOGIN=samuel@vpn.yeico.com"
echo "LVHTTP=$HTTPP:127.0.0.1:80"
echo "LVSSH=$SSHP:127.0.0.1:22"
