package lvsdk

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
)

func EnvironFromFile(log Logger) {
	path := RelativeExtension(".config")
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	log.Info("loading config", path)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Warn("invalid config line", line)
			continue
		}
		log.Info("setting config", line)
		os.Setenv(parts[0], parts[1])
	}
}

func EnvironDefault(log Logger, name string, defval string) {
	val := os.Getenv(name)
	if len(strings.TrimSpace(val)) == 0 {
		log.Info("setting default", name, defval)
		os.Setenv(name, defval)
	} else {
		log.Info("found environ", name, val)
	}
}

func WaitExitSignal(action Action) {
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	action()

	//wait
	exit := make(Channel)
	go func() {
		defer close(exit)
		ioutil.ReadAll(os.Stdin)
	}()
	select {
	case <-ctrlc:
	case <-exit:
	}
}
