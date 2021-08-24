#!/bin/sh -x

#trace|debug|info
export LV_LOGLEVEL="${2:-debug}"
MOD="github.com/samuelventura/laurelview"
go clean -testcache 
go test $MOD/pkg/lvsdk -v -run $1
