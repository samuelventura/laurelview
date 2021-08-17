#!/bin/sh -x

MOD="github.com/samuelventura/laurelview"
#go get github.com/akavel/rsrc
rsrc -ico icon.ico -o build/rsrc.syso
cp build/rsrc.syso cmd/lvnss/
cp build/rsrc.syso cmd/lvnbe/
cp build/rsrc.syso cmd/lvnrt/
go install $MOD/cmd/lvdpm
go install $MOD/cmd/lvnbe
go install $MOD/cmd/lvnrt
go install $MOD/cmd/lvnss
