#!/bin/sh -x

#trace|debug|info
TEST_SCOPE="${1:-Test}"
export LV_LOGLEVEL="${2:-info}"
MOD="github.com/samuelventura/laurelview"
go clean -testcache 
go test $MOD/pkg/lvndb -v -run $TEST_SCOPE
