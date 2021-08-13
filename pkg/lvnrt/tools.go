package lvnrt

import (
	"fmt"
	"runtime/debug"
	"time"
)

func TraceRecover(output Output) {
	r := recover()
	if r != nil {
		output("trace", "recover", r, string(debug.Stack()))
	}
}

func TraceIfError(output Output, err error) {
	if err != nil {
		output("trace", "error", err)
	}
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicF(format string, args ...Any) {
	panic(fmt.Errorf(format, args...))
}

func PanicLN(args ...Any) {
	panic(fmt.Sprintln(args...))
}

func Assert(flag bool, args ...Any) {
	if !flag {
		panic(fmt.Sprintln(args...))
	}
}

func ClearDispatch(dispatchs map[string]Dispatch) {
	for name := range dispatchs {
		delete(dispatchs, name)
	}
}

func MapDispatch(output Output, dispmap map[string]Dispatch) Dispatch {
	return func(mut *Mutation) {
		dispatch, ok := dispmap[mut.Name]
		if !ok {
			output("unknown mutation", mut.Name)
			return
		}
		dispatch(mut)
	}
}

func PrefixOutput(output Output, prefix ...string) Output {
	return func(args ...Any) {
		output(prefix, args)
	}
}

func AsyncDispatch(output Output, dispatch Dispatch) Dispatch {
	disposed := false
	output = PrefixOutput(output, "async")
	queue := make(chan *Mutation)
	dispose := func(name string) {
		if name == "dispose" {
			if !disposed {
				close(queue)
			}
			disposed = true
		}
	}
	loop := func() {
		defer TraceRecover(output)
		for mut := range queue {
			dispatch(mut)
		}
	}
	go loop()
	return func(mut *Mutation) {
		defer dispose(mut.Name)
		if disposed {
			output("already disposed")
			return
		}
		queue <- mut
	}
}

func Millis(ms int64) time.Duration {
	return time.Duration(ms) * time.Millisecond
}

func Future(ms int64) time.Time {
	d := Millis(ms)
	return time.Now().Add(d)
}

func Send(channel Channel, any Any) {
	channel <- any
}
