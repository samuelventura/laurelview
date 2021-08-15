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

func flatPrintln(writer io.Writer) func(args ...Any) {
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

type prefixLog struct {
	prefix []Any
	log    Log
}

func PrefixLogger(log Log, prefix ...Any) Logger {
	l := &prefixLog{}
	l.log = log
	l.prefix = prefix
	return l
}

func (l *prefixLog) Log(level string, args ...Any) {
	l.log(level, l.prefix, args)
}

func (l *prefixLog) Trace(args ...Any) {
	l.log("trace", l.prefix, args)
}

func (l *prefixLog) Debug(args ...Any) {
	l.log("debug", l.prefix, args)
}

func (l *prefixLog) Info(args ...Any) {
	l.log("info", l.prefix, args)
}

func (l *prefixLog) Warn(args ...Any) {
	l.log("warn", l.prefix, args)
}

func (l *prefixLog) Error(args ...Any) {
	l.log("error", l.prefix, args)
}

func (l *prefixLog) Panic(args ...Any) {
	l.log("panic", l.prefix, args)
}

func DefaultOutput() Output {
	done := make(Channel)
	queue := make(chan []Any, 128)
	print := flatPrintln(os.Stdout)
	loop := func() {
		for args := range queue {
			if len(args) == 0 {
				//do not close queue
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
		//wait for flush
		if len(args) == 0 {
			<-done
		}
	}
}

func DefaultLog() Log {
	logLevelFromEnv()
	output := DefaultOutput()
	log := func(level string, args ...Any) {
		if level != "" {
			if isLogPrintable(level) {
				now := time.Now()
				when := now.Format("20060102T150405.000")
				output(when, level, args)
			}
			if level == "panic" {
				//FIXME overlapped output
				buf := new(strings.Builder)
				print := flatPrintln(buf)
				print(args)
				panic(buf.String())
			}
		} else {
			output()
		}
	}
	log("info", "log-level", logLevel)
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
