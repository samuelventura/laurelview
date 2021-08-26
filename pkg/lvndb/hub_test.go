package lvndb

import (
	"testing"
)

func TestDbHubDispose(t *testing.T) {
	testSetupHub(func(to TestOutput, rt Runtime, log Logger) {
		disp := NewHub(rt)
		disp(Mnsa(":add", "tid1", to.Dispatch("callback1")))
		disp(Mnsa(":add", "tid2", to.Dispatch("callback2")))
		to.MatchWait(t, 200, "trace", "hub", "{:add,tid1,")
		to.MatchWait(t, 200, "trace", "hub", "{:add,tid2,")
		disp(Mn(":dispose"))
		to.MatchWait(t, 200, "trace", "hub", "{:dispose,,<nil>,<nil>}")
		to.MatchWait(t, 200, "trace", "callback1", "{:dispose,,<nil>,<nil>}")
		to.MatchWait(t, 200, "trace", "callback2", "{:dispose,,<nil>,<nil>}")
		disp(Mn(":dispose"))
		to.MatchWait(t, 200, "debug", "hub", "{:dispose,,<nil>,<nil>}")
	})
}

func TestDbHubAddRemove(t *testing.T) {
	testSetupHub(func(to TestOutput, rt Runtime, log Logger) {
		disp := NewHub(rt)
		disp(Mnsa(":add", "tid", to.Dispatch("callback")))
		to.MatchWait(t, 200, "trace", "hub", "{:add,tid,")
		disp(Mnsa(":add", "tid", to.Dispatch("callback")))
		to.MatchWait(t, 200, "debug", "hub", "{:add,tid,", "duplicated sid", "tid")
		disp(Mns(":remove", "tid"))
		to.MatchWait(t, 200, "trace", "hub", "{:remove,tid,<nil>,<nil>}")
		to.MatchWait(t, 200, "trace", "callback", "{:remove,tid,<nil>,<nil>}")
		disp(Mns(":remove", "tid"))
		to.MatchWait(t, 200, "debug", "hub", "{:remove,tid,<nil>,<nil>}", "non-existent sid", "tid")
	})
}

func TestDbHubAll(t *testing.T) {
	testSetupHub(func(to TestOutput, rt Runtime, log Logger) {
		disp := NewHub(rt)
		disp(Mnsa(":add", "tid1", to.Dispatch("callback1")))
		disp(Mnsa(":add", "tid2", func(mut Mutation) { PanicLN("tid2", mut) }))
		disp(Mnsa("all", "tid1", []OneArgs{{Id: 1, Name: "name1", Json: "json1"}}))
		to.MatchWait(t, 200, "trace", "hub", "{all,tid1,..lvndb.OneArgs,.{1 name1 json1}.}")
		to.MatchWait(t, 200, "trace", "callback1", "{all,tid1,..lvndb.OneArgs,.{1 name1 json1}.}")
		disp(Mnsa("all", "tid3", []OneArgs{{Id: 1, Name: "name1", Json: "json1"}}))
		to.MatchWait(t, 200, "debug", "hub", "{all,tid3,..lvndb.OneArgs,.{1 name1 json1}.}", "non-existent sid", "tid3")
	})
}

func TestDbHubCreate(t *testing.T) {
	testSetupHub(func(to TestOutput, rt Runtime, log Logger) {
		disp := NewHub(rt)
		disp(Mnsa(":add", "tid1", to.Dispatch("callback1")))
		disp(Mnsa(":add", "tid2", to.Dispatch("callback2")))
		disp(Mnsa("create", "tid1", OneArgs{Id: 1, Name: "name1", Json: "json1"}))
		to.MatchWait(t, 200, "trace", "hub", "{create,tid1,lvndb.OneArgs,{1 name1 json1}}")
		to.MatchWait(t, 200, "trace", "callback1", "{create,tid1,lvndb.OneArgs,{1 name1 json1}}")
		to.MatchWait(t, 200, "trace", "callback2", "{create,tid1,lvndb.OneArgs,{1 name1 json1}}")
	})
}

func TestDbHubUpdate(t *testing.T) {
	testSetupHub(func(to TestOutput, rt Runtime, log Logger) {
		disp := NewHub(rt)
		disp(Mnsa(":add", "tid1", to.Dispatch("callback1")))
		disp(Mnsa(":add", "tid2", to.Dispatch("callback2")))
		disp(Mnsa("update", "tid1", OneArgs{Id: 1, Name: "name1", Json: "json1"}))
		to.MatchWait(t, 200, "trace", "hub", "{update,tid1,lvndb.OneArgs,{1 name1 json1}}")
		to.MatchWait(t, 200, "trace", "callback1", "{update,tid1,lvndb.OneArgs,{1 name1 json1}}")
		to.MatchWait(t, 200, "trace", "callback2", "{update,tid1,lvndb.OneArgs,{1 name1 json1}}")
	})
}

func TestDbHubDelete(t *testing.T) {
	testSetupHub(func(to TestOutput, rt Runtime, log Logger) {
		disp := NewHub(rt)
		disp(Mnsa(":add", "tid1", to.Dispatch("callback1")))
		disp(Mnsa(":add", "tid2", to.Dispatch("callback2")))
		disp(Mnsa("delete", "tid1", uint(1)))
		to.MatchWait(t, 200, "trace", "hub", "{delete,tid1,uint,1}")
		to.MatchWait(t, 200, "trace", "callback1", "{delete,tid1,uint,1}")
		to.MatchWait(t, 200, "trace", "callback2", "{delete,tid1,uint,1}")
	})
}

func testSetupHub(callback func(to TestOutput, rt Runtime, log Logger)) {
	to := NewTestOutput()
	defer to.Close()
	log := to.Logger()
	rt := NewRuntime(to.Log)
	defer WaitClose(rt.Close)
	callback(to, rt, log)
}
