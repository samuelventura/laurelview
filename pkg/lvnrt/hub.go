package lvnrt

import (
	"container/list"
	"fmt"
)

type hubSlaveDso struct {
	request   string
	response  string
	callbacks *list.List
}

type hubSessionDso struct {
	callback func(mutation *Mutation)
	disposer Action
}

func NewHub(rt Runtime) Dispatch {
	log := prefixLogger(rt.Log, "hub")
	dispatchs := make(map[string]Dispatch)
	slaves := make(map[string]*hubSlaveDso)
	sessions := make(map[string]*hubSessionDso)
	dispatchs["dispose"] = func(mut *Mutation) {
		defer disposeArgs(mut.Args)
		clearDispatch(dispatchs)
		for _, session := range sessions {
			session.disposer()
			session.callback(mut)
		}
	}
	dispatchs["add"] = func(mut *Mutation) {
		sid := mut.Sid
		args := mut.Args.(*AddArgs)
		_, ok := sessions[sid]
		assertTrue(!ok, "duplicated sid", sid)
		session := &hubSessionDso{}
		session.callback = args.Callback
		session.disposer = NopAction
		sessions[sid] = session
	}
	dispatchs["remove"] = func(mut *Mutation) {
		sid := mut.Sid
		session, ok := sessions[sid]
		assertTrue(ok, "non-existent sid", sid)
		session.disposer()
		delete(sessions, sid)
		session.callback(mut)
	}
	dispatchs["setup"] = func(mut *Mutation) {
		sid := mut.Sid
		args := mut.Args.(*SetupArgs)
		session, ok := sessions[sid]
		assertTrue(ok, "non-existent sid", sid)
		session.disposer()
		disposers := make([]Action, 0, len(args.Items))
		statuses := make([]Action, 0, len(args.Items))
		for i, it := range args.Items {
			index := uint(i)
			item := it
			address := fmt.Sprintf("%v:%v:%v",
				item.Host, item.Port, item.Slave)
			slave, ok := slaves[address]
			if !ok {
				slave = &hubSlaveDso{}
				slave.callbacks = list.New()
				slaves[address] = slave
			}
			callback := func(sid string, args *StatusArgs) {
				mut := &Mutation{}
				mut.Sid = sid
				mut.Name = "query"
				query := &QueryArgs{}
				query.Index = index
				query.Request = args.Request
				query.Response = args.Response
				mut.Args = query
				session.callback(mut)
			}
			args := &StatusArgs{}
			args.Slave = address
			args.Request = slave.request
			args.Response = slave.response
			statuses = append(statuses, func() {
				callback(sid, args)
			})
			element := slave.callbacks.PushBack(callback)
			disposer := func() {
				slave.callbacks.Remove(element)
				if slave.callbacks.Len() == 0 {
					delete(slaves, address)
				}
			}
			disposers = append(disposers, disposer)
		}
		session.disposer = func() {
			for _, disposer := range disposers {
				disposer()
			}
		}
		for _, action := range statuses {
			action()
		}
	}
	dispatchs["status"] = func(mut *Mutation) {
		sid := mut.Sid
		args := mut.Args.(*StatusArgs)
		slave, ok := slaves[args.Slave]
		assertTrue(ok, "non-existent slave", args.Slave)
		slave.request = args.Request
		slave.response = args.Response
		element := slave.callbacks.Front()
		for element != nil {
			value := element.Value
			callback := value.(func(sid string, args *StatusArgs))
			callback(sid, args)
			element = element.Next()
		}
	}
	return mapDispatch(log.Trace, log.Debug, dispatchs)
}
