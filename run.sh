#!/bin/sh -x

#trace|debug|info
export LV_LOGLEVEL="${1:-debug}"
export LV_NBE_ENDPOINT="${2:-:5001}"
export LV_NRT_ENDPOINT="${2:-:5002}"
MOD="github.com/samuelventura/laurelview"
(go install $MOD/cmd/lvnbe && lvnbe &)
(go install $MOD/cmd/lvnrt && lvnrt &)
wait
