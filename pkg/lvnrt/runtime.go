package lvnrt

//provides
//1) config
//2) factories
//3) dispatchs
//4) output
type Runtime interface {
	Getv(name string) Any
	Setv(name string, value Any)
	Setf(name string, factory Factory)
	Setd(name string, dispatch Dispatch)
	Make(name string) Dispatch
	Post(name string, mut *Mutation)
	Output(args ...Any)
}

type runtimeDso struct {
	output    Output
	values    map[string]Any
	factories map[string]Factory
	dispatchs map[string]Dispatch
}

func NewRuntime(output Output) Runtime {
	rt := &runtimeDso{}
	rt.output = output
	rt.factories = make(map[string]Factory)
	rt.dispatchs = make(map[string]Dispatch)
	return rt
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
	rt.dispatchs[name](mut)
}

func (rt *runtimeDso) Output(args ...Any) {
	rt.output(args...)
}
