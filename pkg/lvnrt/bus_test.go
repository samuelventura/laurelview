package lvnrt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRtBusBasic(t *testing.T) {
	to := newTestOutput()
	defer to.close()
	log := to.logger()
	echo := newTestEcho(log)
	log.Info("echo", "port", echo.port())
	defer echo.close()
	echo.ping()
	rt := NewRuntime(to.push)
	defer rt.ManagedWait()
	rt.Setv("bus.dialtoms", int64(400))
	rt.Setv("bus.writetoms", int64(400))
	rt.Setv("bus.readtoms", int64(400))
	rt.Setv("bus.sleepms", int64(10))
	rt.Setv("bus.retryms", int64(2000))
	rt.Setd("hub", to.dispatch("hub"))
	brt := rt.Overlay("bus")
	disp := asyncDispatch(log.Warn, NewBus(brt))
	brt.Setd("bus", disp)
	disp(&Mutation{Name: "setup", Sid: "tid", Args: &BusArgs{
		Host: "127.0.0.1",
		Port: echo.port(),
	}})
	for i := 1; i < 32; i++ {
		disp(&Mutation{Name: "slave", Sid: "tid", Args: &SlaveArgs{
			Slave: uint(i),
			Count: 1,
		}})
		//first one repeats (should only happen with echo test server)
		to.matchWait(t, 200, "trace", "echo", fmt.Sprintf(".%vB1.0D.", slaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		disp(&Mutation{Name: "query", Sid: "tid", Args: &QueryArgs{
			Index:   uint(i),
			Request: "reset-peak",
		}})
		to.matchWait(t, 200, "trace", "echo", fmt.Sprintf(".%vC3.0D.", slaveId(uint(i))))
		to.matchWait(t, 200, "trace", "echo", fmt.Sprintf(".%vB2.0D.", slaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		disp(&Mutation{Name: "query", Sid: "tid", Args: &QueryArgs{
			Index:   uint(i),
			Request: "reset-valley",
		}})
		to.matchWait(t, 200, "trace", "echo", fmt.Sprintf(".%vC9.0D.", slaveId(uint(i))))
		to.matchWait(t, 200, "trace", "echo", fmt.Sprintf(".%vB3.0D.", slaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		disp(&Mutation{Name: "query", Sid: "tid", Args: &QueryArgs{
			Index:   uint(i),
			Request: "apply-tara",
		}})
		to.matchWait(t, 200, "trace", "echo", fmt.Sprintf(".%vCA.0D.", slaveId(uint(i))))
		to.matchWait(t, 200, "trace", "echo", fmt.Sprintf(".%vB1.0D.", slaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		disp(&Mutation{Name: "query", Sid: "tid", Args: &QueryArgs{
			Index:   uint(i),
			Request: "reset-tara",
		}})
		to.matchWait(t, 200, "trace", "echo", fmt.Sprintf(".%vCB.0D.", slaveId(uint(i))))
		to.matchWait(t, 200, "trace", "echo", fmt.Sprintf(".%vB1.0D.", slaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		disp(&Mutation{Name: "query", Sid: "tid", Args: &QueryArgs{
			Index:   uint(i),
			Request: "reset-cold",
		}})
		to.matchWait(t, 200, "trace", "echo", fmt.Sprintf(".%vC0.0D.", slaveId(uint(i))))
		to.matchWait(t, 200, "trace", "echo", fmt.Sprintf(".%vB1.0D.", slaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		disp(&Mutation{Name: "slave", Sid: "tid", Args: &SlaveArgs{
			Slave: uint(i),
			Count: 2,
		}})
	}
	for i := 1; i < 32; i++ {
		disp(&Mutation{Name: "slave", Sid: "tid", Args: &SlaveArgs{
			Slave: uint(i),
			Count: 0,
		}})
	}
	disp(&Mutation{Name: "dispose", Sid: "tid"})
}

func TestRtSlaveId(t *testing.T) {
	assert.Equal(t, "1", slaveId(1))
	assert.Equal(t, "2", slaveId(2))
	assert.Equal(t, "3", slaveId(3))
	assert.Equal(t, "4", slaveId(4))
	assert.Equal(t, "5", slaveId(5))
	assert.Equal(t, "6", slaveId(6))
	assert.Equal(t, "7", slaveId(7))
	assert.Equal(t, "8", slaveId(8))
	assert.Equal(t, "9", slaveId(9))
	assert.Equal(t, "A", slaveId(10))
	assert.Equal(t, "B", slaveId(11))
	assert.Equal(t, "C", slaveId(12))
	assert.Equal(t, "D", slaveId(13))
	assert.Equal(t, "E", slaveId(14))
	assert.Equal(t, "F", slaveId(15))
	assert.Equal(t, "G", slaveId(16))
	assert.Equal(t, "H", slaveId(17))
	assert.Equal(t, "I", slaveId(18))
	assert.Equal(t, "J", slaveId(19))
	assert.Equal(t, "K", slaveId(20))
	assert.Equal(t, "L", slaveId(21))
	assert.Equal(t, "M", slaveId(22))
	assert.Equal(t, "N", slaveId(23))
	assert.Equal(t, "O", slaveId(24))
	assert.Equal(t, "P", slaveId(25))
	assert.Equal(t, "Q", slaveId(26))
	assert.Equal(t, "R", slaveId(27))
	assert.Equal(t, "S", slaveId(28))
	assert.Equal(t, "T", slaveId(29))
	assert.Equal(t, "U", slaveId(30))
	assert.Equal(t, "V", slaveId(31))
}
