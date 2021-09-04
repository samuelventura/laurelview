#!/bin/sh -x

#trace|debug|info
export LV_LOGLEVEL="${1:-info}"
export LV_NUP_ENDPOINT="${2:-127.0.0.1:5001}"
export LV_DOS_LL=5
export LV_DOS_UL=10
export LV_DOS_FM=5000
export LV_DOS_VM=5000
MOD="github.com/samuelventura/laurelview"
go install $MOD/cmd/lvdos && lvdos
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT
#ssh -R 5001:127.0.0.1:5001 -L 55001:127.0.0.1:5001 165.22.39.156
#./dos.sh trace 127.0.0.1:55001
#LV_NUP_ENDPOINT=127.0.0.1:55001 LV_LOGLEVEL=trace lvnup
