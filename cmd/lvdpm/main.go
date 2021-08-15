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
	rt := DefaultRuntime()
	defer rt.Close()
	log := rt.PrefixLog("main")
	ep := endpoint()
	log.Info("endpoint", ep)
	dpm := lvnrt.NewDpm(log, ep, 400)
	defer dpm.Close()
	log.Info("port", dpm.Port())
	<-ctrlc
}

func endpoint() string {
	ep := os.Getenv("LV_DPM_ENDPOINT")
	if len(strings.TrimSpace(ep)) > 0 {
		return ep
	}
	return ":0"
}
