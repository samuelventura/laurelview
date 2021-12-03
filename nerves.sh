#!/bin/bash -xe

# ssh|sd|emmc|clean
MEDIA="${1:-ssh}"
# upgrade|complete
TYPE="${2:-upgrade}"
# bbb|rpi4
BOARD="${3:-bbb}"
# $HOST|10.77.3.171
HOST="${4:-$HOST}"

case $MEDIA in
    sd|ssh)
    export MIX_ENV=dev
    export MIX_TARGET=$BOARD
    ;;
    emmc)
    export MIX_ENV=prod
    export MIX_TARGET=bbb_emmc
    ;;
    clean)
    # ensure NFW_BIN is reset
    rm -fr nfw/_build
    exit
    ;;
    deps)
    mix local.hex --force
    mix local.rebar --force
    mix archive.install hex nerves_bootstrap --force
    cd nfw
    export MIX_ENV=dev
    export MIX_TARGET=bbb
    mix deps.get
    export MIX_ENV=dev
    export MIX_TARGET=rpi4
    mix deps.get
    export MIX_ENV=prod
    export MIX_TARGET=bbb_emmc
    mix deps.get
    export -n MIX_ENV
    export -n MIX_TARGET
    exit
    ;;
esac

cd nfw

mix firmware

KEY=`pwd`/id_rsa
case $MEDIA in
    ssh) #wont work for bbb_emmc
    ssh-add $KEY
    mix upload $HOST
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
(cd $IMAGES && sftp -oIdentityFile=$KEY $HOST) << EOF
put nfw.fw /tmp/
quit
EOF

ssh -i $KEY $HOST << EOF
cmd "fwup -aU -i /tmp/nfw.fw -d /dev/mmcblk1 -t $TYPE"
cmd "poweroff"
exit
EOF
    ;;
esac
