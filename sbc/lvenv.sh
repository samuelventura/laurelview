#!/bin/bash

MAC=$1

echo $MAC >> /tmp/macs


HTTP=80
SSH=22

case $MAC in
    "WINDOWS")
    HTTPP=31601
    SSHP=31001
    HTTP=31601
    ;;
    "90:59:af:55:b0:aa")
    #bbb
    HTTPP=31601
    SSHP=31001
    ;;
    "80:30:dc:50:71:ea")
    #bbbw
    HTTPP=31601
    SSHP=31001
    ;;
    "b8:27:eb:03:61:49")
    #pi
    HTTPP=31601
    SSHP=31001
    ;;
    "b8:27:eb:90:02:8d")
    #pi with touch de Revilla wlan0
    HTTPP=31601
    SSHP=31001
    ;;
    "b8:27:eb:c5:57:d8")
    #pi with touch de Revilla eth0
    HTTPP=31601
    SSHP=31001
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
echo "LVHTTP=$HTTPP:127.0.0.1:$HTTP"
echo "LVSSH=$SSHP:127.0.0.1:$SSH"
