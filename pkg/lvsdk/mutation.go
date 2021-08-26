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

func M(name string, sid string, args Any) Mutation {
	return Mutation{Name: name, Sid: sid, Args: args}
}

func Mn(name string) Mutation {
	return Mutation{Name: name}
}

func Mns(name string, sid string) Mutation {
	return Mutation{Name: name, Sid: sid}
}

func Mna(name string, args Any) Mutation {
	return Mutation{Name: name, Args: args}
}

func Mnsa(name string, sid string, args Any) Mutation {
	return Mutation{Name: name, Sid: sid, Args: args}
}

func (m Mutation) String() string {
	buf := new(strings.Builder)
	var typ string
	switch m.Args.(type) {
	case nil:
		typ = "<nil>"
	default:
		typ = reflect.TypeOf(m.Args).String()
	}
	fmt.Fprintf(buf, "{%s,%s,%s,%v}", m.Name, m.Sid, typ, m.Args)
	return buf.String()
}
