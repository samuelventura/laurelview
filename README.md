# Laurel View

- [Photo Album](https://photos.app.goo.gl/94ASC7XjtpJcb2LF6)
- [Overview Slides](https://docs.google.com/presentation/d/1xeui3pKBsiawwA66xwWQhE_MmHkpiWV1BbGwXqUr4k0/edit?usp=sharing)
- [Live Demo https://laurelview.io/](https://laurelview.io/)

## What is LaurelView?

- A Web Viewer for [Laurel Electronics devices](https://www.laurels.com/)
- A technological HW/SW exploration
- An arquitectural exploration

## Tech Stack

- BBB SBC
- Nerves firmware
- [Elixir + cowboy + plug setup layer](https://github.com/samuelventura/nerves_backdoor)
- Golang device driver and multiplexer
- ReactJS UIs

## Code Tree

```bash
pkg/lvnrt #node runtime testeable package
pkg/lvnbe #node backend testeable package
cmd/lvnrt #node runtime testeable package
cmd/lvnbe #node backend executable
cmd/lvnup #node up checker executable
cmd/lvnss #node system service
cmd/lvnlk #node uplink ssh client
cmd/lvdpm #node dpm echo server
web/lvnfe #node react frontend
cmd/lvcbe #cloud backend executable
web/lvcfe #cloud react frontend
web/lvclk #cloud link sshd server
cmd/lvsbe #setup backend executable
web/lvsfe #setup react frontend
cmd/lvfix #fixture power loss tester
nfw #nerves firmware folder
msi #batch files to complement inno setup
```

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

## Howto

```bash
#testing go code
./test.sh all|db|rt|sdk
#develop go + react
./node.sh node|cloud #launches http://localhost:3000/
./run.sh info|debug|trace
#develop local proxy
./pack.sh node|cloud #only once
./run.sh        #go to http://localhost:5001/ 5001,5002,5003
./proxy.sh      #go to https://127.0.0.1:31080/proxy/demo/
#nerves BBB firmware
./nerves.sh clean   #once if required
./nerves.sh deps    #once only
./pack.sh node      #node build
./cross.sh bbb|rpi4 #build and zip
./nerves.sh ssh|sd|emmc upgrade|complete bbb|rpi4 nerves.local|10.77.3.171
#nerves.local conflicts between usb and visible eth0
#./nerves.sh ssh upgrade rpi4 10.77.3.171       #rpi4 samuel
#./nerves.sh ssh upgrade bbb 10.77.3.170        #bbb hiram
#./nerves.sh ssh upgrade bbb 10.77.3.182        #bbb samuel
#./nerves.sh ssh upgrade bbb 172.31.219.181     #bbb.usb samuel
#./nerves.sh emmc complete bbb 10.77.3.155      #bbb samuel
#./nerves.sh sd complete bbb
#windows installer
./pack.sh
./build.sh
./inno.sh       #gui
#nerves discovery
./discover.sh
```

## Helpers

```bash
dig dock.laurelview.io TXT
(cd cmd/lvnlk; go install && ~/go/bin/lvnlk)
(cd cmd/lvclk; go install && ~/go/bin/lvclk)
#https://github.com/samuelventura/nerves_backdoor
#nerves_backdoor was integrated into laurelview
ssh nerves.local -i nfw/id_rsa
export MIX_ENV=dev
export MIX_TARGET=bbb|rpi4
(cd nfw; ssh-add `pwd`/id_rsa)
(cd nfw; mix firmware)
(cd nfw; mix upload nerves.local)
(cd nfw; mix do firmware, upload nerves.local)
#elixir development
iex -S mix
recompile
Application.start :nfw
Application.stop :nfw
Application.start :nfw
Application.ensure_all_started :nfw
#elixir environment info
Application.started_applications
Application.loaded_applications
Application.get_all_env :nfw
Application.app_dir :nfw, "priv"
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
iex(1)> :code.priv_dir(:nfw)           
'/srv/erlang/lib/nfw-0.1.0/priv'
#list priva data folder
iex(2)> cmd "ls /srv/erlang/lib/nfw-0.1.0/priv"
test.txt
0
```

## Future

- Serial Port
- Realtime Plot 
- Authentication
- HTTPS certificate
- Kill service test
- Autoscale view screens
- Historic data reporting
- Link node Laurels to cloud
- iPhone bookmark shows react icon
- Node react proxy wont server ws to chrome

## Challenges

- Seemsless experience between intranet and roaming users
- Long term reliability of embedded data store
- Out of the box setup experience
- Remote OTA upgrades

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
