package lvnrt

import (
	"fmt"
	"net"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

type Entry interface {
	Port() int
	Close()
}

type entryDso struct {
	port     int
	id       Id
	dispatch Dispatch
	output   Output
	listener net.Listener
	upgrader websocket.FastHTTPUpgrader
}

type clientDso struct {
	output   Output
	sid      string
	dispatch Dispatch
	conn     *websocket.Conn
	callback chan *Mutation
}

func NewEntry(output Output, dispatch Dispatch, id Id, endpoint string) Entry {
	listener, err := net.Listen("tcp", endpoint)
	panicIfError(err)
	entry := &entryDso{}
	entry.id = id
	entry.output = output
	entry.listener = listener
	entry.dispatch = dispatch
	entry.port = listener.Addr().(*net.TCPAddr).Port
	entry.upgrader = websocket.FastHTTPUpgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     entry.origin,
	}
	go entry.listen()
	return entry
}

func (entry *entryDso) Port() int {
	return entry.port
}

func (entry *entryDso) Close() {
	//hub may have multiple entries
	err := entry.listener.Close()
	panicIfError(err)
}

func (entry *entryDso) listen() {
	defer traceRecover(entry.output)
	defer entry.listener.Close()
	//ignore accept close error on exit
	fasthttp.Serve(entry.listener, entry.handle)
}

func (entry *entryDso) origin(ctx *fasthttp.RequestCtx) bool {
	return true
}

func (entry *entryDso) handle(ctx *fasthttp.RequestCtx) {
	err := entry.upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		defer traceRecover(entry.output)
		defer conn.Close()
		id := entry.id.Next()
		ipp := conn.RemoteAddr().String()
		client := &clientDso{}
		client.conn = conn
		client.output = entry.output
		client.dispatch = entry.dispatch
		client.sid = fmt.Sprintf("%v_%v", id, ipp)
		client.callback = make(chan *Mutation)
		client.loop()
	})
	if err != nil {
		entry.output("trace", "upgrade:", err)
		return
	}
}

func (client *clientDso) loop() {
	mut, err := client.read()
	if err != nil {
		client.output("trace", err)
		return
	}
	_, ok := mut.Args.(*SetupArgs)
	if !ok {
		client.output("trace", "setup expected", mut)
		return
	}
	mut.Sid = client.sid
	defer client.wait()
	defer client.remove()
	client.add()
	go client.reader()
	for mutation := range client.callback {
		bytes := encodeMutation(mutation)
		err := client.conn.WriteMessage(websocket.TextMessage, bytes)
		panicIfError(err)
	}
}

func (client *clientDso) remove() {
	mut := &Mutation{}
	mut.Sid = client.sid
	mut.Name = "remove"
	client.dispatch(mut)
}

func (client *clientDso) add() {
	args := &AddArgs{}
	args.Callback = client.writer
	mut := &Mutation{}
	mut.Sid = client.sid
	mut.Name = "add"
	mut.Args = args
	client.dispatch(mut)
}

func (client *clientDso) wait() {
	for range client.callback {
	}
}

func (client *clientDso) writer(mutation *Mutation) {
	//closing a closed channel panics
	defer traceRecover(client.output)
	switch mutation.Name {
	case "setup", "query":
		client.callback <- mutation
	case "remove", "dispose":
		close(client.callback)
	default:
		client.output("trace", "unknown mutation", mutation.Name)
	}
}

func (client *clientDso) reader() {
	defer traceRecover(client.output)
	defer client.conn.Close()
	defer client.remove()
	for {
		mutation, err := client.read()
		if err != nil {
			client.output("trace", err)
			return
		}
		mutation.Sid = client.sid
		client.dispatch(mutation)
	}
}

func (client *clientDso) read() (mutation *Mutation, err error) {
	mt, msg, err := client.conn.ReadMessage()
	if err != nil {
		err = fmt.Errorf("conn.ReadMessage %w", err)
		return
	}
	if websocket.TextMessage != mt {
		err = fmt.Errorf("websocket.TextMessage !=%v", mt)
		return
	}
	//may throw on invalid json format
	mutation, err = decodeMutation(msg)
	if err != nil {
		err = fmt.Errorf("decodeMutation %w", err)
		return
	}
	return
}
