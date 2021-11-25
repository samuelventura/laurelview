#!/bin/bash -x

mix local.hex --force
mix local.rebar --force
mix archive.install hex nerves_bootstrap --force

cd nfw

export MIX_ENV=dev
export MIX_TARGET=bbb
mix deps.get

export MIX_ENV=dev
export MIX_TARGET=rpi4
mix deps.get

export MIX_ENV=prod
export MIX_TARGET=bbb_emmc
mix deps.get

export -n MIX_ENV
export -n MIX_TARGET
