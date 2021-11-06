#!/bin/bash -xe

SBC="${1:-bbb}"

#passwd changed to Tp4jTk7mpNwN
#export PATH=~/bin:~/go/bin:~/local/go/bin:$PATH
case $SBC in
    bbb)
    USR=debian
    NIF=eth0
    IP="${2:-10.77.3.155}"
    #requires /etc/sudoers.d/nopwd with
    #debian ALL=NOPASSWD: ALL
    ;;
    bbbw)
    USR=debian
    NIF=wlan0
    IP="${2:-10.77.3.146}"
    #requires /etc/sudoers.d/nopwd with
    #debian ALL=NOPASSWD: ALL
    ;;
    pi)
    USR=pi
    NIF=eth0
    IP="${2:-10.77.3.143}"
    ;;
    piw)
    USR=pi
    NIF=wlan0
    IP="${2:-192.168.0.23}"
    #requires /etc/sudoers.d/nopwd with
    #debian ALL=NOPASSWD: ALL
    ;;
esac

#avoid sbc node installation and slow compilation
rsync -r cmd/lvnbe/build $USR@$IP:local/laurelview/cmd/lvnbe/
rsync -r sbc/.ssh $USR@$IP:
rsync -r sbc/bin $USR@$IP:
ssh $USR@$IP "crontab -r || true"
ssh $USR@$IP "echo @reboot /home/$USR/bin/lvreboot.sh | crontab -"
ssh $USR@$IP "dos2unix bin/*.sh .ssh/*"
ssh $USR@$IP "sudo bin/lvsetup.sh"
ssh $USR@$IP "echo NIF=$NIF > bin/lvenv.sh"
ssh $USR@$IP "echo USER=$USR >> bin/lvenv.sh"
