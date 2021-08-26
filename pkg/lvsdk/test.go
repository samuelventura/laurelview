package lvsdk

import (
	"container/list"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

type testOutputDso struct {
	list    *list.List
	mutex   sync.Mutex
	flatten func(cb func(Any), args ...Any)
	log     Log
}

type TestOutput interface {
	AssertEmpty(t *testing.T)
	Log(level string, args ...Any)
	MatchWait(t *testing.T, toms int, args ...string) []string
	MatchNext(t *testing.T, args ...string) []string
	Dispatch(name string) Dispatch
	Logger(prefix ...Any) Logger
	Trace(...Any)
	Debug(...Any)
	Info(...Any)
	Warn(...Any)
	Error(...Any)
	Close()
}

func NewTestOutput() TestOutput {
	to := &testOutputDso{}
	to.flatten = flattenArgs()
	to.log = DefaultLog()
	to.list = list.New()
	return to
}

func (to *testOutputDso) Close() {
	CloseLog(to.log) //wait flush
}

func (to *testOutputDso) Dispatch(name string) Dispatch {
	return func(mut Mutation) {
		to.Log("trace", name, mut)
	}
}

func (to *testOutputDso) Logger(prefix ...Any) Logger {
	return PrefixLogger(to.Log, prefix...)
}

func (to *testOutputDso) Log(level string, args ...Any) {
	to.log(level, args...)
	var array []string
	array = append(array, level)
	cb := func(arg Any) {
		array = append(array, fmt.Sprint(arg))
	}
	to.flatten(cb, args)
	to.pushArray(array)
}

func (to *testOutputDso) Trace(args ...Any) {
	to.Log("trace", args)
}

func (to *testOutputDso) Debug(args ...Any) {
	to.Log("debug", args)
}

func (to *testOutputDso) Info(args ...Any) {
	to.Log("info", args)
}

func (to *testOutputDso) Warn(args ...Any) {
	to.Log("warn", args)
}

func (to *testOutputDso) Error(args ...Any) {
	to.Log("error", args)
}

func (to *testOutputDso) AssertEmpty(t *testing.T) {
	array := to.popArray()
	for array != nil {
		PanicF("not empty %v", array)
	}
}

func (to *testOutputDso) MatchNext(t *testing.T, args ...string) []string {
	array := to.popArray()
	for array == nil {
		PanicF("empty pop")
		return nil
	}
	matchers := to.compile(args)
	if !to.matches(args, matchers) {
		PanicF("%v is no match for %v", array, args)
		return nil
	}
	return array
}

func (to *testOutputDso) MatchWait(t *testing.T, toms int, args ...string) []string {
	matchers := to.compile(args)
	dl := Future(toms)
	array := to.popWait(toms)
	for array != nil {
		if to.matches(array, matchers) {
			return array
		}
		if time.Now().After(dl) {
			break
		}
		array = to.popWait(toms)
	}
	PanicF("no match for %v", args)
	return nil
}

func (to *testOutputDso) compile(args []string) []func(string) bool {
	matchers := make([]func(string) bool, 0, len(args))
	for _, arg := range args {
		parts := strings.Split(arg, "||")
		if len(parts) == 1 {
			matchers = append(matchers, regexp.MustCompile(arg).MatchString)
		} else {
			submatchers := make([]func(string) bool, 0, len(parts))
			for _, part := range parts {
				submatchers = append(submatchers, regexp.MustCompile(part).MatchString)
			}
			matchers = append(matchers, func(value string) bool {
				for _, matcher := range submatchers {
					if !matcher(value) {
						return false
					}
				}
				return true
			})
		}
	}
	return matchers
}

func (to *testOutputDso) matches(array []string, matchers []func(string) bool) bool {
	if len(array) >= len(matchers) {
		count := len(matchers)
		for i, matcher := range matchers {
			value := array[i]
			if matcher(value) {
				count--
			}
		}
		if count == 0 {
			return true
		}
	}
	return false
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

func (to *testOutputDso) popWait(toms int) []string {
	onems := Millis(1)
	dl := Future(toms)
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
