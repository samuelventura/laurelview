package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/samuelventura/laurelview/pkg/lvnbe"
)

func main() {
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)
	output := lvnbe.DefaultOutput()
	defer output("") //wait flush
	defer lvnbe.TraceRecover(output)
	var db = relative("db3")
	dao := lvnbe.NewDao(db)
	state := lvnbe.NewState(dao)
	hub := lvnbe.NewHub(state)
	core := lvnbe.NewCore(hub, output)
	defer core.Close()
	ep := endpoint()
	output("info", "endpoint", ep)
	entry := lvnbe.NewEntry(core, output, ep)
	defer entry.Close()
	output("info", "port", entry.Port())
	<-ctrlc
}

func executable() string {
	exe, err := os.Executable()
	lvnbe.PanicIfError(err)
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
	ep := os.Getenv("LV_ENDPOINT")
	if len(strings.TrimSpace(ep)) > 0 {
		return ep
	}
	return ":0"
}
