package lvnrt

import (
	"container/list"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
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

func defaultOutput() Output {
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

func defaultLog() Log {
	logLevelFromEnv()
	output := defaultOutput()
	log := func(level string, args ...Any) {
		if level != "" {
			if isLogPrintable(level) {
				now := time.Now()
				when := now.Format("20060102T150405.000")
				output(when, level, args)
			}
			if level == "panic" {
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

func levelOutput(log Log, level string) Output {
	return func(args ...Any) {
		log(level, args...)
	}
}

type prefixLog struct {
	prefix []Any
	log    Log
}

func prefixLogger(log Log, prefix ...Any) Logger {
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

type testOutputDso struct {
	list    *list.List
	mutex   sync.Mutex
	flatten func(cb func(Any), args ...Any)
	log     Log
}

type testOutput interface {
	close()
	logger(prefix ...Any) Logger
	push(level string, args ...Any)
	assertEmpty(t *testing.T)
	dispatch(name string) Dispatch
	matchNext(t *testing.T, args ...string) []string
	matchWait(t *testing.T, toms uint64, args ...string) []string
}

func newTestOutput() testOutput {
	to := &testOutputDso{}
	to.flatten = flattenArgs()
	to.log = defaultLog()
	to.list = list.New()
	return to
}

func (to *testOutputDso) close() {
	to.log("") //wait flush
}

func (to *testOutputDso) dispatch(name string) Dispatch {
	return func(mut *Mutation) {
		to.push("trace", name, mut)
	}
}

func (to *testOutputDso) logger(prefix ...Any) Logger {
	return prefixLogger(to.push, prefix...)
}

func (to *testOutputDso) push(level string, args ...Any) {
	to.log(level, args...)
	var array []string
	array = append(array, level)
	cb := func(arg Any) {
		array = append(array, fmt.Sprint(arg))
	}
	to.flatten(cb, args)
	to.pushArray(array)
}

func (to *testOutputDso) compile(args []string) []*regexp.Regexp {
	matchers := make([]*regexp.Regexp, 0, len(args))
	for _, arg := range args {
		matchers = append(matchers, regexp.MustCompile(arg))
	}
	return matchers
}

func (to *testOutputDso) matches(array []string, matchers []*regexp.Regexp) bool {
	if len(array) >= len(matchers) {
		count := len(matchers)
		for i, matcher := range matchers {
			value := array[i]
			if matcher.MatchString(value) {
				count--
			}
		}
		if count == 0 {
			return true
		}
	}
	return false
}

func (to *testOutputDso) assertEmpty(t *testing.T) {
	array := to.popArray()
	for array != nil {
		t.Errorf("not empty %v", array)
	}
}

func (to *testOutputDso) matchNext(t *testing.T, args ...string) []string {
	array := to.popArray()
	for array == nil {
		t.Errorf("empty pop")
		return nil
	}
	matchers := to.compile(args)
	if !to.matches(args, matchers) {
		t.Errorf("%v is no match for %v", array, args)
		return nil
	}
	return array
}

func (to *testOutputDso) matchWait(t *testing.T, toms uint64, args ...string) []string {
	matchers := to.compile(args)
	array := to.popWait(toms)
	for array != nil {
		if to.matches(array, matchers) {
			return array
		}
		array = to.popWait(toms)
	}
	t.Errorf("no match for %v", args)
	return nil
}

func (to *testOutputDso) pushArray(array []string) {
	defer to.mutex.Unlock()
	to.mutex.Lock()
	to.list.PushBack(array)
}

func (to *testOutputDso) popArray() []string {
	defer to.mutex.Unlock()
	to.mutex.Lock()
	e := to.list.Front()
	if e != nil {
		to.list.Remove(e)
		return e.Value.([]string)
	}
	return nil
}

func (to *testOutputDso) popWait(toms uint64) []string {
	onems := time.Duration(1) * time.Millisecond
	tod := time.Duration(toms) * time.Millisecond
	dl := time.Now().Add(tod)
	array := to.popArray()
	for array == nil {
		time.Sleep(onems)
		if time.Now().After(dl) {
			return nil
		}
		array = to.popArray()
	}
	return array
}
