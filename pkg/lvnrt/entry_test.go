package lvnrt

import (
	"fmt"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
)

func TestRtEntryBasic(t *testing.T) {
	testSetupEntry(t, func(to TestOutput, ctx Context, log Logger, conn *websocket.Conn, dp int) {
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
	testSetupEntry(t, func(to TestOutput, ctx Context, log Logger, conn *websocket.Conn, dp int) {
		testEntryPostSetup(conn, dp)
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||count:1||total:1")
		conn.Close()
		to.MatchWait(t, 200, "trace", "state", "{:remove,entry-1")
		to.MatchWait(t, 200, "trace", "hub", "{:remove,entry-1")
		to.MatchWait(t, 200, "trace", "entry-1", "out", "{:remove,entry-1")
	})
}

func TestRtEntryDisposeReceived(t *testing.T) {
	testSetupEntry(t, func(to TestOutput, ctx Context, log Logger, conn *websocket.Conn, dp int) {
		testEntryPostSetup(conn, dp)
		stateDispatch := ctx.GetDispatch("state")
		to.MatchWait(t, 200, "trace", "rsep", "query", "entry-1", "map||index:0||count:1||total:1")
		stateDispatch(Mns(":dispose", "tid"))
		to.MatchWait(t, 200, "trace", "entry-1", "out", "{:dispose,tid")
	})
}

func testSetupEntry(t *testing.T, callback func(to TestOutput, ctx Context, log Logger, conn *websocket.Conn, dpmPort int)) {
	to := NewTestOutput()
	defer to.Close() //wait flush
	log := to.Logger()
	dpm := NewDpm(log, ":0", 0)
	defer WaitClose(dpm.Close)
	log.Info("dpm", "port", dpm.Port())
	dpm.Echo()
	ctx := NewContext(to.Log)
	defer WaitClose(ctx.Close)
	ctx.SetValue("bus.dialtoms", 400)
	ctx.SetValue("bus.writetoms", 400)
	ctx.SetValue("bus.readtoms", 400)
	ctx.SetValue("bus.discardms", 50)
	ctx.SetValue("bus.sleepms", 10)
	ctx.SetValue("bus.retryms", 2000)
	ctx.SetValue("bus.resetms", 0)
	ctx.SetValue("entry.endpoint", ":0")
	ctx.SetValue("entry.buflen", 0)
	ctx.SetValue("entry.wtoms", 0)
	ctx.SetValue("entry.rtoms", 0)
	ctx.SetValue("entry.static", NopHandler)
	ctx.SetFactory("bus", func(ctx Context) Dispatch { return NewBus(ctx) })
	ctx.SetDispatch("hub", AsyncDispatch(log, NewHub(ctx)))
	ctx.SetDispatch("state", AsyncDispatch(log, NewState(ctx)))
	ctx.SetDispatch("check", NewCheck(ctx))
	checkDispatch := ctx.GetDispatch("check")
	defer checkDispatch(Mn(":dispose"))
	ctx.SetDispatch("/ws/test", checkDispatch)
	entry := NewEntry(ctx)
	defer WaitClose(entry.Close)
	log.Info("port", entry.Port())
	conn := testEntryConnect(entry.Port(), "/ws/test")
	defer conn.Close()
	log.Trace("client", conn.LocalAddr())
	go testEntryReadLoop(log.Trace, conn)
	callback(to, ctx, log, conn, dpm.Port())
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
