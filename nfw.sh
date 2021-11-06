#!/bin/bash -x

export NFW_USB=/dev/sdc1
rm -fr nfw/_build
(cd nfw && iex -S mix)
