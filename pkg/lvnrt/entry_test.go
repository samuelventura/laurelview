package lvnrt

import (
	"fmt"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
)

func TestRtEntryBasic(t *testing.T) {
	testSetupEntry(t, func(to TestOutput, rt Runtime, log Logger, conn *websocket.Conn, dp int) {
		testEntryPostSetup(conn, dp)
		//first empty
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||count:1||total:1")
		//read value
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||count:2||total:2||request:read-value||response:.1B1")
		testEntryPostQuery(conn, "reset-valley")
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||request:reset-valley||response:ok")
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||request:read-valley||response:.1B3")
		testEntryPostQuery(conn, "reset-peak")
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||request:reset-peak||response:ok")
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||request:read-peak||response:.1B2")
		testEntryPostQuery(conn, "apply-tara")
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||request:apply-tara||response:ok")
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||request:read-value||response:.1B1")
		testEntryPostQuery(conn, "reset-tara")
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||request:reset-tara||response:ok")
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||request:read-value||response:.1B1")
		conn.Close()
		to.MatchWait(t, 200, "trace", "state", "{:remove,entry-1")
		to.MatchWait(t, 200, "trace", "hub", "{:remove,entry-1")
		to.MatchWait(t, 200, "trace", "entry-1", "out", "{:remove,entry-1")
	})
}

func TestRtEntryRemoveReceived(t *testing.T) {
	testSetupEntry(t, func(to TestOutput, rt Runtime, log Logger, conn *websocket.Conn, dp int) {
		testEntryPostSetup(conn, dp)
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||count:1||total:1")
		conn.Close()
		to.MatchWait(t, 200, "trace", "state", "{:remove,entry-1")
		to.MatchWait(t, 200, "trace", "hub", "{:remove,entry-1")
		to.MatchWait(t, 200, "trace", "entry-1", "out", "{:remove,entry-1")
	})
}

func TestRtEntryDisposeReceived(t *testing.T) {
	testSetupEntry(t, func(to TestOutput, rt Runtime, log Logger, conn *websocket.Conn, dp int) {
		testEntryPostSetup(conn, dp)
		stateDispatch := rt.GetDispatch("state")
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||count:1||total:1")
		stateDispatch(Mns(":dispose", "tid"))
		to.MatchWait(t, 200, "trace", "entry-1", "out", "{:dispose,tid")
	})
}

func testSetupEntry(t *testing.T, callback func(to TestOutput, rt Runtime, log Logger, conn *websocket.Conn, dpmPort int)) {
	to := NewTestOutput()
	defer to.Close() //wait flush
	log := to.Logger()
	dpm := NewDpm(log, ":0", 0)
	defer WaitClose(dpm.Close)
	log.Info("dpm", "port", dpm.Port())
	dpm.Echo()
	rt := NewRuntime(to.Log)
	defer WaitClose(rt.Close)
	rt.SetValue("bus.dialtoms", 400)
	rt.SetValue("bus.writetoms", 400)
	rt.SetValue("bus.readtoms", 400)
	rt.SetValue("bus.discardms", 50)
	rt.SetValue("bus.sleepms", 10)
	rt.SetValue("bus.retryms", 2000)
	rt.SetValue("bus.resetms", 0)
	rt.SetValue("entry.endpoint", ":0")
	rt.SetValue("entry.buflen", 0)
	rt.SetValue("entry.static", NopHandler)
	rt.SetDispatch("hub", AsyncDispatch(log, NewHub(rt)))
	rt.SetDispatch("state", AsyncDispatch(log, NewState(rt)))
	rt.SetDispatch("checkin", NewCheckin(rt))
	checkinDispatch := rt.GetDispatch("checkin")
	defer checkinDispatch(Mn(":dispose"))
	rt.SetFactory("bus", func(rt Runtime) Dispatch { return NewBus(rt) })
	rt.SetDispatch("/ws/test", checkinDispatch)
	entry := NewEntry(rt)
	defer WaitClose(entry.Close)
	log.Info("port", entry.Port())
	conn := testEntryConnect(entry.Port(), "/ws/test")
	defer conn.Close()
	log.Trace("client", conn.LocalAddr())
	go testEntryReadLoop(log.Trace, conn)
	callback(to, rt, log, conn, dpm.Port())
}

func testEntryConnect(port int, path string) *websocket.Conn {
	url := fmt.Sprintf("ws://localhost:%v%v", port, path)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	PanicIfError(err)
	return conn
}

func testEntryReadLoop(output Output, conn *websocket.Conn) {
	for {
		mut := testEntryReadMutation(conn)
		if mut.Name == "" {
			return
		}
		output("rmut", mut)
		output("rsep", mut.Name, mut.Sid, fmt.Sprint(mut.Args))
	}
}

func testEntryReadMutation(conn *websocket.Conn) Mutation {
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 400))
	mt, bytes, err := conn.ReadMessage()
	if err != nil {
		return Mutation{}
	}
	if mt != websocket.TextMessage {
		PanicF("Invalid msg type %v", mt)
	}
	mut, err := DecodeMutation(bytes)
	if err != nil {
		return Mutation{}
	}
	return mut
}

func testEntryPostSetup(conn *websocket.Conn, port int) {
	args := []ItemArgs{{Host: "127.0.0.1", Port: uint(port), Slave: 1}}
	mut := Mna("setup", args)
	WriteMutation(conn, mut)
}

func testEntryPostQuery(conn *websocket.Conn, request string) {
	args := QueryArgs{}
	args.Index = 0
	args.Request = request
	mut := Mna("query", args)
	WriteMutation(conn, mut)
}
