package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"strings"

	"github.com/samuelventura/laurelview/pkg/lvnrt"
)

//NOTICE: Two scenarios needed for proper cleanup
//monitoring stdin is for running as daemon service
//monitoring os.Interrupt is for running interactive
func main() {
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)
	dl := DefaultLog()
	defer dl("")
	defer dl("info", "exited")
	log := PrefixLogger(dl, "main")
	ep := endpoint()
	log.Info("endpoint", ep)
	dpm := lvnrt.NewDpm(log, ep, 400)
	defer dpm.Close(true)
	log.Info("port", dpm.Port())
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
	ep := os.Getenv("LV_DPM_ENDPOINT")
	if len(strings.TrimSpace(ep)) > 0 {
		return ep
	}
	return ":0"
}
