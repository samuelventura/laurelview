package lvnrt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRtDpm(t *testing.T) {
	to := NewTestOutput()
	defer to.Close()
	log := to.Logger()
	dpm := NewDpm(log, ":0", 0)
	defer WaitClose(dpm.Close)
	address := fmt.Sprintf("127.0.0.1:%v", dpm.Port())
	socket := NewSocketDial(address, 400, 13)
	defer socket.Close()
	testDpmWriteLine(socket, busReadValue)
	assert.Equal(t, busReadValue, testDpmReadLine(socket))
	testDpmWriteLine(socket, busReadPeak)
	assert.Equal(t, busReadPeak, testDpmReadLine(socket))
	testDpmWriteLine(socket, busReadValley)
	assert.Equal(t, busReadValley, testDpmReadLine(socket))
	testDpmWriteLine(socket, busResetPeak)
	testDpmWriteLine(socket, busResetValley)
	testDpmWriteLine(socket, busApplyTara)
	testDpmWriteLine(socket, busResetTara)
	testDpmWriteLine(socket, busResetCold)
	testDpmWriteLine(socket, busReadValue)
	assert.Equal(t, busReadValue, testDpmReadLine(socket))
	testDpmWriteLine(socket, busReadValue)
	testDpmDiscard(socket)
	testDpmWriteLine(socket, busReadPeak)
	assert.Equal(t, busReadPeak, testDpmReadLine(socket))
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
