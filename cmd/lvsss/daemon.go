package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func logfp(fp string) string {
	dir := os.Getenv("LV_SSS_LOGS")
	if len(dir) > 0 {
		//filepath handles Windows/Unix separators
		return filepath.Join(dir, filepath.Base(fp))
	}
	return fp
}

func daemon(log Logger, sibling string, exit chan bool) chan bool {
	log = log.PrefixLog("daemon", sibling)
	done := make(chan bool)
	path := RelativeSibling(sibling)
	outp := ChangeExtension(path, ".out.log")
	ff := os.O_APPEND | os.O_WRONLY | os.O_CREATE
	outf, err := os.OpenFile(logfp(outp), ff, 0644)
	PanicIfError(err)
	errp := ChangeExtension(path, ".err.log")
	errf, err := os.OpenFile(logfp(errp), ff, 0644)
	PanicIfError(err)
	go func() {
		defer log.Debug("exited", path)
		defer TraceRecover(log.Error)
		defer close(done)
		defer outf.Close()
		defer errf.Close()
		run := func() {
			defer TraceRecover(log.Error)
			cmd := exec.Command(path)
			cmd.Env = os.Environ()
			cmd.Stdout = outf
			cmd.Stderr = errf
			sin, err := cmd.StdinPipe()
			PanicIfError(err)
			defer sin.Close()
			err = cmd.Start()
			PanicIfError(err)
			go func() {
				defer TraceRecover(log.Error)
				defer sin.Close()
				select {
				case <-exit:
				case <-done:
				}
			}()
			err = cmd.Wait()
			PanicIfError(err)
		}
		count := 0
		for {
			if count > 0 {
				time.Sleep(Millis(2000))
			}
			log.Info(count, path)
			run()
			count++
			select {
			case <-exit:
				return
			default:
				continue
			}
		}
	}()
	return done
}
