package lvnrt

import (
	"container/list"
	"fmt"
)

type hubDso struct {
	output   Output
	disposed bool
	slaves   map[string]*hubSlaveDso
	sessions map[string]*hubSessionDso
}

type hubSlaveDso struct {
	request   string
	response  string
	callbacks *list.List
}

type hubSessionDso struct {
	callback func(mutation *Mutation)
	disposer Action
}

func NewHub(output Output) Dispatch {
	hub := &hubDso{}
	hub.output = output
	hub.slaves = make(map[string]*hubSlaveDso)
	hub.sessions = make(map[string]*hubSessionDso)
	return hub.apply
}

func (hub *hubDso) apply(mut *Mutation) {
	err := hub.applyMutation(mut)
	traceIfError(hub.output, err)
}

func (hub *hubDso) applyMutation(mut *Mutation) error {
	switch mut.Name {
	case "status":
		return hub.applyStatus(mut.Sid, mut.Args.(*StatusArgs))
	case "setup":
		return hub.applySetup(mut.Sid, mut.Args.(*SetupArgs))
	case "add":
		return hub.applyAdd(mut.Sid, mut.Args.(*AddArgs))
	case "remove":
		return hub.applyRemove(mut.Sid, mut.Args.(*RemoveArgs))
	case "dispose":
		return hub.applyDispose(mut.Sid, mut.Args.(*DisposeArgs))
	}
	return fmt.Errorf("unknown mutation %v", mut.Name)
}

func (hub *hubDso) applyAdd(sid string, args *AddArgs) error {
	_, ok := hub.sessions[sid]
	if ok {
		return fmt.Errorf("duplicated sid %v", sid)
	}
	session := &hubSessionDso{}
	session.callback = args.Callback
	session.disposer = NopAction
	hub.sessions[sid] = session
	return nil
}

func (hub *hubDso) applyRemove(sid string, args *RemoveArgs) error {
	session, ok := hub.sessions[sid]
	if !ok {
		return fmt.Errorf("non-existent sid %v", sid)
	}
	session.disposer()
	delete(hub.sessions, sid)
	mut := &Mutation{}
	mut.Sid = sid
	mut.Name = "remove"
	mut.Args = &RemoveArgs{}
	session.callback(mut)
	return nil
}

func (hub *hubDso) applyDispose(sid string, args *DisposeArgs) error {
	if hub.disposed {
		return fmt.Errorf("already disposed")
	}
	hub.disposed = true
	for _, session := range hub.sessions {
		session.disposer()
		mut := &Mutation{}
		mut.Sid = sid
		mut.Name = "dispose"
		mut.Args = &DisposeArgs{}
		session.callback(mut)
	}
	return nil
}

func (hub *hubDso) applySetup(sid string, args *SetupArgs) error {
	session, ok := hub.sessions[sid]
	if !ok {
		return fmt.Errorf("non-existent sid %v", sid)
	}
	session.disposer()
	disposers := make([]Action, 0, len(args.Items))
	statuses := make([]Action, 0, len(args.Items))
	for i, it := range args.Items {
		index := uint(i)
		item := it
		address := fmt.Sprintf("%v:%v:%v",
			item.Host, item.Port, item.Slave)
		slave, ok := hub.slaves[address]
		if !ok {
			slave = &hubSlaveDso{}
			slave.callbacks = list.New()
			hub.slaves[address] = slave
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
				delete(hub.slaves, address)
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
	return nil
}

func (hub *hubDso) applyStatus(sid string, args *StatusArgs) error {
	slave, ok := hub.slaves[args.Slave]
	if !ok {
		return fmt.Errorf("non-existent slave %v", args.Slave)
	}
	slave.request = args.Request
	slave.response = args.Response
	element := slave.callbacks.Front()
	for element != nil {
		value := element.Value
		callback := value.(func(sid string, args *StatusArgs))
		callback(sid, args)
		element = element.Next()
	}
	return nil
}
