#!/bin/bash -xe

LVMAC=`cat /sys/class/net/eth0/address`
echo $LVMAC > /tmp/lvmac
ssh samuel@ssh.laurelview.io bin/lvenv $LVMAC > /tmp/lvenv
. /tmp/lvenv
ssh -N -R $LVHTTP -R $LVSSH $LVLOGIN 
