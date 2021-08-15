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
	log      Logger
	rt       Runtime
	listener net.Listener
	upgrader websocket.FastHTTPUpgrader
}

type clientDso struct {
	sid      string
	log      Logger
	rt       Runtime
	conn     *websocket.Conn
	callback chan *Mutation
}

func NewEntry(rt Runtime, id Id, endpoint string) Entry {
	listener, err := net.Listen("tcp", endpoint)
	PanicIfError(err)
	entry := &entryDso{}
	entry.id = id
	entry.rt = rt
	entry.log = rt.PrefixLog("entry")
	entry.listener = listener
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
	PanicIfError(err)
}

func (entry *entryDso) listen() {
	defer TraceRecover(entry.log.Debug)
	defer entry.listener.Close()
	//ignore accept close error on exit
	fasthttp.Serve(entry.listener, entry.handle)
}

func (entry *entryDso) origin(ctx *fasthttp.RequestCtx) bool {
	return true
}

func (entry *entryDso) handle(ctx *fasthttp.RequestCtx) {
	err := entry.upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		defer TraceRecover(entry.log.Debug)
		defer conn.Close()
		id := entry.id.Next()
		ipp := conn.RemoteAddr().String()
		client := &clientDso{}
		client.conn = conn
		client.rt = entry.rt
		client.sid = fmt.Sprintf("%v_%v", id, ipp)
		client.log = entry.rt.PrefixLog(client.sid)
		client.callback = make(chan *Mutation)
		client.loop()
	})
	if err != nil {
		entry.log.Error("upgrade", err)
		return
	}
}

func (client *clientDso) loop() {
	mut, err := client.read()
	if err != nil {
		client.log.Error(err)
		return
	}
	_, ok := mut.Args.(*SetupArgs)
	if !ok {
		client.log.Error("setup", mut)
		return
	}
	defer client.wait()
	defer client.remove()
	client.add()
	client.rt.Post("state", mut)
	go client.reader()
	for mutation := range client.callback {
		bytes := encodeMutation(mutation)
		err := client.conn.WriteMessage(websocket.TextMessage, bytes)
		PanicIfError(err)
	}
}

func (client *clientDso) remove() {
	mut := &Mutation{}
	mut.Sid = client.sid
	mut.Name = "remove"
	client.rt.Post("state", mut)
}

func (client *clientDso) add() {
	args := &AddArgs{}
	args.Callback = client.writer
	mut := &Mutation{}
	mut.Sid = client.sid
	mut.Name = "add"
	mut.Args = args
	client.rt.Post("state", mut)
}

func (client *clientDso) wait() {
	for range client.callback {
	}
}

//FIXME is remove/dispose being received? test it
func (client *clientDso) writer(mut *Mutation) {
	defer TraceRecover(client.log.Debug)
	client.log.Trace("out", mut)
	switch mut.Name {
	case "setup", "query":
		client.callback <- mut
	case "remove", "dispose":
		close(client.callback)
	default:
		client.log.Debug("unknown mutation", mut.Name)
	}
}

func (client *clientDso) reader() {
	defer TraceRecover(client.log.Debug)
	defer client.conn.Close()
	defer client.remove()
	for {
		mut, err := client.read()
		if err != nil {
			client.log.Trace(err)
			return
		}
		client.rt.Post("state", mut)
	}
}

func (client *clientDso) read() (mut *Mutation, err error) {
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
	//FIXME test needed to ensure incomming misformats properly handled
	mut, err = decodeMutation(msg)
	if err != nil {
		err = fmt.Errorf("decodeMutation %w", err)
		return
	}
	mut.Sid = client.sid
	client.log.Trace("in", mut)
	return
}
