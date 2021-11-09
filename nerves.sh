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

#ensure NSS_FOLDER is reset
rm -fr nss/_build
rm -fr nfw/_build

cd nfw

#first time requires 
#mix archive.install hex nerves_bootstrap
mix deps.get
mix firmware

case $MEDIA in
    sd)
    IMAGENS=_build/bbb_dev/nerves
    sudo fwup -aU -i $IMAGENS/images/nfw.fw -t $TYPE
    sync
    ;;
    emmc)
    IMAGENS=_build/bbb_emmc_prod/nerves/images
(cd $IMAGENS && sftp nerves.local -i id_rsa) << EOF
put nfw.fw /tmp/
quit
EOF

ssh nerves.local -i id_rsa << EOF
cmd "fwup -aU -i /tmp/nfw.fw -d /dev/mmcblk1 -t $TYPE"
cmd "poweroff"
exit
EOF
    ;;
esac
