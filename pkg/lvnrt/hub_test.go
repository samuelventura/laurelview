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
	disp := AsyncDispatch(log, NewHub(rt))
	disp(Mnsa(":add", "tid", to.Dispatch("entry")))
	disp(Mnsa("setup", "tid", []*ItemArgs{{"host", 0, 1}}))
	to.MatchWait(t, 200, "trace", "hub", "{setup,tid,.*ItemArgs,")
	to.MatchWait(t, 200, "trace", "entry", "{query,tid,.*QueryArgs,&{0   1 1")
	disp(Mnsa("status-slave", "tid", &StatusArgs{
		Address:  "host:0:1",
		Request:  "read-value",
		Response: "value",
	}))
	to.MatchWait(t, 200, "trace", "entry", "{query,tid,.*QueryArgs,&{0 read-value value 2 2")
	disp(Mns(":remove", "tid"))
	to.MatchWait(t, 200, "trace", "entry", "{:remove,tid")
	disp(Mns(":dispose", "tid"))
	to.MatchWait(t, 200, "trace", "hub", "{:dispose,tid")
	disp(Mns(":dispose", "tid"))
	to.MatchWait(t, 200, "debug", "hub", "{:dispose,tid")
}
