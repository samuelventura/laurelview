package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func environDefaults() {
	environDefault("LV_NBE_ENDPOINT", ":31601")
	environDefault("LV_NRT_ENDPOINT", ":31602")
}

func environDefault(name string, defval string) {
	val := os.Getenv(name)
	if len(strings.TrimSpace(val)) == 0 {
		logger.Info("Setting default ", name, "=", defval)
		os.Setenv(name, defval)
	}
}

func environFromFile() {
	path := relativeExtension(".config")
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	logger.Info("Loading config ", path)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			logger.Warning("Invalid config line ", line)
			continue
		}
		logger.Info("Setting config ", line)
		os.Setenv(parts[0], parts[1])
	}
}

//different filename same extension
func relativeSibling(sibling string) string {
	exe := executable()
	dir := filepath.Dir(exe)
	base := filepath.Base(exe)
	ext := filepath.Ext(base) //includes .
	file := sibling + ext
	return filepath.Join(dir, file)
}

//same file name different extension
func relativeExtension(ext string) string {
	path := executable()
	return changeExtension(path, ext)
}

func changeExtension(path string, next string) string {
	ext := filepath.Ext(path) //includes .
	npath := strings.TrimSuffix(path, ext)
	return npath + next
}

func executable() string {
	exe, err := os.Executable()
	panicIfError(err)
	return exe
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func traceRecover() {
	r := recover()
	if r != nil {
		logger.Error("recover", r)
	}
}
