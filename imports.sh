#!/bin/sh -x
SRC=pkg/sdk.go 
cp -f $SRC pkg/lvnrt/
cp -f $SRC pkg/lvndb/
cp -f $SRC cmd/lvnbe/
cp -f $SRC cmd/lvnss/
cp -f $SRC cmd/lvdpm/
cp -f $SRC cmd/lvtry/

sed -i '' '1s/.*/package lvnrt/' pkg/lvnrt/sdk.go
sed -i '' '1s/.*/package lvndb/' pkg/lvndb/sdk.go
sed -i '' '1s/.*/package main/' cmd/lvnbe/sdk.go
sed -i '' '1s/.*/package main/' cmd/lvnss/sdk.go
sed -i '' '1s/.*/package main/' cmd/lvdpm/sdk.go
sed -i '' '1s/.*/package main/' cmd/lvtry/sdk.go
