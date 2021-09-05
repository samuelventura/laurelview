package main

import (
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	output := DefaultOutput()
	log := OutputLog(output)
	dos(SimpleLogger(log), output)
}

func dos(log Logger, output Output) {
	defer TraceRecover(log.Warn)
	plus := make(Channel)
	minus := make(Channel)
	ll := ll()
	ul := ul()
	fm := fm()
	vm := vm()
	go func() {
		count := 0
		for {
			for count < ul {
				<-plus
				count++
				log.Info("count", count)
			}
			for count > ll {
				<-minus
				count--
				log.Info("count", count)
			}
			time.Sleep(Millis(1))
		}
	}()
	for {
		SendChannel(plus, true)
		go func() {
			defer TraceRecover(log.Warn)
			defer SendChannel(minus, true)
			path := RelativeSibling("lvnup")
			cmd := exec.Command(path)
			cmd.Stdout = OutputWriter(output)
			cmd.Env = os.Environ()
			sin, err := cmd.StdinPipe()
			PanicIfError(err)
			defer sin.Close()
			err = cmd.Start()
			PanicIfError(err)
			tms := fm + rand.Intn(vm)
			log.Info("pid", cmd.Process.Pid, tms)
			time.Sleep(Millis(tms))
			err = sin.Close()
			PanicIfError(err)
			err = cmd.Wait()
			PanicIfError(err)
		}()
	}
}

func ll() int {
	return evint("LV_DOS_LL")
}

func ul() int {
	return evint("LV_DOS_UL")
}

func fm() int {
	return evint("LV_DOS_FM")
}

func vm() int {
	return evint("LV_DOS_VM")
}

func evint(evn string) int {
	evv := os.Getenv(evn)
	intv, err := strconv.Atoi(evv)
	PanicIfError(err)
	return intv
}
