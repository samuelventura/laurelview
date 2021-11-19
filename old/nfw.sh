#!/bin/bash -x

touch ../../nerves_backdoor/mix.exs
export NFW_USB=/dev/sdc1
rm -fr nfw/_build
(cd nfw && iex -S mix)
