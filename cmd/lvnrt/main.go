package main

import (
	"os"
	"os/signal"
	"strings"

	"github.com/samuelventura/laurelview/pkg/lvnrt"
)

func main() {
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)
	rt := lvnrt.DefaultRuntime()
	log := rt.PrefixLog("main")
	defer log.Log("") //wait flush
	defer rt.TraceRecover()
	rt.SetdAsync("hub", lvnrt.NewHub(rt))
	rt.SetdAsync("state", lvnrt.NewState(rt))
	rt.Setf("bus", func(rt lvnrt.Runtime) lvnrt.Dispatch {
		nrt := rt.Overlay("bus")
		bus := lvnrt.NewBus(nrt)
		nrt.Setd("bus", bus)
		return bus
	})
	ep := endpoint()
	log.Info("endpoint", ep)
	id := lvnrt.NewId("client")
	entry := lvnrt.NewEntry(rt, id, ep)
	defer entry.Close()
	log.Info("port", entry.Port())
	<-ctrlc
}

func endpoint() string {
	ep := os.Getenv("LV_NRT_ENDPOINT")
	if len(strings.TrimSpace(ep)) > 0 {
		return ep
	}
	return ":0"
}
