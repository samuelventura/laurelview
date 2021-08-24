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
	rt.SetValue("bus.dialtoms", int64(20000))
	rt.SetValue("bus.writetoms", int64(1000))
	rt.SetValue("bus.readtoms", int64(1000))
	rt.SetValue("bus.sleepms", int64(10))
	rt.SetValue("bus.retryms", int64(2000))
	rt.SetValue("bus.discardms", int64(100))
	rt.SetValue("bus.resetms", int64(400))
	rt.SetCleaner("bus", NewCleaner(rt.PrefixLog("bus", "clean")))
	rt.SetDispatch("hub", AsyncDispatch(log.Debug, NewHub(rt)))
	stateDispatch := AsyncDispatch(log.Debug, NewState(rt))
	rt.SetDispatch("state", stateDispatch)
	defer stateDispatch(&Mutation{Name: "dispose"})
	rt.SetFactory("bus", func(rt Runtime) Dispatch { return NewBus(rt) })
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
