package lvnrt

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
)

//FIXME state remove get called many times on exit

func TestRtEntryBasic(t *testing.T) {
	testSetupEntry(t, func(to TestOutput, rt Runtime, log Logger, conn *websocket.Conn, dp int) {
		postSetup(conn, dp)
		//first empty
		to.MatchWait(t, 200, "trace", "read", "query", "client-1", "&{0   1 1")
		//read value
		to.MatchWait(t, 200, "trace", "read", "query", "client-1", "&{0 read-value .1B1")
		postQuery(conn, "reset-valley")
		to.MatchWait(t, 200, "trace", "read", "query", "client-1", "&{0 reset-valley ok")
		to.MatchWait(t, 200, "trace", "read", "query", "client-1", "&{0 read-valley .1B3")
		postQuery(conn, "reset-peak")
		to.MatchWait(t, 200, "trace", "read", "query", "client-1", "&{0 reset-peak ok")
		to.MatchWait(t, 200, "trace", "read", "query", "client-1", "&{0 read-peak .1B2")
		postQuery(conn, "apply-tara")
		to.MatchWait(t, 200, "trace", "read", "query", "client-1", "&{0 apply-tara ok")
		to.MatchWait(t, 200, "trace", "read", "query", "client-1", "&{0 read-value .1B1")
		postQuery(conn, "reset-tara")
		to.MatchWait(t, 200, "trace", "read", "query", "client-1", "&{0 reset-tara ok")
		to.MatchWait(t, 200, "trace", "read", "query", "client-1", "&{0 read-value .1B1")
		conn.Close()
		to.MatchWait(t, 200, "trace", "state", "{remove,client-1")
		to.MatchWait(t, 200, "trace", "hub", "{remove,client-1")
		to.MatchWait(t, 200, "trace", "client-1", "out", "{remove,client-1")
	})
}

func TestRtEntryRemoveReceived(t *testing.T) {
	testSetupEntry(t, func(to TestOutput, rt Runtime, log Logger, conn *websocket.Conn, dp int) {
		postSetup(conn, dp)
		conn.Close()
		to.MatchWait(t, 200, "trace", "state", "{remove,client-1")
		to.MatchWait(t, 200, "trace", "hub", "{remove,client-1")
		to.MatchWait(t, 200, "trace", "client-1", "out", "{remove,client-1")
	})
}

func TestRtEntryDisposeReceived(t *testing.T) {
	testSetupEntry(t, func(to TestOutput, rt Runtime, log Logger, conn *websocket.Conn, dp int) {
		postSetup(conn, dp)
		to.MatchWait(t, 200, "trace", "read", "query", "client-1", "&{0   1 1")
		rt.GetDispatch("state")(Mns("dispose", "tid"))
		to.MatchWait(t, 200, "trace", "client-1", "out", "{dispose,tid")
	})
}

func testSetupEntry(t *testing.T, callback func(to TestOutput, rt Runtime, log Logger, conn *websocket.Conn, dp int)) {
	to := NewTestOutput()
	log := to.Logger()
	dpm := NewDpm(log, ":0", 0)
	log.Info("dpm", "port", dpm.Port())
	defer dpm.Close()
	dpm.Echo()
	rt := NewRuntime(to.Log)
	defer rt.Close()
	rt.SetValue("bus.dialtoms", 400)
	rt.SetValue("bus.writetoms", 400)
	rt.SetValue("bus.readtoms", 400)
	rt.SetValue("bus.discardms", 50)
	rt.SetValue("bus.sleepms", 10)
	rt.SetValue("bus.retryms", 2000)
	rt.SetValue("bus.resetms", 0)
	rt.SetCleaner("bus", NewCleaner(log))
	defer log.Log("") //wait flush
	defer TraceRecover(log.Warn)
	rt.SetDispatch("hub", AsyncDispatch(log.Debug, NewHub(rt)))
	rt.SetDispatch("state", AsyncDispatch(log.Debug, NewState(rt)))
	defer rt.GetDispatch("state")(&Mutation{Name: "dispose"})
	rt.SetFactory("bus", func(rt Runtime) Dispatch { return NewBus(rt) })
	id := NewId("client")
	entry := NewEntry(rt, id, ":0")
	defer entry.Close()
	log.Info("port", entry.Port())
	conn := connect(entry.Port())
	defer conn.Close()
	log.Trace("client", conn.LocalAddr())
	go readLoop(log.Trace, conn)
	callback(to, rt, log, conn, int(dpm.Port()))
}

func connect(port int) *websocket.Conn {
	url := fmt.Sprintf("ws://localhost:%v/ws/index", port)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	PanicIfError(err)
	return conn
}

func readLoop(output Output, conn *websocket.Conn) {
	for {
		mut := readMutation(conn)
		if mut.Name == "" {
			return
		}
		output("read", mut.Name, mut.Sid, reflect.ValueOf(mut.Args))
	}
}

func readMutation(conn *websocket.Conn) *Mutation {
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 400))
	mt, bytes, err := conn.ReadMessage()
	if err != nil {
		return &Mutation{}
	}
	if mt != websocket.TextMessage {
		PanicF("Invalid msg type %v", mt)
	}
	mut, err := decodeMutationEx(bytes, true)
	if err != nil {
		return &Mutation{}
	}
	return mut
}

func postSetup(conn *websocket.Conn, port int) {
	args := &SetupArgs{}
	args.Items = []*ItemArgs{{"127.0.0.1", uint(port), 1}}
	mut := &Mutation{Name: "setup", Args: args}
	bytes := encodeMutation(mut)
	err := conn.WriteMessage(websocket.TextMessage, bytes)
	PanicIfError(err)
}

func postQuery(conn *websocket.Conn, request string) {
	args := &QueryArgs{}
	args.Index = 0
	args.Request = request
	mut := &Mutation{Name: "query", Args: args}
	bytes := encodeMutation(mut)
	err := conn.WriteMessage(websocket.TextMessage, bytes)
	PanicIfError(err)
}
