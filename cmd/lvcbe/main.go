package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

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

	//entry
	endpoint := endpoint()
	driver := driver()
	source := source()
	log.Info("endpoint", endpoint)
	log.Info("driver", driver)
	log.Info("source", source)
	ctx.SetValue("apiEndpoint", endpoint)
	ctx.SetValue("dbDriver", driver)
	ctx.SetValue("dbSource", source)
	api := NewApi(ctx)
	defer api(Mn("dispose"))

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
	value := os.Getenv("LV_CBE_ENDPOINT")
	if len(strings.TrimSpace(value)) > 0 {
		return value
	}
	return ":0"
}

func driver() string {
	value := os.Getenv("LV_CBE_DRIVER")
	if len(strings.TrimSpace(value)) > 0 {
		return value
	}
	return "sqlite"
}

func source() string {
	value := os.Getenv("LV_CBE_SOURCE")
	if len(strings.TrimSpace(value)) > 0 {
		return value
	}
	return relative("db3")
}
