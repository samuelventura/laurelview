#!/bin/sh -x

#trace|debug|info
TEST_PKG="${1:-all}"
TEST_SCOPE="${2:-Test}"
export LV_LOGLEVEL="${3:-info}"
MOD="github.com/samuelventura/laurelview"
go clean -testcache 

case $TEST_PKG in
    all)
    go test $MOD/pkg/lvsdk -v -run $TEST_SCOPE
    go test $MOD/pkg/lvndb -v -run $TEST_SCOPE
    go test $MOD/pkg/lvnrt -v -run $TEST_SCOPE
    ;;
    db)
    go test $MOD/pkg/lvndb -v -run $TEST_SCOPE
    ;;
    rt)
    go test $MOD/pkg/lvnrt -v -run $TEST_SCOPE
    ;;
    sdk)
    go test $MOD/pkg/lvsdk -v -run $TEST_SCOPE
    ;;
esac
