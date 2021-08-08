# laurelview

## Code Tree

```bash
pkg/lvnbe #node backend testeable package
cmd/lvnbe #node backend executable
cmd/lvnws #node Windows service
cmd/lvcbe #cloud backend executable
web/lvnfe #node frontend
web/lvcfe #cloud frontend
```

## Mutation Flow

```bash
<-- all
<-> create
<-> update

#open coupled to url params
#close coupled to socket lifecycle
--> mode
<-- status
<-- reading
```

## Components

```
entry <- core <- hub <- state <- dao
```

## Entity Fields

- Host
- Port
- Slave (Device Address)

## v0.0.1

- Filterable list of Laurels
- Single page Index + Control Panel + Dashboard
- Control Panel is single Laurel with control buttons
- Control Buttons are Peak, Valley, Tare, Cold Reset
- Dashboard is view only multi Laurel monitor
- TCP and Serial Port
- Windows service 
- Windows installer https://github.com/mh-cbon/go-msi

## v0.0.2

- https://laurelview.io
- Link node Laurels to cloud

