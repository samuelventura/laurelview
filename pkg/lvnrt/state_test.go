package lvnrt

import (
	"testing"
)

func TestRtStateDispose(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, disp Dispatch) {
		disp(Mns(":dispose", "tid"))
		to.MatchWait(t, 200, "trace", "state", "{:dispose,tid,<nil>,<nil>}")
		to.MatchWait(t, 200, "trace", "hub", "{:dispose,tid,<nil>,<nil>}")
		disp(Mns(":dispose", "tid"))
		to.MatchWait(t, 200, "debug", "state", "{:dispose,tid,<nil>,<nil>}")
	})
}

func TestRtStateBasic(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, disp Dispatch) {
		disp(Mnsa(":add", "tid", to.Dispatch("entry")))
		to.MatchWait(t, 200, "trace", "hub", "{:add,tid")
		disp(Mnsa("setup", "tid", []ItemArgs{{"host", 0, 1}}))
		to.MatchWait(t, 200, "trace", "bus", "{setup,tid,string,host:0}")
		to.MatchWait(t, 200, "trace", "bus", "{slave,tid,lvnrt.SlaveArgs,{1 1}}")
		to.MatchWait(t, 200, "trace", "hub", "{setup,tid")
		disp(Mnsa("query", "tid", QueryArgs{
			Index:   0,
			Request: "read-value",
		}))
		to.MatchWait(t, 200, "trace", "state", "{query,tid,lvnrt.QueryArgs,{0 read-value")
		to.MatchWait(t, 200, "trace", "bus", "{query,tid,lvnrt.QueryArgs,{1 read-value")
		disp(Mns(":remove", "tid"))
		to.MatchWait(t, 200, "trace", "bus", "{slave,tid,lvnrt.SlaveArgs,{1 0}}")
		to.MatchWait(t, 200, "trace", "bus", "{:dispose,tid,<nil>,<nil>}")
		to.MatchWait(t, 200, "trace", "hub", "{:remove,tid,<nil>,<nil>}")
		disp(Mns(":dispose", "tid"))
		to.MatchWait(t, 200, "trace", "state", "{:dispose,tid,<nil>,<nil>}")
		to.MatchWait(t, 200, "trace", "hub", "{:dispose,tid,<nil>,<nil>}")
		disp(Mns(":dispose", "tid"))
		to.MatchWait(t, 200, "debug", "state", "{:dispose,tid,<nil>,<nil>}")
	})
}

func testSetupState(callback func(to TestOutput, ctx Context, disp Dispatch)) {
	to := NewTestOutput()
	defer to.Close()
	log := to.Logger()
	ctx := NewContext(to.Log)
	defer WaitClose(ctx.Close)
	ctx.SetFactory("bus", func(Context) Dispatch { return to.Dispatch("bus") })
	ctx.SetDispatch("hub", to.Dispatch("hub"))
	disp := AsyncDispatch(log, NewState(ctx))
	callback(to, ctx, disp)
}
