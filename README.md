# Laurel View

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
./node.sh 
./run.sh info|debug|trace
#sbc install
./pack.sh #node build
./sbc.sh bbb|bbbw|pi|piw
./build.sh #from sbc
./install.sh #from sbc
#windows
./pack.sh
./build.sh
./inno.sh #gui
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