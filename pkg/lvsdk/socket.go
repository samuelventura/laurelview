package lvsdk

import (
	"fmt"
	"net"
	"strings"
)

type Socket interface {
	WriteLine(req string, toms int) error
	ReadLine(toms int) (string, error)
	Discard(toms int) error
	Close()
}

type socketDso struct {
	conn    net.Conn
	sep     byte
	discard []byte
	input   *strings.Builder
}

func NewSocketConn(conn net.Conn, sep byte) Socket {
	s := &socketDso{}
	s.conn = conn
	s.sep = sep
	s.discard = make([]byte, 512)
	s.input = new(strings.Builder)
	return s
}

func NewSocketDial(address string, toms int, sep byte) Socket {
	to := Millis(toms)
	conn, err := net.DialTimeout("tcp", address, to)
	PanicIfError(err)
	return NewSocketConn(conn, sep)
}

func (s *socketDso) Close() {
	s.conn.Close()
}

func (s *socketDso) Discard(toms int) error {
	err := s.conn.SetReadDeadline(Future(toms))
	if err != nil {
		return err
	}
	_, err = s.conn.Read(s.discard)
	for err == nil {
		_, err = s.conn.Read(s.discard)
	}
	nerr, ok := err.(net.Error)
	if ok && nerr.Timeout() {
		return nil
	}
	return err
}

func (s *socketDso) WriteLine(req string, toms int) error {
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
	bytes = []byte{s.sep}
	n, err = s.conn.Write(bytes)
	if err == nil && n != len(bytes) {
		err = fmt.Errorf("wrote %v of %v", n, len(bytes))
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *socketDso) ReadLine(toms int) (string, error) {
	err := s.conn.SetReadDeadline(Future(toms))
	if err != nil {
		return "", err
	}
	s.input.Reset()
	bytes := []byte{s.sep}
	for {
		n, err := s.conn.Read(bytes)
		if err == nil && n != 1 {
			err = fmt.Errorf("read %v of %v", n, 1)
		}
		if err != nil {
			return "", err
		}
		b := bytes[0]
		if b == s.sep {
			break
		}
		s.input.WriteByte(b)
	}
	return s.input.String(), nil
}
