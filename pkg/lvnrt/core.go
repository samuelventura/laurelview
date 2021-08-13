package lvnrt

type coreDso struct {
	dispatch Dispatch
	disposed bool
	output   Output
	queue    Queue
}

func NewCoreHub(output Output, dispatch Dispatch) Dispatch {
	core := &coreDso{}
	core.dispatch = dispatch
	core.output = output
	core.queue = make(Queue)
	go core.loop()
	return core.apply
}

func (core *coreDso) apply(mut *Mutation) {
	core.queue <- func() {
		defer core.dispose(mut)
		defer traceRecover(core.output)
		core.dispatch(mut)
	}
}

func (core *coreDso) dispose(mut *Mutation) {
	if mut.Name == "dispose" {
		core.disposed = true
	}
}

func (core *coreDso) loop() {
	for action := range core.queue {
		if core.disposed {
			return
		}
		core.run(action)
	}
}

func (core *coreDso) run(action Action) {
	defer traceRecover(core.output)
	action()
}
