#!/bin/bash -xe

source lvenv.sh
LVMAC=`cat /sys/class/net/$NIF/address`
echo $LVMAC > /tmp/lvmac
ssh samuel@ssh.laurelview.io bin/lvenv.sh $LVMAC > /tmp/lvenv
source /tmp/lvenv
ssh -N -R $LVHTTP -R $LVSSH $LVLOGIN 
