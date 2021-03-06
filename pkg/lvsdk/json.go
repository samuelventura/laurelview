package lvsdk

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func EncodeMutation(mut Mutation) (bytes []byte, err error) {
	mm := make(Map)
	mm["sid"] = mut.Sid
	mm["name"] = mut.Name
	switch val := mut.Args.(type) {
	case nil:
		bytes, err = json.Marshal(mm)
		return
	default:
		mm["args"] = val
		bytes, err = json.Marshal(mm)
		return
	}
}

func DecodeMutation(bytes []byte) (mut Mutation, err error) {
	var mmi Any
	err = json.Unmarshal(bytes, &mmi)
	if err != nil {
		return
	}
	mm := mmi.(Map)
	mut = Mutation{}
	name, err := ParseString(mm, "name")
	if err != nil {
		return
	}
	mut.Name = name
	sid, err := MaybeString(mm, "sid", "")
	if err != nil {
		return
	}
	mut.Sid = sid
	args, ok := mm["args"]
	if ok {
		mut.Args = args
	}
	return
}

func ParseString(mm Map, key string) (res string, err error) {
	val, err := getValue(mm, key)
	if err != nil {
		return
	}
	switch cur := val.(type) {
	case string:
		res = cur
		return
	default:
		typ := reflect.TypeOf(val)
		err = fmt.Errorf("invalid type `%v:%v`", key, typ)
		return
	}
}

func ParseUint(mm Map, key string) (res uint, err error) {
	val, err := getValue(mm, key)
	if err != nil {
		return
	}
	switch cur := val.(type) {
	case float64:
		res = uint(cur)
		return
	case string:
		u64, err2 := strconv.ParseUint(cur, 10, 64)
		if err2 != nil {
			err = fmt.Errorf("parse error %v:%w", key, err2)
			return
		}
		res = uint(u64)
		return
	default:
		typ := reflect.TypeOf(val)
		err = fmt.Errorf("invalid type `%v:%v`", key, typ)
		return
	}
}

func MaybeString(mm Map, key string, def string) (res string, err error) {
	val, err := getValue(mm, key)
	if err != nil {
		err = nil
		res = def
		return
	}
	switch cur := val.(type) {
	case string:
		res = cur
		return
	default:
		typ := reflect.TypeOf(val)
		err = fmt.Errorf("invalid type `%v:%v`", key, typ)
		return
	}
}

func MaybeUint(mm Map, key string, def uint) (res uint, err error) {
	val, err := getValue(mm, key)
	if err != nil {
		err = nil
		res = def
		return
	}
	switch cur := val.(type) {
	case float64:
		res = uint(cur)
		return
	case string:
		u64, err2 := strconv.ParseUint(cur, 10, 64)
		if err2 != nil {
			err = fmt.Errorf("parse error %v:%w", key, err2)
			return
		}
		res = uint(u64)
		return
	default:
		typ := reflect.TypeOf(val)
		err = fmt.Errorf("invalid type `%v:%v`", key, typ)
		return
	}
}

func CastUint(val Any, key string) (res uint, err error) {
	switch cur := val.(type) {
	case float64:
		res = uint(cur)
		return
	case string:
		u64, err2 := strconv.ParseUint(cur, 10, 64)
		if err2 != nil {
			err = fmt.Errorf("parse error %v:%w", key, err2)
			return
		}
		res = uint(u64)
		return
	default:
		typ := reflect.TypeOf(val)
		err = fmt.Errorf("invalid type `%v:%v`", key, typ)
		return
	}
}

func getValue(mm Map, key string) (res Any, err error) {
	val, ok := mm[key]
	if !ok {
		err = fmt.Errorf("key not found `%v`", key)
		return
	}
	res = val
	return
}
