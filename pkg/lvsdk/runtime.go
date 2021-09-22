package lvsdk

func DefaultContext() Context {
	return NewContext(DefaultLog())
}

type contextDso struct {
	log       Log
	logger    Logger
	cleaner   Cleaner
	values    map[string]Any
	factories map[string]Factory
	dispatchs map[string]Dispatch
}

//pass a prefixed log instead of assigning an id
func NewContext(log Log) Context {
	ctx := &contextDso{}
	ctx.log = log
	ctx.logger = PrefixLogger(log)
	clogger := PrefixLogger(log, "cleaner")
	ctx.cleaner = NewCleaner(clogger)
	ctx.values = make(map[string]Any)
	ctx.factories = make(map[string]Factory)
	ctx.dispatchs = make(map[string]Dispatch)
	return ctx
}

func (ctx *contextDso) Cleaner() Cleaner {
	return ctx.cleaner
}

func (ctx *contextDso) PrefixLog(prefix ...Any) Logger {
	return PrefixLogger(ctx.log, prefix...)
}

func (ctx *contextDso) LevelOutput(level string) Output {
	return LevelOutput(ctx.log, level)
}

func (ctx *contextDso) SetValue(name string, value Any) {
	ctx.values[name] = value
}

func (ctx *contextDso) SetFactory(name string, factory Factory) {
	ctx.factories[name] = factory
}

func (ctx *contextDso) SetDispatch(name string, dispatch Dispatch) {
	ctx.dispatchs[name] = dispatch
}

func (ctx *contextDso) GetValue(name string) Any {
	value, ok := ctx.values[name]
	if ok {
		return value
	} else {
		PanicLN("value not found", name)
		return nil
	}
}

func (ctx *contextDso) GetFactory(name string) Factory {
	fact, ok := ctx.factories[name]
	if ok {
		return fact
	} else {
		PanicLN("factory not found", name)
		return nil
	}
}

func (ctx *contextDso) GetDispatch(name string) Dispatch {
	disp, ok := ctx.dispatchs[name]
	if ok {
		return disp
	} else {
		PanicLN("dispatch not found", name)
		return nil
	}
}

func (ctx *contextDso) Log(level string, args ...Any) {
	ctx.log(level, args...)
}

func (ctx *contextDso) Close() Channel {
	ctx.cleaner.Close()
	done := make(Channel)
	ctx.cleaner.AddChannel("done", done)
	return done
}
