package lvsdk

func DefaultRuntime() Runtime {
	return NewRuntime(DefaultLog())
}

type runtimeDso struct {
	log       Log
	self      Dispatch
	values    map[string]Any
	factories map[string]Factory
	dispatchs map[string]Dispatch
	cleaners  map[string]Cleaner
}

func NewRuntime(log Log) Runtime {
	rt := &runtimeDso{}
	rt.log = log
	rt.self = NopDispatch
	rt.values = make(map[string]Any)
	rt.factories = make(map[string]Factory)
	rt.dispatchs = make(map[string]Dispatch)
	rt.cleaners = make(map[string]Cleaner)
	return rt
}

func (rt *runtimeDso) Clone() Runtime {
	clone := &runtimeDso{}
	clone.log = rt.log
	clone.self = rt.self
	clone.values = rt.values
	clone.factories = rt.factories
	clone.dispatchs = rt.dispatchs
	clone.cleaners = rt.cleaners
	return clone
}

func (rt *runtimeDso) PrefixLog(prefix ...Any) Logger {
	return PrefixLogger(rt.log, prefix...)
}

func (rt *runtimeDso) LevelOutput(level string) Output {
	return LevelOutput(rt.log, level)
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
	if name == "self" {
		rt.self = dispatch
	} else {
		rt.dispatchs[name] = dispatch
	}
}

func (rt *runtimeDso) Make(name string) Dispatch {
	return rt.factories[name](rt)
}

func (rt *runtimeDso) Post(name string, mut *Mutation) {
	if name == "self" {
		rt.self(mut)
	} else {
		rt.dispatchs[name](mut)
	}
}

func (rt *runtimeDso) Log(level string, args ...Any) {
	rt.log(level, args...)
}

func (rt *runtimeDso) Setc(name string, cleaner Cleaner) {
	rt.cleaners[name] = cleaner
}

func (rt *runtimeDso) Getc(name string) Cleaner {
	return rt.cleaners[name]
}

func (rt *runtimeDso) Close() {
	for _, cleaner := range rt.cleaners {
		cleaner.Close()
	}
}
