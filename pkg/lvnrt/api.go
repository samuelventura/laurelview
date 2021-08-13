package lvnrt

import (
	"fmt"
	"reflect"
	"strings"
)

type Log = func(string, ...Any)
type Output = func(...Any)
type Map = map[string]Any
type Queue = chan Action
type Channel = chan Any
type Any = interface{}
type Action = func()

type Dispatch = func(*Mutation)
type Factory = func(Runtime) Dispatch

type Logger interface {
	Log(string, ...Any)
	Trace(...Any)
	Debug(...Any)
	Info(...Any)
	Warn(...Any)
	Error(...Any)
	Panic(...Any)
}

type Mutation struct {
	Sid  string
	Name string
	Args Any
}

type AddArgs struct {
	Callback Dispatch
}

type SetupArgs struct {
	Items []*ItemArgs
}

type ItemArgs struct {
	Host  string
	Port  uint
	Slave uint
}

type QueryArgs struct {
	Index    uint
	Request  string
	Response string
}

type StatusArgs struct {
	Slave    string
	Request  string
	Response string
}

type BusArgs struct {
	Host string
	Port uint
}

type SlaveArgs struct {
	Slave uint
	Count uint
}

func NopAction()                  {}
func NopOutput(...Any)            {}
func NopDispatch(*Mutation)       {}
func NopFactory(Runtime) Dispatch { return NopDispatch }

func (m *Mutation) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "{%s,%s,%v}", m.Name, m.Sid, reflect.ValueOf(m.Args))
	return buf.String()
}
