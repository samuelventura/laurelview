package lvnbe

import (
	"container/list"
	"fmt"
)

type State interface {
	All() *AllArgs
	Apply(mutation *Mutation) error
	Close()
}

type stateDso struct {
	dao      Dao
	iterator *list.List
	items    map[uint]*list.Element
}

func NewState(dao Dao) State {
	state := &stateDso{}
	state.dao = dao
	state.iterator = list.New()
	state.items = make(map[uint]*list.Element)
	for _, item := range state.dao.All() {
		clone := item //decouple
		element := state.iterator.PushBack(&clone)
		state.items[item.ID] = element
	}
	return state
}

func (state *stateDso) Close() {
	state.dao.Close()
}

func (state *stateDso) Apply(mutation *Mutation) error {
	switch mutation.Name {
	case "create":
		return state.applyCreate(mutation.Args.(*CreateArgs))
	case "delete":
		return state.applyDelete(mutation.Args.(*DeleteArgs))
	case "update":
		return state.applyUpdate(mutation.Args.(*UpdateArgs))
	}
	return fmt.Errorf("unknown mutation %v", mutation.Name)
}

func (state *stateDso) All() *AllArgs {
	all := &AllArgs{}
	all.Items = make([]*OneArgs, 0, len(state.items))
	element := state.iterator.Front()
	for element != nil {
		item := element.Value.(*ItemDro)
		mut := &OneArgs{}
		mut.Id = item.ID
		mut.Name = item.Name
		mut.Json = item.Json
		all.Items = append(all.Items, mut)
		element = element.Next()
	}
	return all
}

func (state *stateDso) applyCreate(args *CreateArgs) error {
	item := state.dao.Create(args.Name, args.Json)
	clone := *item //decouple
	element := state.iterator.PushBack(&clone)
	state.items[item.ID] = element
	args.Id = item.ID
	return nil
}

func (state *stateDso) applyDelete(args *DeleteArgs) error {
	element, ok := state.items[args.Id]
	if !ok {
		return fmt.Errorf("unknown item %v", args.Id)
	}
	state.dao.Delete(args.Id)
	delete(state.items, args.Id)
	state.iterator.Remove(element)
	return nil
}

func (state *stateDso) applyUpdate(args *UpdateArgs) error {
	if _, ok := state.items[args.Id]; !ok {
		return fmt.Errorf("unknown item %v", args.Id)
	}
	item := state.dao.Update(args.Id, args.Name, args.Json)
	state.items[item.ID].Value = item
	return nil
}
