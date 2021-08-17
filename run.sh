#!/bin/sh -x

#trace|debug|info
export LV_LOGLEVEL="${1:-debug}"
export LV_NBE_ENDPOINT="${2:-:5001}"
export LV_NRT_ENDPOINT="${2:-:5002}"
export LV_DPM_ENDPOINT="${2:-:5003}"
MOD="github.com/samuelventura/laurelview"
mkdir -p pkg/lvnbe/build
touch pkg/lvnbe/build/.empty
mkfifo /tmp/lvdpm.fifo #keep lvdpm stdin open
go install $MOD/cmd/lvdpm && (tail -f /tmp/lvdpm.fifo | lvdpm > /tmp/lvdpm.log) &
go install $MOD/cmd/lvnbe
go install $MOD/cmd/lvnrt
go install $MOD/cmd/lvnss && lvnss > /tmp/lvnss.log
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT
wait
