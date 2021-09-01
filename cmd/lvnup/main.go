package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/fasthttp/websocket"
)

func main() {
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	log := DefaultLog()
	defer CloseLog(log)

	logger := PrefixLogger(log)
	ep := os.Getenv("LV_NBE_ENDPOINT")
	logger.Info("LV_NBE_ENDPOINT", ep)
	go run(logger)

	//wait
	exit := make(Channel)
	go stdin(exit)
	select {
	case <-ctrlc:
	case <-exit:
	}
}

func stdin(exit Channel) {
	defer close(exit)
	ioutil.ReadAll(os.Stdin)
}

func connect(port int, path string) *websocket.Conn {
	address := fmt.Sprintf("127.0.0.1:%v", port)
	sconn, err := net.DialTimeout("tcp", address, Millis(2000))
	PanicIfError(err)
	url, err := url.Parse(fmt.Sprintf("ws://localhost:%v%v", port, path))
	PanicIfError(err)
	headers := http.Header{}
	wsconn, _, err := websocket.NewClient(sconn, url, headers, 1024, 1024)
	PanicIfError(err)
	return wsconn
}

func port() int {
	ep := endpoint()
	parts := strings.SplitAfterN(ep, ":", 2)
	port, err := strconv.Atoi(parts[1])
	PanicIfError(err)
	return port
}

func endpoint() string {
	ep := os.Getenv("LV_NBE_ENDPOINT")
	if len(strings.TrimSpace(ep)) > 0 {
		return ep
	}
	return ":0"
}
