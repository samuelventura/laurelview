package lvnrt

import (
	"encoding/json"
	"fmt"
	"math/bits"
	"strconv"
)

func decodeMutation(bytes []byte) (mut *Mutation, err error) {
	return decodeMutationEx(bytes, false)
}

func decodeMutationEx(bytes []byte, ws bool) (mut *Mutation, err error) {
	mut = &Mutation{}
	var mmi Any
	err = json.Unmarshal(bytes, &mmi)
	if err != nil {
		return
	}
	mm := mmi.(Map)
	mut.Name = mm["name"].(string)
	if ws {
		mut.Sid = mm["sid"].(string)
	}
	switch mut.Name {
	//FIXME validate inputs
	case "setup":
		argm := mm["args"].(Map)
		args := &SetupArgs{}
		items := argm["items"].([]Any)
		args.Items = make([]*ItemArgs, 0, len(items))
		for _, imi := range items {
			fm := imi.(Map)
			carg := &ItemArgs{}
			carg.Host = fm["host"].(string)
			carg.Port = parseUint(fm["port"])
			carg.Slave = parseUint(fm["slave"])
			args.Items = append(args.Items, carg)
		}
		mut.Args = args
	case "query":
		argm := mm["args"].(Map)
		args := &QueryArgs{}
		args.Index = parseUint(argm["index"])
		args.Request = argm["request"].(string)
		args.Response = maybeString(argm["response"])
		args.Count = maybeUint(argm["count"])
		args.Total = maybeUint(argm["total"])
		mut.Args = args
	default:
		err = fmt.Errorf("unkown mutation %s", mut.Name)
	}
	return
}

func maybeUint(any Any) uint {
	switch v := any.(type) {
	case float64:
		return uint(v)
	default:
		return 0
	}
}

func maybeString(any Any) string {
	switch v := any.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func parseUint(any Any) uint {
	switch v := any.(type) {
	case float64:
		return uint(v)
	case string:
		id, err := strconv.ParseUint(v, 10, bits.UintSize)
		PanicIfError(err)
		return uint(id)
	default:
		return 0
	}
}

func encodeMutation(mutation *Mutation) []byte {
	mm := make(Map)
	mm["name"] = mutation.Name
	mm["sid"] = mutation.Sid
	args, err := encodeArgs(mutation.Name, mutation.Args)
	PanicIfError(err)
	mm["args"] = args
	bytes, err := json.Marshal(mm)
	PanicIfError(err)
	return bytes
}

func encodeArgs(name string, argi Any) (argm Map, err error) {
	switch name {
	case "setup":
		args := argi.(*SetupArgs)
		argm = make(Map)
		items := make([]Map, 0, len(args.Items))
		for _, item := range args.Items {
			fm := make(Map)
			fm["host"] = item.Host
			fm["port"] = item.Port
			fm["slave"] = item.Slave
			items = append(items, fm)
		}
		argm["items"] = items
	case "query":
		args := argi.(*QueryArgs)
		argm = make(Map)
		argm["index"] = args.Index
		argm["request"] = args.Request
		argm["response"] = args.Response
		argm["count"] = args.Count
		argm["total"] = args.Total
	default:
		err = fmt.Errorf("unkown mutation %s", name)
	}
	return
}
