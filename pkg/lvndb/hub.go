package lvndb

import (
	"container/list"
)

//ensure iteration order mainly for testing purposes
func NewHub(rt Runtime) Dispatch {
	log := PrefixLogger(rt.Log, "hub")
	dispatchs := make(map[string]Dispatch)
	iterator := list.New()
	sessions := make(map[string]*list.Element)
	call := func(e *list.Element, mut Mutation) {
		dispatch := e.Value.(Dispatch)
		dispatch(mut)
	}
	sendall := func(mut Mutation) {
		element := iterator.Front()
		for element != nil {
			call(element, mut)
			element = element.Next()
		}
	}
	dispatchs[":dispose"] = func(mut Mutation) {
		defer DisposeArgs(mut.Args)
		ClearDispatch(dispatchs)
		sendall(mut)
		sessions = nil
		iterator = nil
	}
	dispatchs[":add"] = func(mut Mutation) {
		sid := mut.Sid
		_, ok := sessions[sid]
		if ok {
			log.Debug(mut, "duplicated sid", sid)
			return
		}
		callback := mut.Args.(Dispatch)
		sessions[sid] = iterator.PushBack(callback)
	}
	dispatchs[":remove"] = func(mut Mutation) {
		sid := mut.Sid
		element, ok := sessions[sid]
		if !ok {
			log.Debug(mut, "non-existent sid", sid)
			return
		}
		iterator.Remove(element)
		delete(sessions, sid)
		call(element, mut)
	}
	dispatchs["create"] = func(mut Mutation) {
		sendall(mut)
	}
	dispatchs["update"] = func(mut Mutation) {
		sendall(mut)
	}
	dispatchs["delete"] = func(mut Mutation) {
		sendall(mut)
	}
	dispatchs["all"] = func(mut Mutation) {
		sid := mut.Sid
		element, ok := sessions[sid]
		if !ok {
			log.Debug(mut, "non-existent sid", sid)
			return
		}
		call(element, mut)
	}
	return MapDispatch(log, dispatchs)
}
