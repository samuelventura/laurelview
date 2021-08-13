package lvnrt

import (
	"fmt"
)

type stateChannelDso struct {
	dispatch Dispatch
	slaves   map[uint]Count
	count    int
}

type stateSessionDso struct {
	slaves   map[uint]uint
	channels map[uint]*stateChannelDso
	disposer Action
}

type stateDso struct {
	output   Output
	disposed bool
	factory  Factory
	dispatch Dispatch
	channels map[string]*stateChannelDso
	sessions map[string]*stateSessionDso
}

func NewState(output Output, dispatch Dispatch, factory Factory) Dispatch {
	state := &stateDso{}
	state.output = output
	state.factory = factory
	state.dispatch = dispatch
	state.channels = make(map[string]*stateChannelDso)
	state.sessions = make(map[string]*stateSessionDso)
	return state.apply
}

func (state *stateDso) apply(mut *Mutation) {
	err := state.applyMutation(mut)
	traceIfError(state.output, err)
}

func (state *stateDso) applyMutation(mut *Mutation) error {
	switch mut.Name {
	case "query":
		return state.applyQuery(mut.Sid, mut.Args.(*QueryArgs))
	case "setup":
		defer state.dispatch(mut)
		return state.applySetup(mut.Sid, mut.Args.(*SetupArgs))
	case "add":
		defer state.dispatch(mut)
		return state.applyAdd(mut.Sid, mut.Args.(*AddArgs))
	case "remove":
		defer state.dispatch(mut)
		return state.applyRemove(mut.Sid, mut.Args.(*RemoveArgs))
	case "dispose":
		defer state.dispatch(mut)
		return state.applyDispose(mut.Args.(*DisposeArgs))
	}
	return fmt.Errorf("unknown mutation %v", mut.Name)
}

func (state *stateDso) applyAdd(sid string, args *AddArgs) error {
	_, ok := state.sessions[sid]
	if ok {
		return fmt.Errorf("duplicated sid %v", sid)
	}
	session := &stateSessionDso{}
	session.disposer = NopAction
	session.channels = make(map[uint]*stateChannelDso)
	session.slaves = make(map[uint]uint)
	state.sessions[sid] = session
	return nil
}

func (state *stateDso) applyRemove(sid string, args *RemoveArgs) error {
	session, ok := state.sessions[sid]
	if !ok {
		return fmt.Errorf("non-existent sid %v", sid)
	}
	session.disposer()
	delete(state.sessions, sid)
	return nil
}

func (state *stateDso) applyDispose(args *DisposeArgs) error {
	if state.disposed {
		return fmt.Errorf("already disposed")
	}
	state.disposed = true
	for sid, session := range state.sessions {
		session.disposer()
		delete(state.sessions, sid)
	}
	return nil
}

func (state *stateDso) applySetup(sid string, args *SetupArgs) error {
	session, ok := state.sessions[sid]
	if !ok {
		return fmt.Errorf("non-existent sid %v", sid)
	}
	session.disposer()
	session.channels = make(map[uint]*stateChannelDso)
	session.slaves = make(map[uint]uint)
	disposers := make([]Action, 0, len(args.Items))
	for i, it := range args.Items {
		index := uint(i)
		item := it
		address := fmt.Sprintf("%v:%v",
			item.Host, item.Port)
		channel, ok := state.channels[address]
		if !ok {
			channel = &stateChannelDso{}
			channel.dispatch = state.factory(nil)
			channel.slaves = make(map[uint]Count)
			args := &BusArgs{}
			args.Host = item.Host
			args.Port = item.Port
			mut := &Mutation{}
			mut.Sid = sid
			mut.Name = "bus"
			mut.Args = args
			channel.dispatch(mut)
			state.channels[address] = channel
		}
		channel.count++
		session.channels[index] = channel
		session.slaves[index] = item.Slave
		count, ok := channel.slaves[item.Slave]
		if !ok {
			count = NewCount()
			channel.slaves[item.Slave] = count
		}
		count.Inc()
		args := &SlaveArgs{}
		args.Slave = item.Slave
		args.Count = count.Count()
		mut := &Mutation{}
		mut.Sid = sid
		mut.Name = "slave"
		mut.Args = args
		channel.dispatch(mut)
		disposer := func() {
			count.Dec()
			args := &SlaveArgs{}
			args.Slave = item.Slave
			args.Count = count.Count()
			mut := &Mutation{}
			mut.Sid = sid
			mut.Name = "slave"
			mut.Args = args
			channel.dispatch(mut)
			if count.Count() == 0 {
				delete(channel.slaves, item.Slave)
			}
			if channel.count == 0 {
				delete(state.channels, address)
				mut := &Mutation{}
				mut.Sid = sid
				mut.Name = "dispose"
				channel.dispatch(mut)
			}
		}
		disposers = append(disposers, disposer)
	}
	session.disposer = func() {
		for _, disposer := range disposers {
			disposer()
		}
	}
	return nil
}

func (state *stateDso) applyQuery(sid string, args *QueryArgs) error {
	session, ok := state.sessions[sid]
	if !ok {
		return fmt.Errorf("non-existent sid %v", sid)
	}
	channel, ok := session.channels[args.Index]
	if !ok {
		return fmt.Errorf("non-existent channel %v", args.Index)
	}
	slave, ok := session.slaves[args.Index]
	if !ok {
		return fmt.Errorf("non-existent slave %v", args.Index)
	}
	query := &QueryArgs{}
	query.Index = slave
	query.Request = args.Request
	query.Response = args.Response
	mut := &Mutation{}
	mut.Sid = sid
	mut.Name = "query"
	mut.Args = query
	channel.dispatch(mut)
	return nil
}
