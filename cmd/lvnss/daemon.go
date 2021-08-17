package main

import (
	"os"
	"os/exec"
)

func daemon(sibling string, exit chan bool) chan bool {
	done := make(chan bool)
	path := relativeSibling(sibling)
	logger.Info("Daemon ", sibling, " ", path)
	outp := changeExtension(path, ".out.log")
	ff := os.O_APPEND | os.O_WRONLY | os.O_CREATE
	outf, err := os.OpenFile(outp, ff, 0600)
	panicIfError(err)
	errp := changeExtension(path, ".err.log")
	errf, err := os.OpenFile(errp, ff, 0600)
	panicIfError(err)
	cmd := exec.Command(path)
	cmd.Env = os.Environ()
	cmd.Stdout = outf
	cmd.Stderr = errf
	sin, err := cmd.StdinPipe()
	panicIfError(err)
	err = cmd.Start()
	panicIfError(err)
	go func() {
		defer traceRecover()
		defer close(done)
		defer outf.Close()
		defer errf.Close()
		defer traceRecover()
		go func() {
			defer traceRecover()
			defer sin.Close()
			<-exit
		}()
		err = cmd.Wait()
		panicIfError(err)
	}()
	return done
}
