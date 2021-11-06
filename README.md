# Laurel View

## Dev Environ

```bash
#linux
sudo snap install go #go1.17.2 linux/amd64
sudo snap install node #v16.13.0
#https://elixir-lang.org/install.html#gnulinux 
#based on erlang-solutions release
#don't install elixir from here
#nerves needs elixir with same otp version
apt install esl-erlang #1:24.1.3-1 amd64
#https://github.com/taylor/kiex
kiex install 1.12.3
kiex default 1.12.3
#mix archive.install hex nerves_bootstrap

#macos

```

## Code Tree

```bash
pkg/lvnrt #node runtime testeable package
pkg/lvnbe #node backend testeable package
cmd/lvnrt #node runtime testeable package
cmd/lvnbe #node backend executable
cmd/lvnup #node uplink executable
cmd/lvnss #node system service
cmd/lvdpm #node dpm echo server
web/lvnfe #node react frontend
```

## Howto

```bash
#testing
./test.sh all|db|rt|sdk
#developing 
./node.sh       #launches http://localhost:3000/
./run.sh info|debug|trace
./build.sh      #for elixir
./nss.sh        #launch iex
#BBB firmware
./pack.sh       #node build
./cross.sh      #build and zip
./nerves.sh sd|emmc
#windows installer
./pack.sh
./build.sh
./inno.sh       #gui
```

## Helpers

```bash
ssh nerves.local -i nfw/id_rsa
#elixir development
iex -S mix
recompile
Application.start :nss
Application.stop :nss
#elixir environment info
Application.started_applications
Application.get_all_env :nfw
Application.get_all_env :nss
Application.app_dir :nss, "priv"
#network configuration 10.77.3.167
ifconfig
VintageNet.info
ping "10.77.0.49"
ping "google.com"
VintageNet.configure("eth0", %{type: VintageNetEthernet, ipv4: %{method: :dhcp}})
VintageNet.configure("eth0", %{type: VintageNetEthernet, ipv4: %{ method: :static, address: "10.77.4.165", prefix_length: 8, name_servers: []}})
VintageNet.configure("eth0", %{type: VintageNetEthernet, ipv4: %{ method: :static, address: "10.77.4.165", prefix_length: 8, gateway: "10.77.0.1", name_servers: ["10.77.0.1"]}})
#reboot clear first boot errors
ssh nerves.local << EOF
cmd "reboot"
exit
EOF
#get priv data path
iex(1)> :code.priv_dir(:nss)           
'/srv/erlang/lib/nss-0.1.0/priv'
#list priva data folder
iex(2)> cmd "ls /srv/erlang/lib/nss-0.1.0/priv"
test.txt
0
```

## Future

- Serial Port
- Realtime Plot 
- HTTPS certificate
- Kill service test
- Autoscale view screens
- Link node Laurels to cloud
- iPhone bookmark shows react icon

# Networking Issues

- Golang takes ~10s to detect connection drop if panel powered off
- Cold reset in second daisy chain device (transient errors ~1m50s)
- Second connection attempt makes next connection take ~20s
- Tara/Valley/Peak reset (transiente <400ms)

# Audits

- https://web.dev/measure/
- https://web.dev/vitals/
- https://www.pwabuilder.com/

# Resources

- https://maskable.app/editor
- https://realfavicongenerator.net/
- https://caniuse.com/?search=a2hs
- https://jrsoftware.org/isinfo.php
- https://favicon.io/favicon-converter/
- https://github.com/audreyfeldroy/favicon-cheat-sheet