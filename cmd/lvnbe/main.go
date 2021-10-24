package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/YeicoLabs/laurelview/pkg/lvndb"
	"github.com/YeicoLabs/laurelview/pkg/lvnrt"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	//runtime 0
	dl := DefaultLog()
	defer CloseLog(dl)
	ctx := NewContext(dl)
	defer WaitClose(ctx.Close)
	log := ctx.PrefixLog("main")
	defer TraceRecover(log.Warn)

	//runtime 1 /ws/rt
	ctx1 := NewContext(dl)
	log1 := ctx1.PrefixLog("rt")
	defer log1.Info("exited")
	defer WaitClose(ctx1.Close)
	//gets slow after double connection attempted
	//takes >20s to connect on next attempt
	ctx1.SetValue("bus.dialtoms", 20000)
	ctx1.SetValue("bus.writetoms", 1000)
	ctx1.SetValue("bus.readtoms", 1000)
	ctx1.SetValue("bus.sleepms", 10)
	ctx1.SetValue("bus.retryms", 2000)
	ctx1.SetValue("bus.discardms", 100)
	ctx1.SetValue("bus.resetms", 400)
	ctx1.SetFactory("bus", func(ctx Context) Dispatch { return lvnrt.NewBus(ctx) })
	hub1 := AsyncDispatch(log1, lvnrt.NewHub(ctx1))
	ctx1.SetDispatch("hub", hub1)
	ctx1.SetDispatch("state", AsyncDispatch(log1, lvnrt.NewState(ctx1)))
	check1 := lvnrt.NewCheck(ctx1)
	defer check1(Mn(":dispose"))

	//runtime 2 /ws/db
	var db2 = dbpath()
	log.Info("db", db2)
	var dao2 = NewDao(db2)
	defer dao2.Close()
	ctx2 := NewContext(dl)
	log2 := ctx2.PrefixLog("db")
	defer log2.Info("exited")
	defer WaitClose(ctx2.Close)
	ctx2.SetValue("dao", dao2)
	hub2 := AsyncDispatch(log2, lvndb.NewHub(ctx2))
	ctx2.SetDispatch("hub", hub2)
	ctx2.SetDispatch("state", AsyncDispatch(log2, lvndb.NewState(ctx2)))
	check2 := lvndb.NewCheck(ctx2)
	defer check2(Mn(":dispose"))

	//entry
	ep := endpoint()
	log.Info("endpoint", ep)
	ctx.SetValue("entry.endpoint", ep)
	ctx.SetValue("entry.buflen", 256)
	ctx.SetValue("entry.wtoms", 4000)
	ctx.SetValue("entry.rtoms", 4000)
	ctx.SetValue("entry.static", NewEmbedHandler(log))
	ctx.SetDispatch("/ws/rt", check1)
	ctx.SetDispatch("/ws/db", check2)
	entry := lvnrt.NewEntry(ctx)
	defer WaitClose(entry.Close)
	log.Info("port", entry.Port())
	ticker1 := time.NewTicker(1 * time.Second)
	defer ticker1.Stop()
	go func() {
		for range ticker1.C {
			hub1(Mns(":ping", "main"))
			hub2(Mns(":ping", "main"))
		}
	}()

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		output := make(Channel)
		entry.Status(output)
		for line := range output {
			fmt.Fprintln(w, line)
		}
	})
	go func() {
		//https://pkg.go.dev/net/http/pprof
		//https://golang.org/doc/diagnostics
		ep := os.Getenv("LV_NBE_DEBUG")
		log.Info("pprof", ep)
		log.Debug(http.ListenAndServe(ep, nil))
	}()

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

func dbpath() string {
	db := os.Getenv("LV_NBE_DATABASE")
	if len(strings.TrimSpace(db)) > 0 {
		return db
	}
	return relative("db3")
}
