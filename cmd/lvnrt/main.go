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
	defer CloseLog(log.Log) //wait flush
	defer log.Info("exited")
	defer WaitClose(rt.Close)
	defer TraceRecover(log.Warn)
	//gets slow after double connection attempted
	//takes >20s to connect on next attempt
	rt.SetValue("bus.dialtoms", 20000)
	rt.SetValue("bus.writetoms", 1000)
	rt.SetValue("bus.readtoms", 1000)
	rt.SetValue("bus.sleepms", 10)
	rt.SetValue("bus.retryms", 2000)
	rt.SetValue("bus.discardms", 100)
	rt.SetValue("bus.resetms", 400)
	ep := endpoint()
	log.Info("endpoint", ep)
	rt.SetValue("entry.endpoint", ep)
	rt.SetDispatch("hub", AsyncDispatch(log, NewHub(rt)))
	rt.SetDispatch("state", AsyncDispatch(log, NewState(rt)))
	checkinDispatch := NewCheckin(rt)
	defer checkinDispatch(Mn(":dispose"))
	rt.SetFactory("bus", func(rt Runtime) Dispatch { return NewBus(rt) })
	rt.SetDispatch("/ws/rt", checkinDispatch)
	entry := NewEntry(rt)
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
