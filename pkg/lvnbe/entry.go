package lvnbe

import (
	"embed"
	"fmt"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

//go:embed build/*
var fs embed.FS

type Entry interface {
	Port() int
	Close()
}

type entryDso struct {
	port     int
	core     Core
	output   Output
	listener net.Listener
	upgrader websocket.FastHTTPUpgrader
	cache    map[string]*staticDso
}

type staticDso struct {
	bytes []byte
	mime  string
}

type clientDso struct {
	core     Core
	sid      string
	output   Output
	conn     *websocket.Conn
	callback chan *Mutation
}

func NewEntry(core Core, output Output, endpoint string) Entry {
	listener, err := net.Listen("tcp", endpoint)
	PanicIfError(err)
	entry := &entryDso{}
	entry.core = core
	entry.output = output
	entry.listener = listener
	entry.cache = make(map[string]*staticDso)
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
	defer TraceRecover(entry.output)
	defer entry.listener.Close()
	//ignore accept close error on exit
	fasthttp.Serve(entry.listener, entry.handle)
}

func (entry *entryDso) origin(ctx *fasthttp.RequestCtx) bool {
	return true
}

func (entry *entryDso) handle(ctx *fasthttp.RequestCtx) {
	var err error
	path := string(ctx.Path())
	if path == "/" {
		path = "/index.html"
	}
	ext := filepath.Ext(path)
	//FIXME *.map gets empty mime
	//FIXME add client side standard http file caching
	ct := mime.TypeByExtension(ext)
	entry.output("trace", path, ct)
	if !strings.HasPrefix(path, "/ws/") {
		path = "build" + path
		static, ok := entry.cache[path]
		if !ok {
			bytes, err := fs.ReadFile(path)
			if err != nil {
				entry.output("trace", path, err)
				ctx.Response.SetStatusCode(http.StatusNotFound)
				return
			}
			static = &staticDso{}
			entry.cache[path] = static
			static.bytes = bytes
			static.mime = ct
		}
		ctx.SetContentType(static.mime)
		ctx.Write(static.bytes)
		return
	}
	err = entry.upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		defer TraceRecover(entry.output)
		defer conn.Close()
		id := entry.core.NextId()
		ipp := conn.RemoteAddr().String()
		client := &clientDso{}
		client.conn = conn
		client.core = entry.core
		client.output = entry.output
		client.sid = fmt.Sprintf("%v_%v", id, ipp)
		client.callback = make(chan *Mutation)
		defer client.wait()
		client.loop()
	})
	if err != nil {
		entry.output("trace", "upgrade:", err)
		return
	}
}

func (client *clientDso) loop() {
	defer client.core.Remove(client.sid)
	client.core.Add(client.sid, client.writer)
	go client.reader()
	for mutation := range client.callback {
		bytes := encodeMutation(mutation)
		err := client.conn.WriteMessage(websocket.TextMessage, bytes)
		PanicIfError(err)
	}
}

func (client *clientDso) wait() {
	for range client.callback {
	}
}

func (client *clientDso) writer(mutation *Mutation) {
	//closing a closed channel panics
	defer TraceRecover(client.output)
	switch mutation.Name {
	case "all", "create", "delete", "update":
		client.callback <- mutation
	case "remove", "close":
		close(client.callback)
	default:
		client.output("trace", "unknown mutation", mutation.Name)
	}
}

func (client *clientDso) reader() {
	defer TraceRecover(client.output)
	defer client.conn.Close()
	defer client.core.Remove(client.sid)
	for {
		mutation, err := client.read()
		if err != nil {
			client.output("trace", err)
			return
		}
		mutation.Sid = client.sid
		client.core.Apply(mutation)
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
