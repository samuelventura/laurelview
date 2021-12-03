#!/bin/bash -xe

function tunnel() {
    LVMAC="WINDOWS"
    (ssh samuel@ssh.laurelview.io bin/lvenv.sh $LVMAC > /tmp/lvenv) && \
    source /tmp/lvenv && \
    ssh -N -R $LVHTTP -R $LVSSH $LVLOGIN 
}

while true
do
	tunnel
	sleep 2
done
