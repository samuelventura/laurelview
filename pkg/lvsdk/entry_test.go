package lvsdk

import (
	"fmt"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
)

func TestSdkEntryBasic(t *testing.T) {
	to := NewTestOutput()
	defer to.Close()
	rt := NewRuntime(to.Log)
	defer WaitClose(rt.Close)
	var callback Dispatch
	rt.SetValue("entry.endpoint", ":0")
	rt.SetValue("entry.buflen", 0)
	rt.SetValue("entry.static", NopHandler)
	rt.SetDispatch("/ws/test", func(mut Mutation) {
		to.Trace("disp", mut)
		switch mut.Name {
		case ":add":
			callback = mut.Args.(Dispatch)
		}
	})
	entry := NewEntry(rt)
	defer WaitClose(entry.Close)
	conn := testEntryConnect(entry.Port(), "/ws/test")
	defer conn.Close()
	to.MatchWait(t, 200, "trace", "entry-1-", "path", "/ws/test")
	to.MatchWait(t, 200, "trace", "disp", "{:add,entry-1-")
	callback(Mns("init", "tid"))
	to.MatchWait(t, 200, "trace", "entry-1-", "out", "{init,tid,<nil>,<nil>}")
	go testEntryReadLoop(conn, to)
	to.MatchWait(t, 200, "trace", "test", "read", "{init,tid,<nil>,<nil>}")
	testEntryPostMutation(conn, Mn("query"))
	to.MatchWait(t, 200, "trace", "entry-1-", "in", "{query,,<nil>,<nil>}")
	to.MatchWait(t, 200, "trace", "disp", "{query,entry-1-")
	testEntryPostMutation(conn, Mn(":query"))
	to.MatchWait(t, 200, "trace", "entry-1-", "nop", "{:query,,<nil>,<nil>}")
	conn.Close()
	to.MatchWait(t, 200, "trace", "disp", "{:remove,entry-1-")
	callback(Mns(":remove", "tid"))
	to.MatchWait(t, 200, "trace", "entry-1-", "out", "{:remove,tid,<nil>,<nil>}")
}

func testEntryConnect(port int, path string) *websocket.Conn {
	url := fmt.Sprintf("ws://localhost:%v%v", port, path)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	PanicIfError(err)
	return conn
}

func testEntryReadLoop(conn *websocket.Conn, to TestOutput) {
	for {
		mut := testEntryReadMutation(conn)
		if mut.Name == "" {
			return
		}
		to.Trace("test", "read", mut)
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

func testEntryPostMutation(conn *websocket.Conn, mut Mutation) {
	bytes, err := EncodeMutation(mut)
	PanicIfError(err)
	err = conn.WriteMessage(websocket.TextMessage, bytes)
	PanicIfError(err)
}
