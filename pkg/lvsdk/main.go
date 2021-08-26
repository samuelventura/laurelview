package lvsdk

import (
	"bufio"
	"os"
	"strings"
)

func EnvironFromFile(log Logger) {
	path := RelativeExtension(".config")
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	log.Info("Loading config", path)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Warn("Invalid config line", line)
			continue
		}
		log.Info("Setting config", line)
		os.Setenv(parts[0], parts[1])
	}
}

func EnvironDefault(log Logger, name string, defval string) {
	val := os.Getenv(name)
	if len(strings.TrimSpace(val)) == 0 {
		log.Info("Setting default", name, defval)
		os.Setenv(name, defval)
	} else {
		log.Info("Found environ", name, val)
	}
}
