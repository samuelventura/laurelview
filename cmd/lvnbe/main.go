package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/samuelventura/laurelview/pkg/lvndb"
	"github.com/samuelventura/laurelview/pkg/lvnrt"

	"net/http"
	"net/http/pprof"
)

func main() {
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	//runtime 0
	dl := DefaultLog()
	defer CloseLog(dl)
	rt := NewRuntime(dl)
	defer WaitClose(rt.Close)
	log := rt.PrefixLog("main")
	defer TraceRecover(log.Warn)

	go func() {
		//https://pkg.go.dev/net/http/pprof
		//https://golang.org/doc/diagnostics
		mux := http.NewServeMux()
		mux.HandleFunc("/custom_debug_path/profile", pprof.Profile)
		ep := os.Getenv("LV_NBE_DEBUG")
		log.Info("pprof", ep)
		log.Debug(http.ListenAndServe(ep, mux))
	}()

	//runtime 1 /ws/rt
	rt1 := NewRuntime(dl)
	log1 := rt1.PrefixLog("rt")
	defer log1.Info("exited")
	defer WaitClose(rt1.Close)
	//gets slow after double connection attempted
	//takes >20s to connect on next attempt
	rt1.SetValue("bus.dialtoms", 20000)
	rt1.SetValue("bus.writetoms", 1000)
	rt1.SetValue("bus.readtoms", 1000)
	rt1.SetValue("bus.sleepms", 10)
	rt1.SetValue("bus.retryms", 2000)
	rt1.SetValue("bus.discardms", 100)
	rt1.SetValue("bus.resetms", 400)
	rt1.SetFactory("bus", func(rt Runtime) Dispatch { return lvnrt.NewBus(rt) })
	rt1.SetDispatch("hub", AsyncDispatch(log1, lvnrt.NewHub(rt1)))
	rt1.SetDispatch("state", AsyncDispatch(log1, lvnrt.NewState(rt1)))
	check1 := lvnrt.NewCheck(rt1)
	defer check1(Mn(":dispose"))

	//runtime 2 /ws/db
	var db2 = relative("db3")
	var dao2 = NewDao(db2)
	defer dao2.Close()
	rt2 := NewRuntime(dl)
	log2 := rt2.PrefixLog("db")
	defer log2.Info("exited")
	defer WaitClose(rt2.Close)
	rt2.SetValue("dao", dao2)
	rt2.SetDispatch("hub", AsyncDispatch(log2, lvndb.NewHub(rt2)))
	rt2.SetDispatch("state", AsyncDispatch(log2, lvndb.NewState(rt2)))
	check2 := lvndb.NewCheck(rt2)
	defer check2(Mn(":dispose"))

	//entry
	ep := endpoint()
	log.Info("endpoint", ep)
	rt.SetValue("entry.endpoint", ep)
	rt.SetValue("entry.buflen", 0)
	rt.SetValue("entry.static", NewEmbedHandler(log))
	rt.SetDispatch("/ws/rt", check1)
	rt.SetDispatch("/ws/db", check2)
	entry := lvnrt.NewEntry(rt)
	defer WaitClose(entry.Close)
	log.Info("port", entry.Port())

	//wait
	exit := make(Channel)
	go stdin(exit)
	select {
	case <-ctrlc:
	case <-exit:
	}
}

func stdin(exit Channel) {
	defer close(exit)
	ioutil.ReadAll(os.Stdin)
}

func executable() string {
	exe, err := os.Executable()
	PanicIfError(err)
	return exe
}

func relative(ext string) string {
	exe := executable()
	dir := filepath.Dir(exe)
	base := filepath.Base(exe)
	file := base + "." + ext
	return filepath.Join(dir, file)
}

func endpoint() string {
	ep := os.Getenv("LV_NBE_ENDPOINT")
	if len(strings.TrimSpace(ep)) > 0 {
		return ep
	}
	return ":0"
}
