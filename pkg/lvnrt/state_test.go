package lvnrt

import (
	"testing"
)

func TestRtStateBasic(t *testing.T) {
	to := NewTestOutput()
	defer to.Close()
	log := to.Logger()
	rt := NewRuntime(to.Log)
	defer rt.Close()
	bid := NewId("bus")
	rt.Setf("bus", func(Runtime) Dispatch {
		return to.Dispatch(bid.Next())
	})
	rt.Setd("hub", to.Dispatch("hub"))
	disp := AsyncDispatch(log.Warn, NewState(rt))
	disp(&Mutation{Name: "add", Sid: "tid", Args: &AddArgs{
		Callback: to.Dispatch("entry"),
	}})
	to.MatchWait(t, 200, "trace", "hub", "{add,tid")
	disp(&Mutation{Name: "setup", Sid: "tid", Args: &SetupArgs{
		Items: []*ItemArgs{{"host", 0, 1}},
	}})
	to.MatchWait(t, 200, "trace", "bus-1", "{setup,tid,&{host 0}")
	to.MatchWait(t, 200, "trace", "bus-1", "{slave,tid,&{1 1}")
	to.MatchWait(t, 200, "trace", "hub", "{setup,tid")
	disp(&Mutation{Name: "query", Sid: "tid", Args: &QueryArgs{
		Index:   0,
		Request: "read-value",
	}})
	to.MatchWait(t, 200, "trace", "bus-1", "{query,tid,&{1 read-value }}")
	disp(&Mutation{Name: "remove", Sid: "tid"})
	to.MatchWait(t, 200, "trace", "bus-1", "{slave,tid,&{1 0}}")
	to.MatchWait(t, 200, "trace", "bus-1", "{dispose,tid,")
	to.MatchWait(t, 200, "trace", "hub", "{remove,tid")
	disp(&Mutation{Name: "dispose", Sid: "tid"})
	to.MatchWait(t, 200, "trace", "hub", "{dispose,tid")
}
