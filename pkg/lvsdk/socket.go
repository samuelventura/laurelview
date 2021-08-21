package lvsdk

import (
	"fmt"
	"net"
	"strings"
)

type Socket interface {
	WriteLine(req string, toms int64) error
	ReadLine(toms int64) (string, error)
	Discard(toms int64) error
	Close()
}

type socketDso struct {
	conn net.Conn
}

func NewSocketConn(conn net.Conn) Socket {
	s := &socketDso{}
	s.conn = conn
	//connection drop detection macos=~9s
	//set data asap, do not wait for larger packet
	tcp := conn.(*net.TCPConn)
	tcp.SetNoDelay(true)
	tcp.SetLinger(0)
	tcp.SetKeepAlive(true)
	tcp.SetKeepAlivePeriod(Millis(1000))
	tcp.SetWriteBuffer(0)
	return s
}

func NewSocket(address string, toms int) Socket {
	s := &socketDso{}
	to := Millis(int64(toms))
	conn, err := net.DialTimeout("tcp", address, to)
	PanicIfError(err)
	s.conn = conn
	return s
}

func (s *socketDso) Close() {
	s.conn.Close()
}

func (s *socketDso) Discard(toms int64) error {
	err := s.conn.SetReadDeadline(Future(toms))
	if err != nil {
		return err
	}
	bytes := make([]byte, 512)
	_, err = s.conn.Read(bytes)
	for err == nil {
		_, err = s.conn.Read(bytes)
	}
	nerr, ok := err.(net.Error)
	if ok && nerr.Timeout() {
		return nil
	}
	return err
}

func (s *socketDso) WriteLine(req string, toms int64) error {
	err := s.conn.SetWriteDeadline(Future(toms))
	if err != nil {
		return err
	}
	bytes := []byte(req)
	n, err := s.conn.Write(bytes)
	if err == nil && n != len(bytes) {
		err = fmt.Errorf("wrote %v of %v", n, len(bytes))
	}
	if err != nil {
		return err
	}
	bytes = []byte{byte(13)}
	n, err = s.conn.Write(bytes)
	if err == nil && n != len(bytes) {
		err = fmt.Errorf("wrote %v of %v", n, len(bytes))
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *socketDso) ReadLine(toms int64) (string, error) {
	err := s.conn.SetReadDeadline(Future(toms))
	if err != nil {
		return "", err
	}
	cr := byte(13)
	buf := new(strings.Builder)
	bytes := []byte{cr}
	for {
		n, err := s.conn.Read(bytes)
		if err == nil && n != 1 {
			err = fmt.Errorf("read %v of %v", n, 1)
		}
		if err != nil {
			return "", err
		}
		b := bytes[0]
		if b == cr {
			break
		}
		buf.WriteByte(b)
	}
	return buf.String(), nil
}
