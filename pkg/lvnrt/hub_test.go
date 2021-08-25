package lvnrt

import (
	"testing"
)

func TestRtHubBasic(t *testing.T) {
	to := NewTestOutput()
	defer to.Close()
	log := to.Logger()
	rt := NewRuntime(to.Log)
	defer WaitClose(rt.Close)
	disp := AsyncDispatch(log.Warn, NewHub(rt))
	disp(&Mutation{Name: ":add", Sid: "tid", Args: to.Dispatch("entry")})
	disp(&Mutation{Name: "setup", Sid: "tid", Args: []*ItemArgs{{"host", 0, 1}}})
	to.MatchWait(t, 200, "trace", "hub", "{setup,tid,.*SetupArgs,&{")
	to.MatchWait(t, 200, "trace", "entry", "{query,tid,.*QueryArgs,&{0   1 1")
	disp(&Mutation{Name: "status-slave", Sid: "tid", Args: &StatusArgs{
		Address:  "host:0:1",
		Request:  "read-value",
		Response: "value",
	}})
	to.MatchWait(t, 200, "trace", "entry", "{query,tid,.*QueryArgs,&{0 read-value value 2 2")
	disp(&Mutation{Name: ":remove", Sid: "tid"})
	to.MatchWait(t, 200, "trace", "entry", "{:remove,tid")
	disp(&Mutation{Name: ":dispose", Sid: "tid"})
}
