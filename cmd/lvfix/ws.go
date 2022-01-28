package main

import (
	"time"

	"github.com/sacOO7/gowebsocket"
)

func printWebSocket(scheme string, ip string, seconds int) {

	socket := gowebsocket.New(scheme + "://" + ip + "/ws/db")
	socket.Timeout = time.Second * 1
	socket.OnConnected = func(socket gowebsocket.Socket) {

	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		writeFile(seconds, "Fail")
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		writeFile(seconds, "Pass")
		socket.Close()
	}

	socket.Connect()

	return
}
