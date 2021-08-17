#!/bin/sh -x

(cd web/lvnfe; npm run build)
rm -fr pkg/lvnbe/build/*
mkdir -p pkg/lvnbe/build
touch pkg/lvnbe/build/.empty
cp -fr web/lvnfe/build/* pkg/lvnbe/build/
MOD="github.com/samuelventura/laurelview"
go install $MOD/cmd/lvdpm
go install $MOD/cmd/lvnbe
go install $MOD/cmd/lvnrt
go install $MOD/cmd/lvnsd
