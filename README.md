# Laurel View

## Dev Environ

```bash
#linux
sudo snap install go #go1.17.2 linux/amd64
sudo snap install node #v16.13.0
#https://elixir-lang.org/install.html#gnulinux 
#based on erlang-solutions release
#get and install the deb
#don't install elixir from here
#nerves needs elixir with same otp version
sudo apt update
sudo apt install esl-erlang #1:24.1.3-1 amd64
#https://github.com/taylor/kiex
kiex install 1.12.3
kiex use 1.12.3
#mix archive.install hex nerves_bootstrap
#https://github.com/fwup-home/fwup
#get and install the deb
sudo apt install tio curl wget tmux screen vim
sudo apt install build-essential automake autoconf git squashfs-tools ssh-askpass pkg-config curl libssl-dev libncurses5-dev bc m4 unzip cmake python
sudo gpasswd -a $USER dialout
stty -F /dev/ttyUSB0 115200
tio /dev/ttyUSB0
#Use FTDI TTL-234X-3V3
#J1 UART pinout
#1 GND
#4 RX
#5 TX 

#macos
#https://hexdocs.pm/nerves/installation.html
brew install erlang elixir node go sqlite
brew install fwup squashfs coreutils xz pkg-config
brew deps --tree --installed
brew list
#samuel@svm-mbair ~ % go version  
#go version go1.17.2 darwin/arm64
#samuel@svm-mbair ~ % node --version
#v16.4.1
#samuel@svm-mbair ~ % elixir --version
#Erlang/OTP 24 [erts-12.1.4] [source] [64-bit] [smp:8:8] [ds:8:8:10] [async-threads:1] [dtrace]
#Elixir 1.12.3 (compiled with Erlang/OTP 24)

#both
sudo visudo     #passwordless sudo
#%admin          ALL = (ALL) NOPASSWD: ALL
#sudo gpasswd -a $USER admin
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
#testing go code
./test.sh all|db|rt|sdk
#develop go + react
./node.sh       #launches http://localhost:3000/
./run.sh info|debug|trace
./build.sh      #for elixir
./nss.sh        #launch iex
#develop local proxy
./pack.sh       #only once
./run.sh        #go to http://localhost:5001/
./proxy.sh      #go to https://127.0.0.1:31080/proxy/demo/
#nerves BBB firmware
./deps.sh       #once only
./pack.sh       #node build
./cross.sh      #build and zip
./nerves.sh sd|emmc complete|upgrade
#windows installer
./pack.sh
./build.sh
./inno.sh       #gui
```

## Helpers

```bash
#https://github.com/samuelventura/nerves_backdoor
ssh nerves.local -i nfw/id_rsa
#elixir development
iex -S mix
recompile
Application.start :nss
Application.stop :nss
#elixir environment info
Application.started_applications
Application.loaded_applications
Application.get_all_env :nfw
Application.get_all_env :nss
Application.app_dir :nss, "priv"
#network configuration 10.77.3.167
ifconfig
VintageNet.info
ping "10.77.0.49"
ping "google.com"
VintageNet.get_configuration("eth0")
VintageNet.get(["interface", "eth0", "type"])
VintageNet.get(["interface", "eth0", "state"])
VintageNet.get(["interface", "eth0", "connection"])
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
- Node react proxy wont server ws to chrome

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