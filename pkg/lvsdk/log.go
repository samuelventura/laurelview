package lvsdk

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func flattenArgs() func(cb func(Any), args ...Any) {
	flatten := func(func(Any), ...Any) {}
	flatten = func(cb func(Any), args ...Any) {
		for _, arg := range args {
			switch v := arg.(type) {
			case []Any:
				flatten(cb, v...)
			default:
				cb(arg)
			}
		}
	}
	return flatten
}

func FlatPrintln(writer io.Writer) func(args ...Any) {
	flatten := flattenArgs()
	return func(args ...Any) {
		count := 0
		cb := func(arg Any) {
			if count > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, arg)
			count++
		}
		flatten(cb, args)
		fmt.Fprintln(writer)
	}
}

func LevelOutput(log Log, level string) Output {
	return func(args ...Any) {
		log(level, args...)
	}
}

//idempotent close
//first close waits for flush
//second close returns immediately
func DefaultOutput() Output {
	done := make(Channel)
	queue := make(chan []Any, 128)
	print := FlatPrintln(os.Stdout)
	loop := func() {
		for args := range queue {
			if len(args) == 0 {
				select {
				case <-done:
				default:
					close(done)
				}
			} else {
				print(args)
			}
		}
	}
	go loop()
	return func(args ...Any) {
		queue <- args
		if len(args) == 0 {
			<-done
		}
	}
}

func CloseOutput(output Output) {
	output()
}

func CloseLog(log Log) {
	log(":dispose")
}

func DefaultLog() Log {
	logLevelFromEnv()
	output := DefaultOutput()
	return OutputLog(output)
}

func OutputLog(output Output) Log {
	pid := os.Getpid()
	log := func(level string, args ...Any) {
		switch level {
		case ":dispose":
			CloseOutput(output)
		default:
			if isLogPrintable(level) {
				now := time.Now()
				when := now.Format("20060102T150405.000")
				output(when, pid, level, args)
			}
		}
	}
	log("info", "level", logLevel)
	return log
}

var logLevel = "info" //trace, debug, info

func logLevelFromEnv() {
	logLevelEnv := os.Getenv("LV_LOGLEVEL")
	if len(strings.TrimSpace(logLevelEnv)) > 0 {
		logLevel = logLevelEnv
	}
}

func isLogPrintable(level string) bool {
	switch logLevel {
	case "trace":
		return true
	case "debug":
		switch level {
		case "trace":
			return false
		default:
			return true
		}
	default:
		switch level {
		case "trace", "debug":
			return false
		default:
			return true
		}
	}
}

type outputWriterDso struct {
	output Output
	sb     strings.Builder
}

func (ow *outputWriterDso) Write(bytes []byte) (n int, err error) {
	nl := byte(0x0A)
	for _, b := range bytes {
		if b == nl {
			ow.output(ow.sb.String())
			ow.sb.Reset()
		} else {
			ow.sb.WriteByte(b)
		}
	}
	return len(bytes), nil
}

func OutputWriter(output Output) io.Writer {
	ow := &outputWriterDso{}
	ow.output = output
	return ow
}
