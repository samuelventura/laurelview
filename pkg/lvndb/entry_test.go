package lvndb

import (
	"fmt"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
)

func TestDbEntryCrud(t *testing.T) {
	testSetupEntry(t, func(to TestOutput, rt Runtime, log Logger, dao Dao, conn *websocket.Conn) {
		to.MatchWait(t, 200, "trace", "rmut", "{all,entry-1.*,..interface {},..}")
		testEntryWriteMutation(conn, Mna("create", OneArgs{Name: "name1", Json: "json1"}))
		to.MatchWait(t, 200, "trace", "rsep", "create", "entry-1.*", "id:1||json:json1||name:name1")
		testEntryWriteMutation(conn, Mna("update", OneArgs{Id: 1, Name: "name2", Json: "json2"}))
		to.MatchWait(t, 200, "trace", "rsep", "update", "entry-1.*", "id:1||json:json2||name:name2")
		testEntryWriteMutation(conn, Mna("delete", uint(1)))
		to.MatchWait(t, 200, "trace", "rsep", "delete", "entry-1.*", "1")
	})
}

func TestRtEntryRemoveReceived(t *testing.T) {
	testSetupEntry(t, func(to TestOutput, rt Runtime, log Logger, dao Dao, conn *websocket.Conn) {
		to.MatchWait(t, 200, "trace", "rmut", "{all,entry-1.*,..interface {},..}")
		conn.Close()
		to.MatchWait(t, 200, "trace", "state", "{:remove,entry-1.*,<nil>,<nil>}")
		to.MatchWait(t, 200, "trace", "hub", "{:remove,entry-1.*,<nil>,<nil>}")
		to.MatchWait(t, 200, "trace", "entry-1", "out", "{:remove,entry-1.*,<nil>,<nil>}")
	})
}

func TestRtEntryDisposeReceived(t *testing.T) {
	testSetupEntry(t, func(to TestOutput, rt Runtime, log Logger, dao Dao, conn *websocket.Conn) {
		to.MatchWait(t, 200, "trace", "rmut", "{all,entry-1.*,..interface {},..}")
		stateDispatch := rt.GetDispatch("state")
		stateDispatch(Mns(":dispose", "tid"))
		to.MatchWait(t, 200, "trace", "state", "{:dispose,tid,<nil>,<nil>}")
		to.MatchWait(t, 200, "trace", "hub", "{:dispose,tid,<nil>,<nil>}")
		to.MatchWait(t, 200, "trace", "entry-1", "out", "{:dispose,tid,<nil>,<nil>}")
	})
}

func testSetupEntry(t *testing.T, callback func(to TestOutput, rt Runtime, log Logger, dao Dao, conn *websocket.Conn)) {
	var dao = NewDao(":memory:")
	defer dao.Close()
	to := NewTestOutput()
	defer to.Close() //wait flush
	log := to.Logger()
	rt := NewRuntime(to.Log)
	defer WaitClose(rt.Close)
	rt.SetValue("dao", dao)
	rt.SetValue("entry.endpoint", ":0")
	rt.SetValue("entry.buflen", 0)
	rt.SetValue("entry.static", NopHandler)
	rt.SetDispatch("hub", AsyncDispatch(log, NewHub(rt)))
	rt.SetDispatch("state", AsyncDispatch(log, NewState(rt)))
	rt.SetDispatch("checkin", NewCheckin(rt))
	checkinDispatch := rt.GetDispatch("checkin")
	defer checkinDispatch(Mn(":dispose"))
	rt.SetDispatch("/ws/test", checkinDispatch)
	entry := NewEntry(rt)
	defer WaitClose(entry.Close)
	log.Info("port", entry.Port())
	conn := testEntryConnect(entry.Port(), "/ws/test")
	defer conn.Close()
	log.Trace("client", conn.LocalAddr())
	go testEntryReadLoop(log.Trace, conn)
	callback(to, rt, log, dao, conn)
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

func testEntryWriteMutation(conn *websocket.Conn, mut Mutation) {
	bytes, err := EncodeMutation(mut)
	PanicIfError(err)
	err = conn.WriteMessage(websocket.TextMessage, bytes)
	PanicIfError(err)
}