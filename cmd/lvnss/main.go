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
	dbe = daemon(logger, "lvnbe", exit)
	return nil
}

func (p *program) Stop(s service.Service) error {
	close(exit)
	select {
	case <-dbe:
	case <-time.After(3 * time.Second):
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
		log.Fatal(err)
	}
	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	slog, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
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
	EnvironDefault(log, "LV_NBE_ENDPOINT", "0.0.0.0:31601")
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
