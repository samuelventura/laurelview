package lvnrt

import (
	"fmt"
)

type stateBusDso struct {
	dispatch Dispatch
	//json.slave vs ref count
	slaves map[uint]Count
}

type stateSessionDso struct {
	//index vs json.slave
	slaves map[uint]uint
	//index vs busDso
	buses    map[uint]*stateBusDso
	disposer Action
	setup    Flag
}

func NewState(rt Runtime) Dispatch {
	log := PrefixLogger(rt.Log, "state")
	hubDispatch := rt.GetDispatch("hub")
	busFactory := rt.GetFactory("bus")
	dispatchs := make(map[string]Dispatch)
	sessions := make(map[string]*stateSessionDso)
	buses := make(map[string]*stateBusDso)
	dispatchs[":dispose"] = func(mut Mutation) {
		ClearDispatch(dispatchs)
		for sid, session := range sessions {
			session.disposer()
			delete(sessions, sid)
		}
		DisposeArgs(mut.Args)
		hubDispatch(mut)
	}
	dispatchs[":add"] = func(mut Mutation) {
		sid := mut.Sid
		_, ok := sessions[sid]
		AssertTrue(!ok, "duplicated sid", sid)
		session := &stateSessionDso{}
		session.disposer = NopAction
		session.buses = make(map[uint]*stateBusDso)
		session.slaves = make(map[uint]uint)
		session.setup = NewFlag()
		sessions[sid] = session
		hubDispatch(mut)
	}
	dispatchs[":remove"] = func(mut Mutation) {
		sid := mut.Sid
		session, ok := sessions[sid]
		AssertTrue(ok, "non-existent sid", sid)
		session.disposer()
		delete(sessions, sid)
		hubDispatch(mut)
	}
	dispatchs["setup"] = func(mut Mutation) {
		sid := mut.Sid
		args := mut.Args.([]ItemArgs)
		session, ok := sessions[sid]
		AssertTrue(ok, "non-existent sid", sid)
		AssertTrue(!session.setup.Get(), "re-setup sid", sid)
		session.setup.Set(true)
		disposers := make([]Action, 0, len(args))
		for i, it := range args {
			index := uint(i)
			item := it
			address := fmt.Sprintf("%v:%v", item.Host, item.Port)
			bus, ok := buses[address]
			if !ok {
				bus = &stateBusDso{}
				bus.dispatch = busFactory(rt)
				bus.slaves = make(map[uint]Count)
				bus.dispatch(Mnsa("setup", sid, address))
				buses[address] = bus
			}
			session.buses[index] = bus
			session.slaves[index] = item.Slave
			count, ok := bus.slaves[item.Slave]
			if !ok {
				count = NewCount()
				bus.slaves[item.Slave] = count
			}
			count.Inc()
			args := SlaveArgs{}
			args.Slave = item.Slave
			args.Count = count.Get()
			bus.dispatch(Mnsa("slave", sid, args))
			disposer := func() {
				count.Dec()
				args := SlaveArgs{}
				args.Slave = item.Slave
				args.Count = count.Get()
				bus.dispatch(Mnsa("slave", sid, args))
				if count.Get() == 0 {
					delete(bus.slaves, item.Slave)
				}
				if len(bus.slaves) == 0 {
					delete(buses, address)
					bus.dispatch(Mns(":dispose", sid))
				}
			}
			disposers = append(disposers, disposer)
		}
		session.disposer = func() {
			for _, disposer := range disposers {
				disposer()
			}
		}
		hubDispatch(mut)
	}
	dispatchs["query"] = func(mut Mutation) {
		sid := mut.Sid
		args := mut.Args.(QueryArgs)
		session, ok := sessions[sid]
		AssertTrue(ok, "non-existent sid", sid, mut)
		AssertTrue(session.setup.Get(), "non-setup sid", sid)
		bus, ok := session.buses[args.Index]
		AssertTrue(ok, "non-existent bus", args.Index, mut)
		slave, ok := session.slaves[args.Index]
		AssertTrue(ok, "non-existent slave", args.Index, mut)
		nargs := args
		nmut := mut
		nargs.Index = slave
		nmut.Args = nargs
		bus.dispatch(nmut)
	}
	return MapDispatch(log, dispatchs)
}
