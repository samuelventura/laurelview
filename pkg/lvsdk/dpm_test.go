package lvsdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSdkDpm(t *testing.T) {
	to := NewTestOutput()
	defer to.Close()
	log := to.Logger()
	dpm := NewDpm(log, ":0", 0)
	address := fmt.Sprintf("127.0.0.1:%v", dpm.Port())
	socket := NewSocket(address, 400)
	defer socket.Close()
	testDpmWriteLine(socket, "*1B1")
	assert.Equal(t, "*1B1", testDpmReadLine(socket))
	testDpmWriteLine(socket, "*1B2")
	assert.Equal(t, "*1B2", testDpmReadLine(socket))
	testDpmWriteLine(socket, "*1B3")
	assert.Equal(t, "*1B3", testDpmReadLine(socket))
	testDpmWriteLine(socket, "*1C3")
	testDpmWriteLine(socket, "*1C9")
	testDpmWriteLine(socket, "*1CA")
	testDpmWriteLine(socket, "*1CB")
	testDpmWriteLine(socket, "*1C0")
	testDpmWriteLine(socket, "*1B1")
	assert.Equal(t, "*1B1", testDpmReadLine(socket))
	testDpmWriteLine(socket, "*1B1")
	testDpmDiscard(socket)
	testDpmWriteLine(socket, "*1B2")
	assert.Equal(t, "*1B2", testDpmReadLine(socket))
}

func testDpmReadLine(socket Socket) string {
	res, err := socket.ReadLine(400)
	PanicIfError(err)
	return res
}

func testDpmWriteLine(socket Socket, req string) {
	err := socket.WriteLine(req, 400)
	PanicIfError(err)
}

func testDpmDiscard(socket Socket) {
	err := socket.Discard(100)
	PanicIfError(err)
}
