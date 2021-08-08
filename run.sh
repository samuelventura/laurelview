#!/bin/sh -x

#trace|debug|info
export LV_LOGLEVEL="${1:-debug}"
export LV_ENDPOINT="${2:-:5001}"
MOD="github.com/samuelventura/laurelview"
(go install $MOD/cmd/lvnbe && lvnbe)
