#!/bin/bash -xe

SRC=pkg/sdk.go

cp -f $SRC pkg/lvnrt/
cp -f $SRC pkg/lvndb/
cp -f $SRC cmd/lvnbe/
cp -f $SRC cmd/lvnss/
cp -f $SRC cmd/lvnup/
cp -f $SRC cmd/lvdpm/
cp -f $SRC cmd/lvtry/
cp -f $SRC cmd/lvdos/

sed -i '' '1s/.*/package lvnrt/' pkg/lvnrt/sdk.go
sed -i '' '1s/.*/package lvndb/' pkg/lvndb/sdk.go
sed -i '' '1s/.*/package main/' cmd/lvnbe/sdk.go
sed -i '' '1s/.*/package main/' cmd/lvnss/sdk.go
sed -i '' '1s/.*/package main/' cmd/lvnup/sdk.go
sed -i '' '1s/.*/package main/' cmd/lvdpm/sdk.go
sed -i '' '1s/.*/package main/' cmd/lvtry/sdk.go
sed -i '' '1s/.*/package main/' cmd/lvdos/sdk.go
