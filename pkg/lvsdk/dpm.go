package lvsdk

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"time"
)

type Dpm interface {
	Port() uint
	Echo()
	Close()
}

type dpmDso struct {
	id      Id
	log     Logger
	listen  net.Listener
	cleaner Cleaner
	done    Channel
	delay   int
	echos   *regexp.Regexp
}

func NewDpm(log Logger, address string, delay int) Dpm {
	te := &dpmDso{}
	listen, err := net.Listen("tcp", address)
	PanicIfError(err)
	te.log = PrefixLogger(log.Log, "dpm")
	clog := PrefixLogger(log.Log, "dpm", "cleaner")
	te.cleaner = NewCleaner(clog)
	te.echos = regexp.MustCompile(`^\*.B\d\r$`)
	te.id = NewId("dpm")
	te.done = make(Channel)
	te.delay = delay
	te.listen = listen
	go te.aloop()
	return te
}

func (e *dpmDso) Port() uint {
	return uint(e.listen.Addr().(*net.TCPAddr).Port)
}

func (e *dpmDso) Close() {
	e.listen.Close()
	<-e.done //wait aloop
	e.cleaner.Close()
}

func (e *dpmDso) Echo() {
	address := fmt.Sprintf("127.0.0.1:%v", e.Port())
	socket := NewSocket(address, 400)
	defer socket.Close()
	req := "*1B1"
	err := socket.WriteLine(req, 400)
	PanicIfError(err)
	res, err := socket.ReadLine(400 + e.delay)
	PanicIfError(err)
	AssertTrue(req == res, "mismatch", req, res)
}

func (e *dpmDso) aloop() {
	e.cleaner.AddChannel("accept", e.done)
	defer e.cleaner.Remove("accept")
	for {
		conn, err := e.listen.Accept()
		if err != nil {
			return
		}
		cid := e.id.Next()
		go e.cloop(cid, conn)
	}
}

func (e *dpmDso) cloop(cid string, conn net.Conn) {
	defer TraceRecover(e.log.Trace)
	cid += "-" + conn.RemoteAddr().String()
	e.cleaner.AddCloser(cid, conn)
	defer e.cleaner.Remove(cid)
	cr := byte(13)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		req, err := reader.ReadString(cr)
		if err != nil {
			return
		}
		echos := e.echos.MatchString(req)
		e.log.Trace(echos, Readable(req))
		if echos {
			time.Sleep(Millis(e.delay))
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
}
