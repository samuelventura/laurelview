#!/bin/sh -x

#trace|debug|info
export LV_LOGLEVEL="${1:-info}"
export LV_NBE_ENDPOINT=":5001"
export LV_DPM_ENDPOINT=":5002"
MOD="github.com/samuelventura/laurelview"
mkdir -p cmd/lvnbe/build
touch cmd/lvnbe/build/.empty
go install $MOD/cmd/lvdpm
go install $MOD/cmd/lvnbe
go install $MOD/cmd/lvnup
go install $MOD/cmd/lvnss && lvnss > /tmp/lvnss.log
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT
