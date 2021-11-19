#!/bin/bash -x

MEDIA="${1:-sd}"
#complete|upgrade
TYPE="${2:-complete}"

case $MEDIA in
    sd)
    export MIX_ENV=dev
    export MIX_TARGET=bbb
    ;;
    emmc)
    export MIX_ENV=prod
    export MIX_TARGET=bbb_emmc
    ;;
esac

#ensure NFW_BIN is reset
rm -fr nfw/_build

cd nfw

#first time requires 
#mix archive.install hex nerves_bootstrap
mix deps.get
mix firmware

KEY=`pwd`/id_rsa
case $MEDIA in
    sd)
    IMAGES=_build/bbb_dev/nerves
    sudo fwup -aU -i $IMAGES/images/nfw.fw -t $TYPE
    sync
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
