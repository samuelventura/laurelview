package lvnrt

import (
	"testing"
)

func TestRtHubBasic(t *testing.T) {
	to := NewTestOutput()
	defer to.Close()
	log := to.Logger()
	rt := NewRuntime(to.Log)
	defer rt.Close()
	disp := AsyncDispatch(log.Warn, NewHub(rt))
	disp(&Mutation{Name: "add", Sid: "tid", Args: &AddArgs{
		Callback: to.Dispatch("entry"),
	}})
	disp(&Mutation{Name: "setup", Sid: "tid", Args: &SetupArgs{
		Items: []*ItemArgs{{"host", 0, 1}},
	}})
	to.MatchWait(t, 200, "trace", "entry", "{query,tid,&{0   1 1}}")
	disp(&Mutation{Name: "status", Sid: "tid", Args: &StatusArgs{
		Slave:    "host:0:1",
		Request:  "read-value",
		Response: "value",
	}})
	to.MatchWait(t, 200, "trace", "entry", "{query,tid,&{0 read-value value 2 2}}")
	disp(&Mutation{Name: "remove", Sid: "tid"})
	to.MatchWait(t, 200, "trace", "entry", "{remove,tid")
	disp(&Mutation{Name: "dispose", Sid: "tid"})
}
