package lvnbe

import (
	"fmt"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/stretchr/testify/assert"
)

func TestEntry(t *testing.T) {
	out := newTestOutput()
	defer out.close()
	var dao = NewDao(":memory:")
	var state = NewState(dao)
	var hub = NewHub(state)
	var core = NewCore(hub, out.out)
	var entry = NewEntry(core, out.out, ":0")
	defer entry.Close()
	out.out("trace", "Entry port", entry.Port())
	conn := connect(entry.Port())
	out.out("trace", "Client", conn.LocalAddr())
	all := readMutation(conn)
	out.out("trace", "all", all)
	assert.Equal(t, "all", all.Name)
	assert.Equal(t, 0, len(all.Args.(*AllArgs).Items))

	postCreate(conn, "name1", "json1")
	create1 := readMutation(conn)
	id1 := create1.Args.(*CreateArgs).Id
	out.out("trace", "create1", create1)
	assert.Equal(t, "create", create1.Name)
	assert.Equal(t, "name1", create1.Args.(*CreateArgs).Name)
	assert.Equal(t, "json1", create1.Args.(*CreateArgs).Json)
	assert.Equal(t, uint(1), id1)

	postCreate(conn, "name2", "json2")
	create2 := readMutation(conn)
	id2 := create2.Args.(*CreateArgs).Id
	out.out("trace", "create2", create2)
	assert.Equal(t, "create", create2.Name)
	assert.Equal(t, "name2", create2.Args.(*CreateArgs).Name)
	assert.Equal(t, "json2", create2.Args.(*CreateArgs).Json)
	assert.Equal(t, uint(2), id2)

	postDelete(conn, id2)
	delete := readMutation(conn)
	out.out("trace", "delete", delete)
	assert.Equal(t, "delete", delete.Name)
	assert.Equal(t, id2, delete.Args.(*DeleteArgs).Id)

	postUpdate(conn, id1, "name3", "json3")
	update := readMutation(conn)
	out.out("trace", "update", update)
	assert.Equal(t, "update", update.Name)
	assert.Equal(t, id1, update.Args.(*UpdateArgs).Id)
	assert.Equal(t, "name3", update.Args.(*UpdateArgs).Name)
	assert.Equal(t, "json3", update.Args.(*UpdateArgs).Json)
}

func connect(port int) *websocket.Conn {
	url := fmt.Sprintf("ws://localhost:%v/ws/index", port)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	PanicIfError(err)
	return conn
}

func readMutation(conn *websocket.Conn) *Mutation {
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 400))
	mt, bytes, err := conn.ReadMessage()
	PanicIfError(err)
	if mt != websocket.TextMessage {
		PanicF("Invalid msg type %v", mt)
	}
	mut, err := decodeMutation(bytes)
	PanicIfError(err)
	return mut
}

func postCreate(conn *websocket.Conn, name string, json string) {
	args := &CreateArgs{}
	args.Json = json
	args.Name = name
	mut := &Mutation{Name: "create", Args: args}
	bytes := encodeMutation(mut)
	err := conn.WriteMessage(websocket.TextMessage, bytes)
	PanicIfError(err)
}

func postUpdate(conn *websocket.Conn, id uint, name string, json string) {
	args := &UpdateArgs{}
	args.Id = id
	args.Name = name
	args.Json = json
	mut := &Mutation{Name: "update", Args: args}
	bytes := encodeMutation(mut)
	err := conn.WriteMessage(websocket.TextMessage, bytes)
	PanicIfError(err)
}

func postDelete(conn *websocket.Conn, id uint) {
	args := &DeleteArgs{}
	args.Id = id
	mut := &Mutation{Name: "delete", Args: args}
	bytes := encodeMutation(mut)
	err := conn.WriteMessage(websocket.TextMessage, bytes)
	PanicIfError(err)
}
