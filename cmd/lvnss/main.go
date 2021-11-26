package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kardianos/service"
)

var exit chan bool
var dbe chan bool
var dup chan bool
var dpm chan bool
var logger Logger

type program struct{}

func (p *program) Start(s service.Service) (err error) {
	exit = make(chan bool)
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("recover %v", r)
		}
	}()
	dpm = daemon(logger, "lvdpm", exit)
	dbe = daemon(logger, "lvnbe", exit)
	dup = daemon(logger, "lvnup", exit)
	return nil
}

func (p *program) Stop(s service.Service) error {
	close(exit)
	select {
	case <-dbe:
	case <-time.After(3 * time.Second):
	}
	select {
	case <-dup:
	case <-time.After(1 * time.Second):
	}
	select {
	case <-dpm:
	case <-time.After(1 * time.Second):
	}
	return nil
}

func main() {
	//-service install, uninstall, start, stop, restart
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()
	svcConfig := &service.Config{
		Name:        "LaurelView",
		DisplayName: "LaurelView Service",
		Description: "LaurelView https://laurelview.io",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Panicln(err)
	}
	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Panicln(err)
		}
		return
	}
	slog, err := s.Logger(nil)
	if err != nil {
		log.Panicln(err)
	}
	logger = Wrap(slog)
	//after logger created
	EnvironFromFile(logger)
	environDefaults(logger)
	err = s.Run()
	if err != nil {
		slog.Error(err)
	}
}

func environDefaults(log Logger) {
	//defaults for windows installer to work out of the box
	EnvironDefault(log, "LV_NBE_ENDPOINT", "0.0.0.0:31601")
	EnvironDefault(log, "LV_DPM_ENDPOINT", "127.0.0.1:31602")
	EnvironDefault(log, "LV_NUP_ENDPOINT", "127.0.0.1:31601")
	EnvironDefault(log, "LV_NBE_DEBUG", "127.0.0.1:31001")
	EnvironDefault(log, "LV_NSS_LOGS", "")
}

func Wrap(slog service.Logger) Logger {
	var sb strings.Builder
	print := FlatPrintln(&sb)
	log := func(level string, args ...Any) {
		sb.Reset()
		print(args)
		switch level {
		case "warn":
			slog.Warning(sb.String())
		case "error":
			slog.Error(sb.String())
		default:
			slog.Info(sb.String())
		}
	}
	return SimpleLogger(log)
}
