package lvnrt

import (
	"bufio"
	"fmt"
	"net"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
	"unicode"
)

func traceRecover(output Output) {
	r := recover()
	if r != nil {
		output("trace", "recover", r, string(debug.Stack()))
	}
}

func traceIfError(output Output, err error) {
	if err != nil {
		output("trace", "error", err)
	}
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func panicF(format string, args ...Any) {
	panic(fmt.Errorf(format, args...))
}

func panicLN(args ...Any) {
	panic(fmt.Sprintln(args...))
}

func assertTrue(flag bool, args ...Any) {
	if !flag {
		panic(fmt.Sprintln(args...))
	}
}

func readable(s string) string {
	b := new(strings.Builder)
	for _, c := range s {
		if unicode.IsControl(c) || unicode.IsSpace(c) {
			h := fmt.Sprintf("[%02X]", int(c))
			b.WriteString(h)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func clearDispatch(dispatchs map[string]Dispatch) {
	for name := range dispatchs {
		delete(dispatchs, name)
	}
}

func mapDispatch(log Logger, dispmap map[string]Dispatch) Dispatch {
	return func(mut *Mutation) {
		dispatch, ok := dispmap[mut.Name]
		if !ok {
			log.Warn("unknown mutation", mut.Name)
			return
		}
		log.Trace(mut.Name, mut.Sid, toMap(mut.Args))
		dispatch(mut)
	}
}

func asyncDispatch(log Log, dispatch Dispatch) Dispatch {
	logger := prefixLogger(log, "async")
	queue := make(chan *Mutation)
	dispose := func(string) {}
	dispose = func(name string) {
		if name == "dispose" {
			dispose = func(string) {}
		}
	}
	loop := func() {
		defer traceRecover(logger.Warn)
		for mut := range queue {
			dispatch(mut)
		}
	}
	go loop()
	return func(mut *Mutation) {
		defer dispose(mut.Name)
		queue <- mut
	}
}

func millis(ms int64) time.Duration {
	return time.Duration(ms) * time.Millisecond
}

func future(ms int64) time.Time {
	d := millis(ms)
	return time.Now().Add(d)
}

func send(channel Channel, any Any) {
	channel <- any
}

func toMap(any Any) Map {
	m := make(Map)
	e := reflect.ValueOf(any).Elem()
	t := e.Type()
	m["$type"] = t.Name()
	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		ft := t.Field(i)
		m[ft.Name] = f.Interface()
	}
	return m
}

type testEcho interface {
	port() uint
	close()
}

type testEchoDso struct {
	log    Logger
	listen net.Listener
}

func newTestEcho(log Logger) testEcho {
	te := &testEchoDso{}
	listen, err := net.Listen("tcp", ":0")
	panicIfError(err)
	te.log = log
	te.listen = listen
	go te.loop()
	return te
}

func (te *testEchoDso) port() uint {
	return uint(te.listen.Addr().(*net.TCPAddr).Port)
}

func (te *testEchoDso) close() {
	te.listen.Close()
}

func (te *testEchoDso) loop() {
	for {
		conn, err := te.listen.Accept()
		if err != nil {
			return
		}
		go te.echo(conn)
	}
}

func (te *testEchoDso) echo(conn net.Conn) {
	defer conn.Close()
	cr := byte(13)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		req, err := reader.ReadString(cr)
		if err != nil {
			return
		}
		te.log.Trace("echo", readable(req))
		_, err = writer.WriteString(req)
		if err != nil {
			return
		}
		err = writer.Flush()
		if err != nil {
			return
		}
	}
}
