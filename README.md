# Laurel View

## Code Tree

```bash
pkg/lvnrt #node runtime testeable package
pkg/lvnbe #node backend testeable package
cmd/lvnrt #node runtime executable
cmd/lvnbe #node backend executable
cmd/lvnss #node system service
cmd/lvdpm #node dpm echo server
web/lvnfe #node react frontend
```

## v0.0.1

- Filterable list of Laurels
- Single page Index + Control Panel
- Control Panel is single Laurel with control buttons
- Value, Peak/Reset, Valley/Reset, Tare/Reset, Cold Reset
- Cross platform service https://github.com/kardianos/service
- Windows installer https://github.com/mh-cbon/go-msi
- TCP with slave support
- Multi Laurel view
- Banner and favicon
- Fit cell screen size

## v0.0.2

- Latency and error in display
- Dial error to all displays
- Bus pull from feed channel
- Select all/clear buttons
- Disabled multiview button
- Multiview fluid grid layout
- Basic progressive PWA
- Toggle fullscreen
- Headers valign
- Round icon for Android https://maskable.app/editor
- Square icon for iOS https://maskable.app/editor
- Resilient view modals
- Golang fast conn drop detection
- Latest Laurel node FW
- Dark mode autoswitch
- Globe chrome favicon

## Future

- Branding
- Serial Port
- Realtime Plot 
- HTTPS certificate
- Progressive PWA
- MacOS installer
- Single entry port
- Kill service test
- https://laurelview.io
- First dpm duplicate update (only in low latency?)
- Link node Laurels to cloud
- iPhone bookmark shows react icon
- Executable version https://github.com/josephspurrier/goversioninfo

# Networking Issues

- Golang takes ~10 to detect connection drop if panel powered of
- Cold reset in second daisy chain device (transient errors ~1m50s)
- Tara/Valley/Peak reset (transiente <400ms)
- Second connection attempt makes next connection take ~20s
- Connection drop
