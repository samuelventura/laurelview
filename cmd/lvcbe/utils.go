package main

import (
	"log"
	"os"
	"strings"
)

func getenv(name string, defval string) string {
	value := os.Getenv(name)
	trimmed := strings.TrimSpace(value)
	if len(trimmed) > 0 {
		log.Println(name, value)
		return value
	}
	log.Println(name, defval)
	return defval
}
