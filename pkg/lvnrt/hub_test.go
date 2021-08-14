package lvnrt

import (
	"testing"
)

func TestRtHubBasic(t *testing.T) {
	to := newTestOutput()
	defer to.close()
	log := to.logger()
	rt := NewRuntime(to.push)
	defer rt.ManagedWait()
	disp := asyncDispatch(log.Warn, NewHub(rt))
	disp(&Mutation{Name: "add", Sid: "tid", Args: &AddArgs{
		Callback: to.dispatch("entry"),
	}})
	disp(&Mutation{Name: "setup", Sid: "tid", Args: &SetupArgs{
		Items: []*ItemArgs{{"host", 0, 1}},
	}})
	to.matchWait(t, 200, "trace", "entry", "{query,tid,&{0  }}")
	disp(&Mutation{Name: "status", Sid: "tid", Args: &StatusArgs{
		Slave:    "host:0:1",
		Request:  "read-value",
		Response: "value",
	}})
	to.matchWait(t, 200, "trace", "entry", "{query,tid,&{0 read-value value}}")
	disp(&Mutation{Name: "remove", Sid: "tid"})
	to.matchWait(t, 200, "trace", "entry", "{remove,tid")
	disp(&Mutation{Name: "dispose", Sid: "tid"})
}
