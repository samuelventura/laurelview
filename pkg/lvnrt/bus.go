package lvnrt

import (
	"container/list"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type busQueryDso struct {
	sid     string
	slave   uint
	request string
}

/*
- A 400ms delay is needed even for no-response commands
because bus gets unresponsive after any of the resets
- A sync read takes 80s to detect a connection drop, this
was tested from mbair to nucmeg disconnecting wifi. It got down
to 13s adding the following setup:
	tcp := conn.(*net.TCPConn)
	tcp.SetNoDelay(true)
	tcp.SetLinger(0)
	tcp.SetKeepAlive(true)
	tcp.SetKeepAlivePeriod(millis(1000))
	tcp.SetWriteBuffer(0)
- A long 20s dial timeout is needed to ensure connection
success on first attempt for Laurel ethernet modules
- Disconnect on first read timeout is needed to avoid getting
in a read loop where reads timeout before EOF is ever detected.
- Even if Laurels nodes are enforced with slave 1, the user
still may not add slave 1 to the list of items making it
impossible to detect the timeout was on the node master
as to be smart on where to apply the disconnect on first read
timeout policy.
*/
func NewBus(rt Runtime) Dispatch {
	dispose := NopAction
	hubDispatch := rt.GetDispatch("hub")
	log := PrefixLogger(rt.Log, "bus")
	cleaner := NewCleaner(PrefixLogger(rt.Log, "bus", "cleaner"))
	dispatchs := make(map[string]Dispatch)
	dispatchs[":dispose"] = func(mut Mutation) {
		defer DisposeArgs(mut.Args)
		defer dispose()
		ClearDispatch(dispatchs)
	}
	dispatchs["setup"] = func(mut Mutation) {
		dialtoms := rt.GetValue("bus.dialtoms").(int)
		writetoms := rt.GetValue("bus.writetoms").(int)
		readtoms := rt.GetValue("bus.readtoms").(int)
		sleepms := rt.GetValue("bus.sleepms").(int)
		retryms := rt.GetValue("bus.retryms").(int)
		resetms := rt.GetValue("bus.resetms").(int)
		discardms := rt.GetValue("bus.discardms").(int)
		bus := mut.Args.(BusArgs)
		address := fmt.Sprintf("%v:%v", bus.Host, bus.Port)
		log := PrefixLogger(rt.Log, "bus", address)
		exit := make(Channel)
		cleaner.AddChannel("exit", exit)
		dispose = func() {
			cleaner.Close()
		}
		queries := list.New()
		slaves := make(map[uint]*list.Element)
		front := func() *list.Element {
			ls := len(slaves)
			lq := queries.Len()
			element := queries.Front()
			AssertTrue(ls == lq, "slaves != queries", ls, lq, element)
			AssertTrue(element != nil || ls == 0, "ls > 0 and nil element", ls, element)
			return element
		}
		push := func(sid string, slave uint, request string) {
			element, ok := slaves[slave]
			if ok {
				delete(slaves, slave)
				queries.Remove(element)
			}
			query := &busQueryDso{}
			query.sid = sid
			query.request = request
			query.slave = slave
			slaves[slave] = queries.PushBack(query)
		}
		remove := func(slave uint) {
			queries.Remove(slaves[slave])
			delete(slaves, slave)
		}
		var mutex sync.Mutex
		pop := func() *busQueryDso {
			mutex.Lock()
			defer mutex.Unlock()
			element := front()
			if element != nil {
				query := element.Value.(*busQueryDso)
				request := busNextRequest(query.request)
				push(query.sid, query.slave, request)
				return query
			}
			return nil
		}
		dispatchs["slave"] = func(mut Mutation) {
			mutex.Lock()
			defer mutex.Unlock()
			args := mut.Args.(SlaveArgs)
			_, ok := slaves[args.Slave]
			if args.Count == 0 && ok {
				remove(args.Slave)
			}
			if args.Count > 0 && !ok {
				push(mut.Sid, args.Slave, "read-value")
			}
		}
		dispatchs["query"] = func(mut Mutation) {
			mutex.Lock()
			defer mutex.Unlock()
			args := mut.Args.(QueryArgs)
			_, ok := slaves[args.Index]
			AssertTrue(ok, "slave not found", args.Index)
			push(mut.Sid, args.Index, args.Request)
		}
		status_slave := func(query *busQueryDso, response string, err error) {
			mut := Mutation{}
			mut.Sid = query.sid
			mut.Name = "status-slave"
			mut.Args = StatusArgs{
				Address:  fmt.Sprintf("%v:%v:%v", bus.Host, bus.Port, query.slave),
				Request:  query.request,
				Response: response,
				Error:    ErrorString(err),
			}
			hubDispatch(mut)
		}
		status_buserr := func(err error) {
			mut := Mutation{}
			mut.Name = "status-bus"
			mut.Args = StatusArgs{
				Address:  fmt.Sprintf("%v:%v", bus.Host, bus.Port),
				Request:  "Dial",
				Response: "error",
				Error:    ErrorString(err),
			}
			hubDispatch(mut)
		}
		read := func(conn net.Conn) bool {
			defer cleaner.Remove(address)
			cleaner.AddCloser(address, conn)
			socket := NewSocketConn(conn, 13)
			defer socket.Close()
			for {
				select {
				case <-exit:
					return true
				default:
					query := pop()
					if query == nil {
						time.Sleep(Millis(sleepms))
						continue
					}
					cmd := busRequestCode(query.request, query.slave)
					err := socket.Discard(discardms)
					TraceIfError(log.Trace, err)
					if err != nil {
						status_slave(query, "error1", err)
						return false
					}
					err = socket.WriteLine(cmd, writetoms)
					TraceIfError(log.Trace, err)
					if err != nil {
						status_slave(query, "error2", err)
						return false
					}
					res := "ok"
					if strings.HasPrefix(query.request, "read-") {
						res, err = socket.ReadLine(readtoms)
						TraceIfError(log.Trace, err)
						if err != nil {
							status_slave(query, "error3", err)
							return false
						}
					} else {
						//bus get unresponsive after resets 400ms works
						time.Sleep(Millis(resetms))
					}
					status_slave(query, strings.TrimSpace(res), nil)
				}
			}
		}
		loop := func() {
			defer TraceRecover(log.Debug)
			for {
				conn, err := net.DialTimeout("tcp", address, Millis(dialtoms))
				TraceIfError(log.Trace, err)
				if err != nil {
					status_buserr(err)
					to := Future(retryms)
					for time.Now().Before(to) {
						select {
						case <-exit:
							return
						default:
							time.Sleep(Millis(sleepms))
						}
					}
					continue
				}
				if read(conn) {
					return
				}
			}
		}
		go loop()
	}
	return MapDispatch(log, dispatchs)
}

func busSlaveId(slave uint) string {
	ids := "123456789ABCDEFGHIJKLMNOPQRSTUV"
	if slave > 0 && slave < 32 {
		return ids[slave-1 : slave]
	}
	PanicLN("invalid slave", slave)
	return "invalid"
}

const busReadValue = "*#B1"
const busReadPeak = "*#B2"
const busReadValley = "*#B3"
const busResetPeak = "*#C3"
const busResetValley = "*#C9"
const busApplyTara = "*#CA"
const busResetTara = "*#CB"
const busResetCold = "*#C0"

func busRequestCode(request string, slave uint) string {
	id := busSlaveId(slave)
	switch request {
	case "read-value":
		return strings.Replace(busReadValue, "#", id, 1)
	case "read-peak":
		return strings.Replace(busReadPeak, "#", id, 1)
	case "read-valley":
		return strings.Replace(busReadValley, "#", id, 1)
	case "reset-peak":
		return strings.Replace(busResetPeak, "#", id, 1)
	case "reset-valley":
		return strings.Replace(busResetValley, "#", id, 1)
	case "apply-tara":
		return strings.Replace(busApplyTara, "#", id, 1)
	case "reset-tara":
		return strings.Replace(busResetTara, "#", id, 1)
	case "reset-cold":
		return strings.Replace(busResetCold, "#", id, 1)
	}
	PanicLN("invalid request", request)
	return "invalid"
}

func busNextRequest(request string) string {
	switch request {
	case "read-value":
		return "read-value"
	case "read-peak":
		return "read-peak"
	case "read-valley":
		return "read-valley"
	case "reset-peak":
		return "read-peak"
	case "reset-valley":
		return "read-valley"
	case "apply-tara":
		return "read-value"
	case "reset-tara":
		return "read-value"
	case "reset-cold":
		return "read-value"
	}
	PanicLN("invalid request", request)
	return "invalid"
}
