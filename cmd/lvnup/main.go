package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/fasthttp/websocket"
)

func main() {
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	log := DefaultLog()
	defer CloseLog(log)

	logger := PrefixLogger(log)
	ep := os.Getenv("LV_NUP_ENDPOINT")
	logger.Info("LV_NUP_ENDPOINT", ep)
	go run(logger, ep)

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

func connect(ep string, path string) *websocket.Conn {
	sconn, err := net.DialTimeout("tcp", ep, Millis(2000))
	PanicIfError(err)
	url, err := url.Parse(fmt.Sprintf("ws://%v%v", ep, path))
	PanicIfError(err)
	headers := http.Header{}
	wsconn, _, err := websocket.NewClient(sconn, url, headers, 1024, 1024)
	PanicIfError(err)
	return wsconn
}

func endpoint() string {
	ep := os.Getenv("LV_NUP_ENDPOINT")
	if len(strings.TrimSpace(ep)) > 0 {
		return ep
	}
	return "127.0.0.1:31601"
}
