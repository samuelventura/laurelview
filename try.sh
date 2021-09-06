#!/bin/bash -x

MOD="github.com/samuelventura/laurelview"
go install $MOD/cmd/lvtry && lvtry
