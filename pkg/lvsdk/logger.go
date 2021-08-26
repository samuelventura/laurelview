package lvsdk

import (
	"log"
	"os"
	"strings"
)

type nopLogger struct {
}

func NopLogger() Logger {
	logger := &nopLogger{}
	return logger
}

func (log *nopLogger) Log(level string, args ...Any) {
}

func (log *nopLogger) Trace(args ...Any) {
}

func (log *nopLogger) Debug(args ...Any) {
}

func (log *nopLogger) Info(args ...Any) {
}

func (log *nopLogger) Warn(args ...Any) {
}

func (log *nopLogger) Error(args ...Any) {
}

type simpleLogger struct {
	log Log
}

func SimpleLogger(log Log) Logger {
	logger := &simpleLogger{}
	logger.log = log
	return logger
}

func (log *simpleLogger) Log(level string, args ...Any) {
	log.log(level, args)
}

func (log *simpleLogger) Trace(args ...Any) {
	log.log("trace", args)
}

func (log *simpleLogger) Debug(args ...Any) {
	log.log("debug", args)
}

func (log *simpleLogger) Info(args ...Any) {
	log.log("info", args)
}

func (log *simpleLogger) Warn(args ...Any) {
	log.log("warn", args)
}

func (log *simpleLogger) Error(args ...Any) {
	log.log("error", args)
}

type prefixLogger struct {
	prefix []Any
	log    Log
}

func PrefixLogger(log Log, prefix ...Any) Logger {
	logger := &prefixLogger{}
	logger.log = log
	logger.prefix = prefix
	return logger
}

func (log *prefixLogger) Log(level string, args ...Any) {
	log.log(level, log.prefix, args)
}

func (log *prefixLogger) Trace(args ...Any) {
	log.log("trace", log.prefix, args)
}

func (log *prefixLogger) Debug(args ...Any) {
	log.log("debug", log.prefix, args)
}

func (log *prefixLogger) Info(args ...Any) {
	log.log("info", log.prefix, args)
}

func (log *prefixLogger) Warn(args ...Any) {
	log.log("warn", log.prefix, args)
}

func (log *prefixLogger) Error(args ...Any) {
	log.log("error", log.prefix, args)
}

func GoLogLogger() Logger {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lmicroseconds)
	var sb strings.Builder
	print := FlatPrintln(&sb)
	logger := func(level string, args ...Any) {
		sb.Reset()
		print(level, args)
		log.Print(sb.String())
	}
	return SimpleLogger(logger)
}
