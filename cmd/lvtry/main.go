package main

import (
	"bufio"
	"net"
)

func main() {
	tryout1(GoLogLogger())
}

func tryout1(log Logger) {
	defer TraceRecover(log.Info)
	log.Info("Connecting...")
	//192.168.1.77 nucmeg
	//192.168.1.82 fedora
	conn, err := net.Dial("tcp", "192.168.1.77:5000")
	PanicIfError(err)
	tcp := conn.(*net.TCPConn)
	tcp.SetNoDelay(true)
	tcp.SetLinger(0)
	tcp.SetKeepAlive(true)
	tcp.SetKeepAlivePeriod(Millis(1000))
	tcp.SetWriteBuffer(0)
	log.Info("Connected")
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString(byte(10))
		PanicIfError(err)
		log.Info(line[:len(line)-1])
	}
}
