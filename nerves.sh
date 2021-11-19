#!/bin/bash -x

#sd|emmc|sdssh|emmcssh
MEDIA="${1:-sd}"
#upgrade|complete
TYPE="${2:-upgrade}"

case $MEDIA in
    sd*)
    export MIX_ENV=dev
    export MIX_TARGET=bbb
    ;;
    emmc*)
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
    *ssh)
    #sd upgraded ok, emmc upgrade marks image as NOT VALIDATED
    #nfw 0.1.0 (4761f97d-5d54-5f73-9f37-3a8fee26fecb) arm bbb
    #nfw 0.1.0 (8654fa33-8f2f-5046-cbf9-fa7064a9bc75) arm bbb
    ssh-add $KEY
    mix upload nerves.local
    ;;
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
