package lvnbe

import (
	"container/list"
	"fmt"
)

type Hub interface {
	Add(sid string, callback func(mutation *Mutation)) error
	Remove(sid string) error
	Apply(mutation *Mutation) error
	NextId() string
	Close()
}

type hubDso struct {
	id       Id
	state    State
	iterator *list.List
	sessions map[string]*list.Element
}

//ensure iteration order mainly for testing purposes
func NewHub(state State) Hub {
	hub := &hubDso{}
	hub.id = NewId("lv")
	hub.state = state
	hub.iterator = list.New()
	hub.sessions = make(map[string]*list.Element)
	return hub
}

func (hub *hubDso) NextId() string {
	return hub.id.Next()
}

func (hub *hubDso) Add(sid string, callback func(mutation *Mutation)) error {
	_, ok := hub.sessions[sid]
	if ok {
		return fmt.Errorf("duplicated sid %s", sid)
	}
	hub.sessions[sid] = hub.iterator.PushBack(callback)
	mutation := &Mutation{}
	mutation.Sid = sid
	mutation.Name = "all"
	mutation.Args = hub.state.All()
	callback(mutation)
	return nil
}

func (hub *hubDso) Remove(sid string) error {
	element, ok := hub.sessions[sid]
	if !ok {
		return fmt.Errorf("non-existent sid %s", sid)
	}
	hub.iterator.Remove(element)
	delete(hub.sessions, sid)
	hub.call(element, &Mutation{Name: "remove"})
	return nil
}

func (hub *hubDso) Apply(mutation *Mutation) error {
	err := hub.state.Apply(mutation)
	if err != nil {
		return err
	}
	hub.send(mutation)
	return nil
}

func (hub *hubDso) Close() {
	hub.state.Close()
	mutation := &Mutation{Name: "close"}
	hub.send(mutation)
	hub.iterator = nil
	hub.sessions = nil
}

func (hub *hubDso) send(mutation *Mutation) {
	element := hub.iterator.Front()
	for element != nil {
		hub.call(element, mutation)
		element = element.Next()
	}
}

func (hub *hubDso) call(e *list.Element, m *Mutation) {
	e.Value.(Callback)(m)
}
