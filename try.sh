#!/bin/bash -x

export PATH=~/go/bin:$PATH
MOD="github.com/YeicoLabs/laurelview"
go install $MOD/cmd/lvtry && lvtry
