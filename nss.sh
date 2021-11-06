#!/bin/bash -x

export NSS_FOLDER=~/go/bin
rm -fr nss/_build
(cd nss && iex -S mix)
