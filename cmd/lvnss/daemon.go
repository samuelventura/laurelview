package main

import (
	"os"
	"os/exec"
)

//FIXME implement reload strategy
func daemon(log Logger, sibling string, exit chan bool) chan bool {
	done := make(chan bool)
	path := RelativeSibling(sibling)
	log.Info("daemon", sibling, path)
	outp := ChangeExtension(path, ".out.log")
	ff := os.O_APPEND | os.O_WRONLY | os.O_CREATE
	outf, err := os.OpenFile(outp, ff, 0644)
	PanicIfError(err)
	errp := ChangeExtension(path, ".err.log")
	errf, err := os.OpenFile(errp, ff, 0644)
	PanicIfError(err)
	cmd := exec.Command(path)
	cmd.Env = os.Environ()
	cmd.Stdout = outf
	cmd.Stderr = errf
	sin, err := cmd.StdinPipe()
	PanicIfError(err)
	err = cmd.Start()
	PanicIfError(err)
	go func() {
		defer log.Debug("exited", path)
		defer TraceRecover(log.Error)
		defer close(done)
		defer outf.Close()
		defer errf.Close()
		defer TraceRecover(log.Error)
		go func() {
			defer TraceRecover(log.Error)
			defer sin.Close()
			<-exit
		}()
		err = cmd.Wait()
		PanicIfError(err)
	}()
	return done
}
