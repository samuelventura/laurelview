package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"strings"

	"github.com/samuelventura/laurelview/pkg/lvsdk"
)

func main() {
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)
	rt := lvsdk.DefaultRuntime()
	log := rt.PrefixLog("main")
	defer log.Log("") //wait flush
	defer log.Info("exited")
	defer rt.Close()
	defer TraceRecover(log.Warn)
	//gets slow after double connection attempted
	//takes >20s to connect on next attempt
	rt.Setv("bus.dialtoms", int64(20000))
	rt.Setv("bus.writetoms", int64(1000))
	rt.Setv("bus.readtoms", int64(1000))
	rt.Setv("bus.sleepms", int64(10))
	rt.Setv("bus.retryms", int64(2000))
	rt.Setv("bus.discardms", int64(100))
	rt.Setv("bus.resetms", int64(400))
	rt.Setc("bus", NewCleaner(rt.PrefixLog("bus", "clean")))
	rt.Setd("hub", AsyncDispatch(log.Debug, NewHub(rt)))
	rt.Setd("state", AsyncDispatch(log.Debug, NewState(rt)))
	defer rt.Post("state", &Mutation{Name: "dispose"})
	rt.Setf("bus", func(rt Runtime) Dispatch { return NewBus(rt) })
	ep := endpoint()
	log.Info("endpoint", ep)
	id := NewId("client")
	entry := NewEntry(rt, id, ep)
	defer entry.Close()
	log.Info("port", entry.Port())
	exit := make(chan bool)
	go stdin(exit)
	select {
	case <-ctrlc:
	case <-exit:
	}
}

func stdin(exit chan bool) {
	defer close(exit)
	ioutil.ReadAll(os.Stdin)
}

func endpoint() string {
	ep := os.Getenv("LV_NRT_ENDPOINT")
	if len(strings.TrimSpace(ep)) > 0 {
		return ep
	}
	return ":0"
}
