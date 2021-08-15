package main

import (
	"os"
	"os/signal"
	"strings"

	"github.com/samuelventura/laurelview/pkg/lvsdk"
)

func main() {
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)
	rt := lvsdk.DefaultRuntime()
	rt.Setv("bus.dialtoms", int64(800))
	rt.Setv("bus.writetoms", int64(800))
	rt.Setv("bus.readtoms", int64(800))
	rt.Setv("bus.sleepms", int64(10))
	rt.Setv("bus.retryms", int64(2000))
	rt.Setv("bus.discardms", int64(200))
	log := rt.PrefixLog("main")
	rt.Setc("bus", NewCleaner(rt.PrefixLog("bus", "clean")))
	defer rt.Close()
	defer log.Log("") //wait flush
	defer TraceRecover(log.Warn)
	rt.Setd("hub", AsyncDispatch(log.Trace, NewHub(rt)))
	rt.Setd("state", AsyncDispatch(log.Trace, NewState(rt)))
	defer rt.Post("state", &Mutation{Name: "dispose"})
	rt.Setf("bus", func(rt Runtime) Dispatch {
		nrt := rt.Clone()
		bus := AsyncDispatch(log.Trace, NewBus(nrt))
		nrt.Setd("self", bus)
		return bus
	})
	ep := endpoint()
	log.Info("endpoint", ep)
	id := NewId("client")
	entry := NewEntry(rt, id, ep)
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
