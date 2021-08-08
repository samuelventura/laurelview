package lvnbe

import (
	"fmt"
	"runtime/debug"
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
