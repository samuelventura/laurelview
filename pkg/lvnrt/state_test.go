package lvnrt

import (
	"testing"
)

func TestRtStateBasic(t *testing.T) {
	to := NewTestOutput()
	defer to.Close()
	log := to.Logger()
	rt := NewRuntime(to.Log)
	defer WaitClose(rt.Close)
	rt.SetFactory("bus", func(Runtime) Dispatch {
		return to.Dispatch("bus")
	})
	rt.SetDispatch("hub", to.Dispatch("hub"))
	disp := AsyncDispatch(log, NewState(rt))
	disp(Mnsa(":add", "tid", to.Dispatch("entry")))
	to.MatchWait(t, 200, "trace", "hub", "{:add,tid")
	disp(Mnsa("setup", "tid", []*ItemArgs{{"host", 0, 1}}))
	to.MatchWait(t, 200, "trace", "bus", "{setup,tid,.*BusArgs,&{host 0}")
	to.MatchWait(t, 200, "trace", "bus", "{slave,tid,.*SlaveArgs,&{1 1}")
	to.MatchWait(t, 200, "trace", "hub", "{setup,tid")
	disp(Mnsa("query", "tid", &QueryArgs{
		Index:   0,
		Request: "read-value",
	}))
	to.MatchWait(t, 200, "trace", "bus", "{query,tid,.*QueryArgs,&{1 read-value")
	disp(Mns(":remove", "tid"))
	to.MatchWait(t, 200, "trace", "bus", "{slave,tid,.*SlaveArgs,&{1 0}}")
	to.MatchWait(t, 200, "trace", "bus", "{:dispose,tid,")
	to.MatchWait(t, 200, "trace", "hub", "{:remove,tid")
	disp(Mns(":dispose", "tid"))
	to.MatchWait(t, 200, "trace", "state", "{:dispose,tid")
	to.MatchWait(t, 200, "trace", "hub", "{:dispose,tid")
	disp(Mns(":dispose", "tid"))
	to.MatchWait(t, 200, "debug", "state", "{:dispose,tid")
}
