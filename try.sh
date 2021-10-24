#!/bin/bash -x

MOD="github.com/YeicoLabs/laurelview"
go install $MOD/cmd/lvtry && lvtry
