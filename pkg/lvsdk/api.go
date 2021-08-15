package lvsdk

//Guidelines
//- catch invalid params on entry
//- global launch setting to info level
//- invalid rt conditions to debug level
//- recover to debug level
//- missing dispatch to debug level
//- trace is to see everything
//- debug is to see suspicious interactions
//- info is to announce launch setup
//- warn if to alert of suspicious global stats
//- error if to notify of suspicious global stats
//- error should also highlight hacking activity
// FIXME implement global stats monitor (managed, ...)

type Log = func(string, ...Any)
type Output = func(...Any)
type Map = map[string]Any
type Queue = chan Action
type Channel = chan Any
type Any = interface{}
type Action = func()
type Dispatch = func(*Mutation)
type Factory = func(Runtime) Dispatch

//Runtime Provides
//1) config
//2) factories
//3) dispatchs
//4) log (contextualized)
//5) self overlay
//6) cleaners

type Runtime interface {
	Getv(name string) Any
	Setv(name string, value Any)
	Setf(name string, factory Factory)
	Setd(name string, dispatch Dispatch)
	Setc(name string, cleaner Cleaner)
	Getc(name string) Cleaner
	Make(name string) Dispatch
	Post(name string, mut *Mutation)
	Log(level string, args ...Any)
	LevelOutput(level string) Output
	PrefixLog(prefix ...Any) Logger
	Clone() Runtime
	Close()
}

type Logger interface {
	Log(string, ...Any)
	Trace(...Any)
	Debug(...Any)
	Info(...Any)
	Warn(...Any)
	Error(...Any)
	Panic(...Any)
}

func NopAction()                  {}
func NopOutput(...Any)            {}
func NopDispatch(*Mutation)       {}
func NopFactory(Runtime) Dispatch { return NopDispatch }
