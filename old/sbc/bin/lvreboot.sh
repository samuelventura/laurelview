#!/bin/bash

source lvenv.sh
echo `date` >> $HOME/.reboot
export PATH=$HOME/bin:$PATH

sleep 3
lvstart.sh
