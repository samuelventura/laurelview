package lvndb

import (
	"container/list"
	"strings"
)

func NewState(rt Runtime) Dispatch {
	dao := rt.GetValue("dao").(Dao)
	hubDispatch := rt.GetDispatch("hub")
	log := PrefixLogger(rt.Log, "state")
	dispatchs := make(map[string]Dispatch)
	iterator := list.New()
	items := make(map[uint]*list.Element)
	for _, item := range dao.All() {
		items[item.ID] = iterator.PushBack(item)
	}
	dispatchs[":dispose"] = func(mut Mutation) {
		defer DisposeArgs(mut.Args)
		ClearDispatch(dispatchs)
		hubDispatch(mut)
	}
	dispatchs[":add"] = func(mut Mutation) {
		hubDispatch(mut)
		list := make([]OneArgs, 0, len(items))
		element := iterator.Front()
		for element != nil {
			item := element.Value.(ItemDro)
			one := OneArgs{}
			one.Id = item.ID
			one.Name = item.Name
			one.Json = item.Json
			list = append(list, one)
			element = element.Next()
		}
		all := Mnsa("all", mut.Sid, list)
		hubDispatch(all)
	}
	dispatchs[":remove"] = func(mut Mutation) {
		hubDispatch(mut)
	}
	dispatchs["create"] = func(mut Mutation) {
		args := mut.Args.(OneArgs)
		args.Name = strings.TrimSpace(args.Name)
		if len(args.Name) == 0 {
			log.Debug(mut, "name cannot be empty")
			return
		}
		item := dao.Create(args.Name, args.Json)
		items[item.ID] = iterator.PushBack(item)
		nmut := mut
		nmut.Args = OneArgs{
			Id:   item.ID,
			Name: item.Name,
			Json: item.Json,
		}
		hubDispatch(nmut)
	}
	dispatchs["delete"] = func(mut Mutation) {
		id := mut.Args.(uint)
		element, ok := items[id]
		if !ok {
			log.Debug(mut, "item not found", id)
			return
		}
		dao.Delete(id)
		delete(items, id)
		iterator.Remove(element)
		hubDispatch(mut)
	}
	dispatchs["update"] = func(mut Mutation) {
		args := mut.Args.(OneArgs)
		_, ok := items[args.Id]
		if !ok {
			log.Debug(mut, "item not found", args.Id)
			return
		}
		args.Name = strings.TrimSpace(args.Name)
		if len(args.Name) == 0 {
			log.Debug(mut, "name cannot be empty")
			return
		}
		item := dao.Update(args.Id, args.Name, args.Json)
		items[item.ID].Value = item
		hubDispatch(mut)
	}
	return MapDispatch(log, dispatchs)
}
