#!/bin/bash -xe

SBC="${1:-bbbw}"

#passwd changed to Tp4jTk7mpNwN
#export PATH=~/bin:~/go/bin:~/local/go/bin:$PATH
case $SBC in
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
esac

rsync -r cmd/lvnbe/build $USR@$IP:local/laurelview/cmd/lvnbe/
rsync -r sbc/.ssh $USR@$IP:
rsync -r sbc/bin $USR@$IP:
ssh $USR@$IP "crontab -r || true"
ssh $USR@$IP "echo @reboot /home/$USR/bin/lvreboot.sh | crontab -"
ssh $USR@$IP "dos2unix bin/*.sh .ssh/*"
ssh $USR@$IP "sudo bin/lvsetup.sh"
ssh $USR@$IP "echo NIF=$NIF > bin/lvenv.sh"
ssh $USR@$IP "echo USER=$USR >> bin/lvenv.sh"
