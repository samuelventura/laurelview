#!/bin/bash

#install with
#crontab -r
#echo "@reboot /home/pi/bin/lvreboot" | crontab -

UH=/home/pi

echo `date` >> $UH/.reboot
echo $HOME >> $UH/.reboot
echo $USER >> $UH/.reboot #not set on reboot

export USER=pi
export PATH=$HOME/bin:$PATH

sleep 3
lvstart.sh
