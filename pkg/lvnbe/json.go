package lvnbe

import (
	"encoding/json"
	"fmt"
	"math/bits"
	"strconv"
)

func decodeMutation(bytes []byte) (mut *Mutation, err error) {
	mut = &Mutation{}
	var mmi Any
	err = json.Unmarshal(bytes, &mmi)
	if err != nil {
		return
	}
	mm := mmi.(Map)
	mut.Name = mm["name"].(string)
	switch mut.Name {
	case "all":
		argm := mm["args"].(Map)
		args := &AllArgs{}
		items := argm["items"].([]Any)
		args.Items = make([]*OneArgs, 0, len(items))
		for _, fmi := range items {
			fm := fmi.(Map)
			carg := &OneArgs{}
			carg.Id = parseUint(fm["id"])
			carg.Name = fm["name"].(string)
			carg.Json = fm["json"].(string)
		}
		mut.Args = args
	case "one":
		argm := mm["args"].(Map)
		args := &OneArgs{}
		args.Id = parseUint(argm["id"])
		args.Name = argm["name"].(string)
		args.Json = argm["json"].(string)
		mut.Args = args
	case "create":
		argm := mm["args"].(Map)
		args := &CreateArgs{}
		args.Id = maybeUint(argm["id"])
		args.Name = argm["name"].(string)
		args.Json = argm["json"].(string)
		mut.Args = args
	case "delete":
		argm := mm["args"].(Map)
		args := &DeleteArgs{}
		args.Id = parseUint(argm["id"])
		mut.Args = args
	case "update":
		argm := mm["args"].(Map)
		args := &UpdateArgs{}
		args.Id = parseUint(argm["id"])
		args.Name = argm["name"].(string)
		args.Json = argm["json"].(string)
		mut.Args = args
	default:
		err = fmt.Errorf("unkown mutation %s", mut.Name)
	}
	return
}

func maybeUint(id Any) uint {
	switch v := id.(type) {
	case float64:
		return uint(v)
	default:
		return 0
	}
}

func parseUint(id Any) uint {
	switch v := id.(type) {
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
	case "all":
		args := argi.(*AllArgs)
		argm = make(Map)
		items := make([]Map, 0, len(args.Items))
		for _, file := range args.Items {
			fm := make(Map)
			fm["id"] = file.Id
			fm["name"] = file.Name
			fm["json"] = file.Json
			items = append(items, fm)
		}
		argm["items"] = items
	case "one":
		args := argi.(*OneArgs)
		argm = make(Map)
		argm["id"] = args.Id
		argm["name"] = args.Name
		argm["json"] = args.Json
	case "create":
		args := argi.(*CreateArgs)
		argm = make(Map)
		argm["id"] = args.Id
		argm["name"] = args.Name
		argm["json"] = args.Json
	case "delete":
		args := argi.(*DeleteArgs)
		argm = make(Map)
		argm["id"] = args.Id
	case "update":
		args := argi.(*UpdateArgs)
		argm = make(Map)
		argm["id"] = args.Id
		argm["name"] = args.Name
		argm["json"] = args.Json
	default:
		err = fmt.Errorf("unkown mutation %s", name)
	}
	return
}
