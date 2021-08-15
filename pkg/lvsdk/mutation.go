package lvsdk

import (
	"fmt"
	"reflect"
	"strings"
)

type Mutation struct {
	Sid  string
	Name string
	Args Any
}

func M(name string, sid string, args Any) *Mutation {
	return &Mutation{Name: name, Sid: sid, Args: args}
}

func Mn(name string) *Mutation {
	return &Mutation{Name: name}
}

func Mns(name string, sid string) *Mutation {
	return &Mutation{Name: name, Sid: sid}
}

func Mna(name string, args Any) *Mutation {
	return &Mutation{Name: name, Args: args}
}

func (m *Mutation) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "{%s,%s,%v}", m.Name, m.Sid, reflect.ValueOf(m.Args))
	return buf.String()
}
