#!/bin/bash -x

TYPE="${1:-sd}"

case $TYPE in
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

mix deps.get
mix firmware

case $TYPE in
    sd)
    sudo fwup -aU -i _build/bbb_dev/nerves/images/nfw.fw -t complete
    sync
    ;;
    emmc)
(cd _build/bbb_emmc_prod/nerves/images/ && sftp nerves.local) << EOF
put nfw.fw /tmp/
quit
EOF

ssh nerves.local << EOF
cmd "fwup -aU -i /tmp/nfw.fw -d /dev/mmcblk1 -t complete"
cmd "poweroff"
exit
EOF
    ;;
esac
