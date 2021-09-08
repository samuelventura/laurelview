#!/bin/bash -xe

IP="${1:-10.77.3.143}"
rsync -r cmd/lvnbe/build pi@$IP:local/laurelview/cmd/lvnbe/
rsync -r sbc/.ssh pi@$IP:
rsync -r sbc/bin pi@$IP:
ssh pi@$IP 'crontab -r'
ssh pi@$IP 'echo "@reboot /home/pi/bin/lvreboot.sh" | crontab -'
ssh pi@$IP 'dos2unix bin/*.sh .ssh/*'
ssh pi@$IP 'sudo bin/lvsetup.sh'
