package lvsdk

import (
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

type Entry interface {
	Port() int
	Status(Channel)
	Close() Channel
}

type entryDso struct {
	port     int
	id       Id
	log      Logger
	ctx      Context
	buflen   int
	rtoms    int
	wtoms    int
	endpoint string
	cleaner  Cleaner
	listener net.Listener
	upgrader websocket.FastHTTPUpgrader
	static   Handler
}

type clientDso struct {
	sid      string
	log      Logger
	rtoms    int
	wtoms    int
	conn     *websocket.Conn
	dispatch Dispatch
	callback chan Mutation
}

func NewEntry(ctx Context) Entry {
	buflen := ctx.GetValue("entry.buflen").(int)
	wtoms := ctx.GetValue("entry.wtoms").(int)
	rtoms := ctx.GetValue("entry.rtoms").(int)
	static := ctx.GetValue("entry.static").(Handler)
	endpoint := ctx.GetValue("entry.endpoint").(string)
	listener, err := net.Listen("tcp", endpoint)
	PanicIfError(err)
	cleaner := NewCleaner(ctx.PrefixLog("entry", "cleaner"))
	entry := &entryDso{}
	entry.id = NewId("entry")
	entry.ctx = ctx
	entry.buflen = buflen
	entry.wtoms = wtoms
	entry.rtoms = rtoms
	entry.static = static
	entry.endpoint = endpoint
	entry.log = ctx.PrefixLog("entry")
	entry.cleaner = cleaner
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

func (entry *entryDso) Status(output Channel) {
	//do not pass by reference to log
	entry.cleaner.Status(func(any Any) {
		defer close(output)
		c := any.(*cleanerDso)
		output <- len(c.items)
		keys := make([]string, 0, len(c.items))
		for k := range c.items {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			output <- k
		}
	})
}

func (entry *entryDso) Close() Channel {
	entry.cleaner.Close()
	done := make(Channel)
	entry.cleaner.AddChannel("done", done)
	return done
}

func (entry *entryDso) listen() {
	defer TraceRecover(entry.log.Error)
	defer entry.cleaner.Remove("listen")
	entry.cleaner.AddCloser("listen", entry.listener)
	fasthttp.Serve(entry.listener, entry.handle)
	//ignore accept close error on exit
}

func (entry *entryDso) origin(ctx *fasthttp.RequestCtx) bool {
	//FIXME authentication token?
	return true
}

func (entry *entryDso) handle(ctx *fasthttp.RequestCtx) {
	defer TraceRecover(entry.log.Debug)
	path := string(ctx.Path())
	if !strings.HasPrefix(path, "/ws/") {
		entry.static(ctx)
		return
	}
	err := entry.upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		defer TraceRecover(entry.log.Debug)
		defer conn.Close()
		id := entry.id.Next()
		dispatch := entry.ctx.GetDispatch(path) //panics
		ipp := conn.RemoteAddr().String()
		sid := fmt.Sprintf("%s_%s_%s", id, ipp, path)
		log := entry.ctx.PrefixLog(sid)
		log.Trace("path", path)
		defer entry.cleaner.Remove(sid)
		entry.cleaner.AddCloser(sid, conn)
		client := &clientDso{}
		client.log = log
		client.sid = sid
		client.conn = conn
		client.wtoms = entry.wtoms
		client.rtoms = entry.rtoms
		client.dispatch = dispatch
		client.callback = make(chan Mutation, entry.buflen)
		client.loop()
	})
	if err != nil {
		entry.log.Debug("upgrade", err)
		return
	}
}

func (client *clientDso) loop() {
	defer client.wait()
	defer client.remove()
	client.add()
	go client.reader()
	mt := websocket.TextMessage
	for mut := range client.callback {
		bytes, err := EncodeMutation(mut)
		PanicIfError(err)
		if client.wtoms > 0 {
			wdl := Future(client.wtoms)
			err = client.conn.SetWriteDeadline(wdl)
			PanicIfError(err)
		}
		err = client.conn.WriteMessage(mt, bytes)
		PanicIfError(err)
	}
}

func (client *clientDso) add() {
	mut := Mutation{}
	mut.Sid = client.sid
	mut.Name = ":add"
	mut.Args = client.writer
	client.dispatch(mut)
}

func (client *clientDso) remove() {
	mut := Mutation{}
	mut.Sid = client.sid
	mut.Name = ":remove"
	client.dispatch(mut)
}

func (client *clientDso) wait() {
	for range client.callback {
	}
}

func (client *clientDso) writer(mut Mutation) {
	defer TraceRecover(client.log.Debug)
	client.log.Trace("out", mut)
	switch mut.Name {
	case ":remove", ":dispose":
		close(client.callback)
	default:
		client.callback <- mut
	}
}

func (client *clientDso) reader() {
	defer TraceRecover(client.log.Debug)
	defer client.conn.Close()
	defer client.remove()
	for {
		mut, err := client.read()
		if err != nil {
			client.log.Debug(err)
			return
		}
		if !strings.HasPrefix(mut.Name, ":") {
			client.log.Trace("in", mut)
			mut.Sid = client.sid
			client.dispatch(mut)
		} else {
			client.log.Trace("nop", mut)
		}
	}
}

func (client *clientDso) read() (mut Mutation, err error) {
	if client.rtoms > 0 {
		rdl := Future(client.rtoms)
		client.conn.SetReadDeadline(rdl)
	}
	mt, msg, err := client.conn.ReadMessage()
	if err != nil {
		err = fmt.Errorf("conn.ReadMessage %w", err)
		return
	}
	if websocket.TextMessage != mt {
		err = fmt.Errorf("websocket.TextMessage != %v", mt)
		return
	}
	mut, err = DecodeMutation(msg)
	if err != nil {
		err = fmt.Errorf("decode %w", err)
		return
	}
	return
}
