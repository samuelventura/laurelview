package lvnrt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRtBusDpm(t *testing.T) {
	to := NewTestOutput()
	defer to.Close()
	log := to.Logger()
	dpm := NewDpm(log, ":0", 0)
	log.Info("dpm", "port", dpm.Port())
	defer dpm.Close()
	dpm.Echo()
	rt := NewRuntime(to.Log)
	defer rt.Close()
	rt.SetValue("bus.dialtoms", 400)
	rt.SetValue("bus.writetoms", 400)
	rt.SetValue("bus.readtoms", 400)
	rt.SetValue("bus.discardms", 0)
	rt.SetValue("bus.sleepms", 10)
	rt.SetValue("bus.retryms", 2000)
	rt.SetValue("bus.resetms", 0)
	rt.SetCleaner("bus", NewCleaner(rt.PrefixLog("bus", "clean")))
	rt.SetDispatch("hub", to.Dispatch("hub"))
	rt.SetDispatch("bus", AsyncDispatch(log.Debug, NewBus(rt)))
	busDispatch := rt.GetDispatch("bus")
	busDispatch(M("setup", "tid", &BusArgs{
		Host: "127.0.0.1",
		Port: dpm.Port(),
	}))
	for i := 1; i < 32; i++ {
		busDispatch(M("slave", "tid", &SlaveArgs{
			Slave: uint(i),
			Count: 1,
		}))
		//first one repeats (should only happen with echo test server)
		to.MatchWait(t, 200, "trace", "dpm", "true", fmt.Sprintf(".%vB1.0D.", busSlaveId(uint(i))))
		to.MatchWait(t, 200, "trace", "hub", fmt.Sprintf("read-value .%vB1", busSlaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		busDispatch(M("query", "tid", &QueryArgs{
			Index:   uint(i),
			Request: "reset-peak",
		}))
		//lack of order warranty prevents from testing hub as well
		to.MatchWait(t, 200, "trace", "dpm", "false", fmt.Sprintf(".%vC3.0D.", busSlaveId(uint(i))))
		to.MatchWait(t, 200, "trace", "dpm", "true", fmt.Sprintf(".%vB2.0D.", busSlaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		busDispatch(M("query", "tid", &QueryArgs{
			Index:   uint(i),
			Request: "reset-valley",
		}))
		to.MatchWait(t, 200, "trace", "dpm", "false", fmt.Sprintf(".%vC9.0D.", busSlaveId(uint(i))))
		to.MatchWait(t, 200, "trace", "dpm", "true", fmt.Sprintf(".%vB3.0D.", busSlaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		busDispatch(M("query", "tid", &QueryArgs{
			Index:   uint(i),
			Request: "apply-tara",
		}))
		to.MatchWait(t, 200, "trace", "dpm", "false", fmt.Sprintf(".%vCA.0D.", busSlaveId(uint(i))))
		to.MatchWait(t, 200, "trace", "dpm", "true", fmt.Sprintf(".%vB1.0D.", busSlaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		busDispatch(M("query", "tid", &QueryArgs{
			Index:   uint(i),
			Request: "reset-tara",
		}))
		to.MatchWait(t, 200, "trace", "dpm", "false", fmt.Sprintf(".%vCB.0D.", busSlaveId(uint(i))))
		to.MatchWait(t, 200, "trace", "dpm", "true", fmt.Sprintf(".%vB1.0D.", busSlaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		busDispatch(M("query", "tid", &QueryArgs{
			Index:   uint(i),
			Request: "reset-cold",
		}))
		to.MatchWait(t, 200, "trace", "dpm", "false", fmt.Sprintf(".%vC0.0D.", busSlaveId(uint(i))))
		to.MatchWait(t, 200, "trace", "dpm", "true", fmt.Sprintf(".%vB1.0D.", busSlaveId(uint(i))))
	}
	for i := 1; i < 32; i++ {
		busDispatch(M("slave", "tid", &SlaveArgs{
			Slave: uint(i),
			Count: 2,
		}))
	}
	for i := 1; i < 32; i++ {
		busDispatch(M("slave", "tid", &SlaveArgs{
			Slave: uint(i),
			Count: 0,
		}))
	}
	busDispatch(Mns("dispose", "tid"))
}

func TestRtSlaveId(t *testing.T) {
	assert.Equal(t, "1", busSlaveId(1))
	assert.Equal(t, "2", busSlaveId(2))
	assert.Equal(t, "3", busSlaveId(3))
	assert.Equal(t, "4", busSlaveId(4))
	assert.Equal(t, "5", busSlaveId(5))
	assert.Equal(t, "6", busSlaveId(6))
	assert.Equal(t, "7", busSlaveId(7))
	assert.Equal(t, "8", busSlaveId(8))
	assert.Equal(t, "9", busSlaveId(9))
	assert.Equal(t, "A", busSlaveId(10))
	assert.Equal(t, "B", busSlaveId(11))
	assert.Equal(t, "C", busSlaveId(12))
	assert.Equal(t, "D", busSlaveId(13))
	assert.Equal(t, "E", busSlaveId(14))
	assert.Equal(t, "F", busSlaveId(15))
	assert.Equal(t, "G", busSlaveId(16))
	assert.Equal(t, "H", busSlaveId(17))
	assert.Equal(t, "I", busSlaveId(18))
	assert.Equal(t, "J", busSlaveId(19))
	assert.Equal(t, "K", busSlaveId(20))
	assert.Equal(t, "L", busSlaveId(21))
	assert.Equal(t, "M", busSlaveId(22))
	assert.Equal(t, "N", busSlaveId(23))
	assert.Equal(t, "O", busSlaveId(24))
	assert.Equal(t, "P", busSlaveId(25))
	assert.Equal(t, "Q", busSlaveId(26))
	assert.Equal(t, "R", busSlaveId(27))
	assert.Equal(t, "S", busSlaveId(28))
	assert.Equal(t, "T", busSlaveId(29))
	assert.Equal(t, "U", busSlaveId(30))
	assert.Equal(t, "V", busSlaveId(31))
}
