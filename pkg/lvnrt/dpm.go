package lvnrt

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"time"
)

type Dpm interface {
	Close() Channel
	Port() int
	Echo()
}

type dpmDso struct {
	id      Id
	delay   int
	log     Logger
	port    int
	cleaner Cleaner
	listen  net.Listener
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
	te.port = listen.Addr().(*net.TCPAddr).Port
	te.id = NewId("dpm")
	te.delay = delay
	te.listen = listen
	go te.aloop()
	return te
}

func (dpm *dpmDso) Port() int {
	return dpm.port
}

func (dpm *dpmDso) Close() Channel {
	dpm.cleaner.Close()
	done := make(Channel)
	dpm.cleaner.AddChannel("done", done)
	return done
}

func (dpm *dpmDso) Echo() {
	address := fmt.Sprintf("127.0.0.1:%v", dpm.Port())
	socket := NewSocket(address, 400)
	defer socket.Close()
	req := "*1B1"
	err := socket.WriteLine(req, 400)
	PanicIfError(err)
	res, err := socket.ReadLine(400 + dpm.delay)
	PanicIfError(err)
	AssertTrue(req == res, "mismatch", req, res)
}

func (dpm *dpmDso) aloop() {
	dpm.cleaner.AddCloser("accept", dpm.listen)
	defer dpm.cleaner.Remove("accept")
	for {
		conn, err := dpm.listen.Accept()
		if err != nil {
			return
		}
		cid := dpm.id.Next()
		go dpm.cloop(cid, conn)
	}
}

func (dpm *dpmDso) cloop(cid string, conn net.Conn) {
	defer TraceRecover(dpm.log.Trace)
	cid += "-" + conn.RemoteAddr().String()
	dpm.cleaner.AddCloser(cid, conn)
	defer dpm.cleaner.Remove(cid)
	cr := byte(13)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		req, err := reader.ReadString(cr)
		if err != nil {
			return
		}
		echos := dpm.echos.MatchString(req)
		dpm.log.Trace(echos, Readable(req))
		if echos {
			time.Sleep(Millis(dpm.delay))
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
