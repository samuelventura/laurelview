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
//- warn is to alert of suspicious global stats
//- error is to notify of suspicious global stats
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
//4) log (prefixed)
//5) cleaners
//6) self overlay (removed)

type Runtime interface {
	SetValue(name string, value Any)
	SetFactory(name string, factory Factory)
	SetDispatch(name string, dispatch Dispatch)
	SetCleaner(name string, cleaner Cleaner)
	GetValue(name string) Any
	GetCleaner(name string) Cleaner
	GetFactory(name string) Factory
	GetDispatch(name string) Dispatch
	Log(level string, args ...Any)
	LevelOutput(level string) Output
	PrefixLog(prefix ...Any) Logger
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
