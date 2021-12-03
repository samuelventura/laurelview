#!/bin/bash -xe

MOD="github.com/YeicoLabs/laurelview"
go install $MOD/cmd/lvndc
~/go/bin/lvndc
