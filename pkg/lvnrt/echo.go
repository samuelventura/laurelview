package lvnrt

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"sync"
)

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
