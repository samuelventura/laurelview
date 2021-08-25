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
	bid := NewId("bus")
	rt.SetFactory("bus", func(Runtime) Dispatch {
		return to.Dispatch(bid.Next())
	})
	rt.SetDispatch("hub", to.Dispatch("hub"))
	disp := AsyncDispatch(log.Warn, NewState(rt))
	disp(&Mutation{Name: ":add", Sid: "tid", Args: to.Dispatch("entry")})
	to.MatchWait(t, 200, "trace", "hub", "{:add,tid")
	disp(&Mutation{Name: "setup", Sid: "tid", Args: []*ItemArgs{{"host", 0, 1}}})
	to.MatchWait(t, 200, "trace", "bus-1", "{setup,tid,.*BusArgs,&{host 0}")
	to.MatchWait(t, 200, "trace", "bus-1", "{slave,tid,.*SlaveArgs,&{1 1}")
	to.MatchWait(t, 200, "trace", "hub", "{setup,tid")
	disp(&Mutation{Name: "query", Sid: "tid", Args: &QueryArgs{
		Index:   0,
		Request: "read-value",
	}})
	to.MatchWait(t, 200, "trace", "bus-1", "{query,tid,.*QueryArgs,&{1 read-value")
	disp(&Mutation{Name: ":remove", Sid: "tid"})
	to.MatchWait(t, 200, "trace", "bus-1", "{slave,tid,.*SlaveArgs,&{1 0}}")
	to.MatchWait(t, 200, "trace", "bus-1", "{:dispose,tid,")
	to.MatchWait(t, 200, "trace", "hub", "{:remove,tid")
	disp(&Mutation{Name: ":dispose", Sid: "tid"})
	to.MatchWait(t, 200, "trace", "hub", "{:dispose,tid")
}
