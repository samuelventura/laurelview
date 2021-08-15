package lvsdk

import (
	"fmt"
)

//FIXME cumulative panic recover
func TraceRecover(output Output) {
	r := recover()
	if r != nil {
		output("recover", r)
		//output("recover", r, string(debug.Stack()))
	}
}

func TraceIfError(output Output, err error) {
	if err != nil {
		output("error", err)
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

func AssertTrue(flag bool, args ...Any) {
	if !flag {
		PanicLN(args...)
	}
}
