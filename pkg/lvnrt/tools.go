package lvnrt

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"time"
	"unicode"
)

func traceRecover(output Output) {
	r := recover()
	if r != nil {
		output("recover", r, string(debug.Stack()))
	}
}

func traceIfError(output Output, err error) {
	if err != nil {
		output("error", err)
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
		if ok {
			log.Trace(mut)
			dispatch(mut)
		} else {
			log.Debug(mut)
		}
	}
}

func asyncDispatch(output Output, dispatch Dispatch) Dispatch {
	queue := make(chan *Mutation)
	loop := func() {
		defer traceRecover(output)
		for mut := range queue {
			dispatch(mut)
		}
	}
	go loop()
	return func(mut *Mutation) {
		//do not close queue nor state dispose
		//let map dispatch report the ignore
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

func sendChannel(channel Channel, any Any) {
	channel <- any
}

func closeChannel(channel Channel) {
	select {
	case <-channel:
	default:
		close(channel)
	}
}

func waitChannel(channel Channel, output Output) {
	output("waiting channel...")
	<-channel
	output("waiting channel done")
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

func disposeArgs(arg Any) {
	action, ok := arg.(Action)
	if ok {
		action()
	}
	channel, ok := arg.(Channel)
	if ok {
		close(channel)
	}
}

type testEcho interface {
	port() uint
	ping()
	close()
}

type testEchoDso struct {
	log    Logger
	listen net.Listener
	mutex  *sync.Mutex
	conns  *list.List
	done   Channel
}

func newTestEcho(log Logger) testEcho {
	te := &testEchoDso{}
	listen, err := net.Listen("tcp", ":0")
	panicIfError(err)
	te.log = prefixLogger(log.Log, "echo")
	te.conns = list.New()
	te.mutex = new(sync.Mutex)
	te.done = make(Channel)
	te.listen = listen
	go te.loop()
	return te
}

func (te *testEchoDso) port() uint {
	return uint(te.listen.Addr().(*net.TCPAddr).Port)
}

func (te *testEchoDso) close() {
	te.listen.Close()
	<-te.done
}

func (te *testEchoDso) ping() {
	address := fmt.Sprintf("127.0.0.1:%v", te.port())
	conn, err := net.DialTimeout("tcp", address, millis(400))
	panicIfError(err)
	defer conn.Close()
	n, err := conn.Write([]byte{13})
	panicIfError(err)
	m, err := conn.Read([]byte{13})
	panicIfError(err)
	assertTrue(n == 1, "Wrote n", n)
	assertTrue(m == 1, "Read m", m)
}

func (te *testEchoDso) loop() {
	count := 0
	done := make(Channel)
	defer closeChannel(te.done)
	defer te.wait(count, done)
	defer te.clear()
	for {
		conn, err := te.listen.Accept()
		if err != nil {
			return
		}
		count++
		go te.echo(conn, done)
	}
}

func (te *testEchoDso) wait(count int, done Channel) {
	for count > 0 {
		te.log.Trace("wait", "count", count)
		<-done
		count--
	}
	te.log.Trace("wait", "done")
}

func (te *testEchoDso) clear() {
	defer te.mutex.Unlock()
	te.mutex.Lock()
	e := te.conns.Front()
	for e != nil {
		conn := e.Value.(net.Conn)
		defer conn.Close()
		e = e.Next()
	}
}

func (te *testEchoDso) add(conn net.Conn) *list.Element {
	defer te.mutex.Unlock()
	te.mutex.Lock()
	return te.conns.PushBack(conn)
}

func (te *testEchoDso) remove(e *list.Element) {
	defer te.mutex.Unlock()
	te.mutex.Lock()
	conn := e.Value.(net.Conn)
	defer conn.Close()
	te.conns.Remove(e)
}

func (te *testEchoDso) echo(conn net.Conn, done Channel) {
	defer sendChannel(done, true)
	e := te.add(conn)
	defer te.remove(e)
	cr := byte(13)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		req, err := reader.ReadString(cr)
		if err != nil {
			return
		}
		te.log.Trace(readable(req))
		//FIXME do not reply resets
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
