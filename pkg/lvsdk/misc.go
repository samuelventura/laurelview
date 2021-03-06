package lvsdk

import "reflect"

func ToMap(any Any) Map {
	m := make(Map)
	e := reflect.ValueOf(any).Elem()
	t := e.Type()
	m["$type"] = t.Name()
	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		ft := t.Field(i)
		m[ft.Name] = f.Interface()
	}
	return m
}
