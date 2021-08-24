package lvsdk

func DefaultRuntime() Runtime {
	return NewRuntime(DefaultLog())
}

type runtimeDso struct {
	log       Log
	logger    Logger
	cleaner   Cleaner
	closed    Channel
	ids       map[string]Id
	values    map[string]Any
	factories map[string]Factory
	dispatchs map[string]Dispatch
	cleaners  map[string]Cleaner
}

func NewRuntime(log Log) Runtime {
	return NewRuntimeId(log, "drt")
}

func NewRuntimeId(log Log, id string) Runtime {
	rt := &runtimeDso{}
	rt.log = log
	rt.logger = PrefixLogger(log, id)
	rt.cleaner = NewCleaner(rt.logger)
	rt.closed = make(Channel)
	rt.values = make(map[string]Any)
	rt.factories = make(map[string]Factory)
	rt.dispatchs = make(map[string]Dispatch)
	rt.cleaners = make(map[string]Cleaner)
	return rt
}

func (rt *runtimeDso) PrefixLog(prefix ...Any) Logger {
	return PrefixLogger(rt.log, prefix...)
}

func (rt *runtimeDso) LevelOutput(level string) Output {
	return LevelOutput(rt.log, level)
}

func (rt *runtimeDso) SetId(name string, id Id) {
	rt.ids[name] = id
}

func (rt *runtimeDso) SetValue(name string, value Any) {
	rt.values[name] = value
}

func (rt *runtimeDso) SetFactory(name string, factory Factory) {
	rt.factories[name] = factory
}

func (rt *runtimeDso) SetDispatch(name string, dispatch Dispatch) {
	rt.dispatchs[name] = dispatch
}

func (rt *runtimeDso) SetCleaner(name string, cleaner Cleaner) {
	rt.cleaners[name] = cleaner
	rt.cleaner.AddAction(name, cleaner.Close)
}

func (rt *runtimeDso) GetId(name string) Id {
	id, ok := rt.ids[name]
	if ok {
		return id
	} else {
		PanicLN("id not found", name)
		return nil
	}
}

func (rt *runtimeDso) GetValue(name string) Any {
	value, ok := rt.values[name]
	if ok {
		return value
	} else {
		PanicLN("value not found", name)
		return nil
	}
}

func (rt *runtimeDso) GetCleaner(name string) Cleaner {
	cleaner, ok := rt.cleaners[name]
	if ok {
		return cleaner
	} else {
		PanicLN("cleaner not found", name)
		return nil
	}
}

func (rt *runtimeDso) GetFactory(name string) Factory {
	fact, ok := rt.factories[name]
	if ok {
		return fact
	} else {
		PanicLN("factory not found", name)
		return nil
	}
}

func (rt *runtimeDso) GetDispatch(name string) Dispatch {
	disp, ok := rt.dispatchs[name]
	if ok {
		return disp
	} else {
		PanicLN("dispatch not found", name)
		return nil
	}
}

func (rt *runtimeDso) Log(level string, args ...Any) {
	rt.log(level, args...)
}

func (rt *runtimeDso) Closed() Channel {
	return rt.closed
}

func (rt *runtimeDso) Close() {
	close(rt.closed)
	rt.cleaner.Close()
}
