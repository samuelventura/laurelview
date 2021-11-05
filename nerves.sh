#!/bin/bash -x

cd nfw

TYPE="${1:-sd}"

case $TYPE in
    sd)
    export MIX_ENV=dev
    export MIX_TARGET=bbb
    ;;
    emmc)
    export MIX_ENV=dev
    export MIX_TARGET=bbb_emmc
    ;;
esac

mix deps.get
mix firmware

case $TYPE in
    sd)
    sudo fwup -aU -i _build/bbb_dev/nerves/images/nfw.fw -t complete
    sync
    ;;
esac
