package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/kardianos/service"
)

var logger service.Logger
var exit chan bool
var dbe chan bool

type program struct{}

func (p *program) Start(s service.Service) (err error) {
	exit = make(chan bool)
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("recover %v", r)
		}
	}()
	dbe = daemon("lvnbe", exit)
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
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	//after logger created
	environFromFile()
	environDefaults()
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
