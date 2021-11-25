#!/bin/bash -x

# ssh|sd|emmc
MEDIA="${1:-ssh}"
# upgrade|complete
TYPE="${2:-upgrade}"
# bbb|rpi4
BOARD="${3:-bbb}"

case $MEDIA in
    sd|ssh)
    export MIX_ENV=dev
    export MIX_TARGET=$BOARD
    ;;
    emmc)
    export MIX_ENV=prod
    export MIX_TARGET=bbb_emmc
    ;;
esac

# ensure NFW_BIN is reset
rm -fr nfw/_build

cd nfw

# first time requires 
# mix archive.install hex nerves_bootstrap
mix deps.get
mix firmware

KEY=`pwd`/id_rsa
case $MEDIA in
    ssh) #wont work for bbb_emmc
    ssh-add $KEY
    mix upload nerves.local
    ;;
    sd)
    case $TYPE in
        complete)
        mix firmware.burn
        sync
        ;;
        upgrade)
        IMAGES=_build/${BOARD}_dev/nerves
        sudo fwup -aU -i $IMAGES/images/nfw.fw -t $TYPE        
        sync
        ;;
    esac
    ;;
    emmc)
    IMAGES=_build/bbb_emmc_prod/nerves/images
(cd $IMAGES && sftp -oIdentityFile=$KEY nerves.local) << EOF
put nfw.fw /tmp/
quit
EOF

ssh -i $KEY nerves.local << EOF
cmd "fwup -aU -i /tmp/nfw.fw -d /dev/mmcblk1 -t $TYPE"
cmd "poweroff"
exit
EOF
    ;;
esac
