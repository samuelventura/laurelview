package lvnrt

import (
	"container/list"
	"fmt"
)

type hubSlaveDso struct {
	lrequest  string
	lresponse string
	lerror    string
	callbacks *list.List
	parent    *list.List
	self      *list.Element
}

type hubSessionDso struct {
	callback func(mutation Mutation)
	disposer Action
}

func NewHub(ctx Context) Dispatch {
	log := ctx.PrefixLog("hub")
	dispatchs := make(map[string]Dispatch)
	slaves := make(map[string]*hubSlaveDso)
	buses := make(map[string]*list.List)
	sessions := make(map[string]*hubSessionDso)
	dispatchs[":dispose"] = func(mut Mutation) {
		defer DisposeArgs(mut.Args)
		ClearDispatch(dispatchs)
		for _, session := range sessions {
			session.disposer()
			session.callback(mut)
		}
	}
	dispatchs[":add"] = func(mut Mutation) {
		sid := mut.Sid
		args := mut.Args.(Dispatch)
		_, ok := sessions[sid]
		AssertTrue(!ok, "duplicated sid", sid)
		session := &hubSessionDso{}
		session.callback = args
		session.disposer = NopAction
		sessions[sid] = session
	}
	dispatchs[":remove"] = func(mut Mutation) {
		sid := mut.Sid
		session, ok := sessions[sid]
		if ok { //duplicated cleanup
			session.disposer()
			delete(sessions, sid)
			session.callback(mut)
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
		disposers := make([]Action, 0, len(args))
		statuses := make([]Action, 0, len(args))
		total := NewCount()
		for i, it := range args {
			index := uint(i)
			item := it
			baddr := fmt.Sprintf("%v:%v",
				item.Host, item.Port)
			saddr := fmt.Sprintf("%v:%v:%v",
				item.Host, item.Port, item.Slave)
			slave, ok := slaves[saddr]
			if !ok {
				slave = &hubSlaveDso{}
				slave.callbacks = list.New()
				slaves[saddr] = slave
				parent, ok := buses[baddr]
				if !ok {
					parent = list.New()
					buses[baddr] = parent
				}
				slave.parent = parent
				slave.self = parent.PushBack(slave)
			}
			count := NewCount()
			callback := func(sid string, args StatusArgs) {
				count.Inc()
				total.Inc()
				query := QueryArgs{}
				query.Index = index
				query.Request = args.Request
				query.Response = args.Response
				query.Error = args.Error
				query.Count = count.Get()
				query.Total = total.Get()
				session.callback(Mnsa("query", sid, query))
			}
			args := StatusArgs{}
			args.Address = saddr
			args.Request = slave.lrequest
			args.Response = slave.lresponse
			args.Error = slave.lerror
			statuses = append(statuses, func() {
				callback(sid, args)
			})
			element := slave.callbacks.PushBack(callback)
			disposer := func() {
				slave.callbacks.Remove(element)
				if slave.callbacks.Len() == 0 {
					delete(slaves, saddr)
					slave.parent.Remove(slave.self)
					if slave.parent.Len() == 0 {
						delete(buses, baddr)
					}
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
	status := func(sid string, slave *hubSlaveDso, args StatusArgs) {
		slave.lrequest = args.Request
		slave.lresponse = args.Response
		slave.lerror = args.Error
		element := slave.callbacks.Front()
		for element != nil {
			value := element.Value
			callback := value.(func(sid string, args StatusArgs))
			callback(sid, args)
			element = element.Next()
		}
	}
	dispatchs["status-slave"] = func(mut Mutation) {
		sid := mut.Sid
		args := mut.Args.(StatusArgs)
		slave, ok := slaves[args.Address]
		if ok {
			status(sid, slave, args)
		} else {
			log.Debug(mut)
		}
	}
	dispatchs["status-bus"] = func(mut Mutation) {
		sid := mut.Sid
		args := mut.Args.(StatusArgs)
		parent, ok := buses[args.Address]
		if ok {
			element := parent.Front()
			for element != nil {
				slave := element.Value.(*hubSlaveDso)
				status(sid, slave, args)
				element = element.Next()
			}
		} else {
			log.Debug(mut)
		}
	}
	dispatchs[":ping"] = func(mut Mutation) {
		for _, session := range sessions {
			session.callback(mut)
		}
	}
	return MapDispatch(log, dispatchs)
}
