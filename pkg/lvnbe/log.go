package lvnbe

import (
	"container/list"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

func DefaultOutput() Output {
	LogLevelFromEnv()
	var queue = make(chan Action, 128)
	go func() {
		for action := range queue {
			action()
		}
	}()
	output := func(level string, args ...Any) {
		if level != "" {
			if IsLogPrintable(level) {
				when := time.Now().Format("20060102T150405.000")
				queue <- func() {
					//log pkg won't flush on exit
					fmt.Fprint(os.Stdout, when, " ", level, " ")
					fmt.Fprintln(os.Stdout, args...)
				}
			}
		} else {
			done := make(Channel)
			queue <- func() {
				defer close(done)
				close(queue)
			}
			<-done
		}
	}
	output("info", "log-level", LogLevel)
	return output
}

var LogLevel = "info" //trace, debug, info

func LogLevelFromEnv() {
	logLevel := os.Getenv("LV_LOGLEVEL")
	if len(strings.TrimSpace(logLevel)) > 0 {
		LogLevel = logLevel
	}
}

func IsLogPrintable(level string) bool {
	switch LogLevel {
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

type testOutputState struct {
	list   *list.List
	mutex  sync.Mutex
	output Output
}

type testOutput interface {
	close()
	out(string, ...Any)
	assertEmpty(t *testing.T)
	matchNext(t *testing.T, args ...string) []string
	matchWait(t *testing.T, toms uint64, args ...string) []string
}

func newTestOutput() testOutput {
	to := &testOutputState{}
	to.output = DefaultOutput()
	to.list = list.New()
	return to
}

func (to *testOutputState) close() {
	to.output("") //wait flush
}

func (to *testOutputState) out(level string, args ...Any) {
	to.output(level, args...)
	array := make([]string, 0, 1+len(args))
	array = append(array, level)
	for _, arg := range args {
		array = append(array, fmt.Sprint(arg))
	}
	to.push(array)
}

func (to *testOutputState) compile(args []string) []*regexp.Regexp {
	matchers := make([]*regexp.Regexp, 0, len(args))
	for _, arg := range args {
		matchers = append(matchers, regexp.MustCompile(arg))
	}
	return matchers
}

func (to *testOutputState) matches(array []string, matchers []*regexp.Regexp) bool {
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

func (to *testOutputState) assertEmpty(t *testing.T) {
	array := to.pop()
	for array != nil {
		t.Errorf("not empty %v", array)
	}
}

func (to *testOutputState) matchNext(t *testing.T, args ...string) []string {
	array := to.pop()
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

func (to *testOutputState) matchWait(t *testing.T, toms uint64, args ...string) []string {
	matchers := to.compile(args)
	array := to.popWait(toms)
	for array != nil {
		to.output("matchWait", array)
		if to.matches(args, matchers) {
			return array
		}
		array = to.popWait(toms)
	}
	t.Errorf("No match for %v", args)
	return nil
}

func (to *testOutputState) push(array []string) {
	defer to.mutex.Unlock()
	to.mutex.Lock()
	to.list.PushBack(array)
}

func (to *testOutputState) pop() []string {
	defer to.mutex.Unlock()
	to.mutex.Lock()
	e := to.list.Front()
	if e != nil {
		to.list.Remove(e)
		return e.Value.([]string)
	}
	return nil
}

func (to *testOutputState) popWait(toms uint64) []string {
	onems := time.Duration(1) * time.Millisecond
	tod := time.Duration(toms) * time.Millisecond
	dl := time.Now().Add(tod)
	array := to.pop()
	for array == nil {
		time.Sleep(onems)
		if time.Now().After(dl) {
			return nil
		}
		array = to.pop()
	}
	return array
}
