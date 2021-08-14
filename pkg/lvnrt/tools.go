package lvnrt

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
	"unicode"
)

func traceRecover(output Output) {
	r := recover()
	if r != nil {
		output("recover", r, string(debug.Stack()))
	}
}

func traceIfError(output Output, err error) {
	if err != nil {
		output("error", err)
	}
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func panicF(format string, args ...Any) {
	panic(fmt.Errorf(format, args...))
}

func panicLN(args ...Any) {
	panic(fmt.Sprintln(args...))
}

func assertTrue(flag bool, args ...Any) {
	if !flag {
		panic(fmt.Sprintln(args...))
	}
}

func readable(s string) string {
	b := new(strings.Builder)
	for _, c := range s {
		if unicode.IsControl(c) || unicode.IsSpace(c) {
			h := fmt.Sprintf("[%02X]", int(c))
			b.WriteString(h)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func clearDispatch(dispatchs map[string]Dispatch) {
	for name := range dispatchs {
		delete(dispatchs, name)
	}
}

func mapDispatch(found Output, nf404 Output, dispmap map[string]Dispatch) Dispatch {
	return func(mut *Mutation) {
		dispatch, ok := dispmap[mut.Name]
		if ok {
			found(mut)
			dispatch(mut)
		} else {
			nf404(mut)
		}
	}
}

func asyncDispatch(output Output, dispatch Dispatch) Dispatch {
	queue := make(chan *Mutation)
	loop := func() {
		defer traceRecover(output)
		for mut := range queue {
			dispatch(mut)
		}
	}
	go loop()
	return func(mut *Mutation) {
		//do not close queue nor state dispose
		//let map dispatch report the ignore
		queue <- mut
	}
}

func millis(ms int64) time.Duration {
	return time.Duration(ms) * time.Millisecond
}

func future(ms int64) time.Time {
	d := millis(ms)
	return time.Now().Add(d)
}

func sendChannel(channel Channel, any Any) {
	channel <- any
}

func closeChannel(channel Channel) {
	select {
	case <-channel:
	default:
		close(channel)
	}
}

func waitChannel(channel Channel, output Output) {
	output("waiting channel...")
	<-channel
	output("waiting channel done")
}

func toMap(any Any) Map {
	m := make(Map)
	e := reflect.ValueOf(any).Elem()
	t := e.Type()
	m["$type"] = t.Name()
	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		ft := t.Field(i)
		m[ft.Name] = f.Interface()
	}
	return m
}

func disposeArgs(arg Any) {
	action, ok := arg.(Action)
	if ok {
		action()
	}
	channel, ok := arg.(Channel)
	if ok {
		close(channel)
	}
}
