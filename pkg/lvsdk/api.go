package lvsdk

import "github.com/valyala/fasthttp"

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
type Dispatch = func(Mutation)
type Factory = func(Runtime) Dispatch
type Handler = func(ctx *fasthttp.RequestCtx)

//Runtime Provides
//1) config values
//2) factories
//3) dispatchs
//4) log (prefixed)
//5) cleaners (removed)
//6) ids (removed)
//7) self overlay (removed)

type Runtime interface {
	Cleaner() Cleaner
	SetValue(name string, value Any)
	SetFactory(name string, factory Factory)
	SetDispatch(name string, dispatch Dispatch)
	GetValue(name string) Any
	GetFactory(name string) Factory
	GetDispatch(name string) Dispatch
	Log(level string, args ...Any)
	LevelOutput(level string) Output
	PrefixLog(prefix ...Any) Logger
	Close() Channel
}

//panic conflicts with nop logger
//should a nop logger panic or not?
//keep panics separated until resolved
type Logger interface {
	PrefixLog(prefix ...Any) Logger
	Log(string, ...Any)
	Trace(...Any)
	Debug(...Any)
	Info(...Any)
	Warn(...Any)
	Error(...Any)
}

func NopAction()                          {}
func NopOutput(...Any)                    {}
func NopDispatch(Mutation)                {}
func NopFactory(Runtime) Dispatch         { return NopDispatch }
func NopHandler(ctx *fasthttp.RequestCtx) {}
