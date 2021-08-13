package lvnrt

import (
	"sync"
	"time"
)

//FIXME self overlay
//provides
//1) config
//2) factories
//3) dispatchs
//4) log
//5) self overlay
type Runtime interface {
	Getv(name string) Any
	Setv(name string, value Any)
	Setf(name string, factory Factory)
	Setd(name string, dispatch Dispatch)
	Make(name string) Dispatch
	Post(name string, mut *Mutation)
	Log(level string, args ...Any)
	Overlay(self string) Runtime
	Managed(mid string) Action
	ManagedWait()
}

type runtimeDso struct {
	log       Log
	self      string
	values    map[string]Any
	factories map[string]Factory
	dispatchs map[string]Dispatch
	managed   map[string]Count
	mutex     sync.Mutex
}

func NewRuntime(log Log) Runtime {
	rt := &runtimeDso{}
	rt.log = log
	rt.self = "self"
	rt.values = make(map[string]interface{})
	rt.factories = make(map[string]Factory)
	rt.dispatchs = make(map[string]Dispatch)
	rt.managed = make(map[string]Count)
	return rt
}

func (rt *runtimeDso) Overlay(self string) Runtime {
	n := *rt
	n.self = self
	assertTrue(rt.self == "self", "Not cloned")
	return &n
}

func (rt *runtimeDso) Managed(mid string) Action {
	defer rt.mutex.Unlock()
	rt.mutex.Lock()
	count, ok := rt.managed[mid]
	if !ok {
		count = NewCount()
		rt.managed[mid] = count
	}
	count.Inc()
	return func() {
		defer rt.mutex.Unlock()
		rt.mutex.Lock()
		count.Dec()
		if count.Count() == 0 {
			delete(rt.managed, mid)
		}
	}
}

func (rt *runtimeDso) ManagedWait() {
	for {
		rt.mutex.Lock()
		count := len(rt.managed)
		rt.mutex.Unlock()
		if count == 0 {
			return
		}
		time.Sleep(millis(1))
	}
}

func (rt *runtimeDso) Getv(name string) Any {
	return rt.values[name]
}

func (rt *runtimeDso) Setv(name string, value Any) {
	rt.values[name] = value
}

func (rt *runtimeDso) Setf(name string, factory Factory) {
	rt.factories[name] = factory
}

func (rt *runtimeDso) Setd(name string, dispatch Dispatch) {
	rt.dispatchs[name] = dispatch
}

func (rt *runtimeDso) Make(name string) Dispatch {
	return rt.factories[name](rt)
}

func (rt *runtimeDso) Post(name string, mut *Mutation) {
	if name == "self" {
		name = rt.self
	}
	rt.dispatchs[name](mut)
}

func (rt *runtimeDso) Log(level string, args ...Any) {
	rt.log(level, args...)
}
