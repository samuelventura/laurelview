package lvnbe

type coreDso struct {
	hub    Hub
	output Output
	queue  chan func()
}

type Core interface {
	Add(sid string, callback func(mutation *Mutation))
	Remove(sid string)
	Apply(mutation *Mutation)
	NextId() string
	Close()
}

func NewCore(hub Hub, output Output) *coreDso {
	core := &coreDso{}
	core.hub = hub
	core.output = output
	core.queue = make(Queue)
	go core.loop()
	return core
}

func (core *coreDso) NextId() string {
	return core.hub.NextId()
}

func (core *coreDso) Add(sid string, callback func(mutation *Mutation)) {
	core.queue <- func() {
		err := core.hub.Add(sid, callback)
		TraceIfError(core.output, err)
	}
}

func (core *coreDso) Remove(sid string) {
	core.queue <- func() {
		err := core.hub.Remove(sid)
		TraceIfError(core.output, err)
	}
}

func (core *coreDso) Apply(mutation *Mutation) {
	core.queue <- func() {
		err := core.hub.Apply(mutation)
		TraceIfError(core.output, err)
	}
}

func (core *coreDso) Close() {
	//chained close to avoid race condition
	//core -> hub -> state -> dao
	core.queue <- func() {
		core.hub.Close()
		close(core.queue)
	}
}

func (core *coreDso) loop() {
	for action := range core.queue {
		core.run(action)
	}
}

func (core *coreDso) run(action Action) {
	defer TraceRecover(core.output)
	action()
}
