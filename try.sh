#!/bin/sh -x

MOD="github.com/samuelventura/laurelview"
go install $MOD/cmd/lvtry && lvtry
