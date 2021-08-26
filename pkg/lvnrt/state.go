package lvnrt

import (
	"fmt"
)

type stateBusDso struct {
	dispatch Dispatch
	slaves   map[uint]Count
}

type stateSessionDso struct {
	slaves   map[uint]uint
	buses    map[uint]*stateBusDso
	disposer Action
}

func NewState(rt Runtime) Dispatch {
	log := PrefixLogger(rt.Log, "state")
	hubDispatch := rt.GetDispatch("hub")
	dispatchs := make(map[string]Dispatch)
	sessions := make(map[string]*stateSessionDso)
	buses := make(map[string]*stateBusDso)
	dispatchs[":dispose"] = func(mut Mutation) {
		defer DisposeArgs(mut.Args)
		ClearDispatch(dispatchs)
		for sid, session := range sessions {
			session.disposer()
			delete(sessions, sid)
		}
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
		sessions[sid] = session
		hubDispatch(mut)
	}
	dispatchs[":remove"] = func(mut Mutation) {
		sid := mut.Sid
		session, ok := sessions[sid]
		if ok { //duplicate cleanup
			session.disposer()
			delete(sessions, sid)
			hubDispatch(mut)
		} else {
			log.Debug(mut)
		}
	}
	dispatchs["setup"] = func(mut Mutation) {
		sid := mut.Sid
		args := mut.Args.([]ItemArgs)
		session, ok := sessions[sid]
		AssertTrue(ok, "non-existent sid", sid)
		session.disposer()
		session.buses = make(map[uint]*stateBusDso)
		session.slaves = make(map[uint]uint)
		disposers := make([]Action, 0, len(args))
		for i, it := range args {
			index := uint(i)
			item := it
			address := fmt.Sprintf("%v:%v",
				item.Host, item.Port)
			bus, ok := buses[address]
			if !ok {
				bus = &stateBusDso{}
				bus.dispatch = rt.GetFactory("bus")(rt)
				bus.slaves = make(map[uint]Count)
				args := BusArgs{}
				args.Host = item.Host
				args.Port = item.Port
				mut := Mutation{}
				mut.Sid = sid
				mut.Name = "setup"
				mut.Args = args
				bus.dispatch(mut)
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
			args.Count = count.Count()
			mut := Mutation{}
			mut.Sid = sid
			mut.Name = "slave"
			mut.Args = args
			bus.dispatch(mut)
			disposer := func() {
				count.Dec()
				args := SlaveArgs{}
				args.Slave = item.Slave
				args.Count = count.Count()
				mut := Mutation{}
				mut.Sid = sid
				mut.Name = "slave"
				mut.Args = args
				bus.dispatch(mut)
				if count.Count() == 0 {
					delete(bus.slaves, item.Slave)
				}
				if len(bus.slaves) == 0 {
					delete(buses, address)
					mut := Mutation{}
					mut.Sid = sid
					mut.Name = ":dispose"
					bus.dispatch(mut)
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
		AssertTrue(ok, "non-existent sid", sid)
		bus, ok := session.buses[args.Index]
		AssertTrue(ok, "non-existent bus", args.Index)
		slave, ok := session.slaves[args.Index]
		AssertTrue(ok, "non-existent slave", args.Index)
		nargs := args
		nmut := mut
		nargs.Index = slave
		nmut.Args = nargs
		bus.dispatch(nmut)
	}
	return MapDispatch(log, dispatchs)
}
