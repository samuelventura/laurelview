package lvnrt

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"strings"
	"time"
)

type busQueryDso struct {
	sid     string
	slave   uint
	request string
}

func slaveId(slave uint) string {
	ids := "123456789ABCDEFGHIJKLMNOPQRSTUV"
	if slave > 0 && slave < 32 {
		return ids[slave-1 : slave]
	}
	panicLN("invalid slave", slave)
	return "invalid"
}

func NewBus(rt Runtime) Dispatch {
	dispose := NopAction
	log := prefixLogger(rt.Log, "bus")
	dispatchs := make(map[string]Dispatch)
	dispatchs["dispose"] = func(mut *Mutation) {
		defer disposeArgs(mut.Args)
		defer dispose()
		clearDispatch(dispatchs)
	}
	dispatchs["bus"] = func(mut *Mutation) {
		delete(dispatchs, "bus")
		dialtoms := rt.Getv("bus.dialtoms").(int64)
		writetoms := rt.Getv("bus.writetoms").(int64)
		readtoms := rt.Getv("bus.readtoms").(int64)
		sleepms := rt.Getv("bus.sleepms").(int64)
		retryms := rt.Getv("bus.retryms").(int64)
		bus := mut.Args.(*BusArgs)
		address := fmt.Sprintf("%v:%v", bus.Host, bus.Port)
		log := prefixLogger(rt.Log, "bus", address)
		//size = 1 may be in reconnecting loop
		queue := make(chan *busQueryDso, 1)
		exit := make(Channel)
		busy := false
		status := func(query *busQueryDso, response string) {
			mut := &Mutation{}
			mut.Sid = query.sid
			mut.Name = "status"
			mut.Args = &StatusArgs{
				Slave:    fmt.Sprintf("%v:%v:%v", bus.Host, bus.Port, query.slave),
				Request:  query.request,
				Response: response,
			}
			rt.Post("self", mut)
		}
		next := func(request string) string {
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
			panicLN("invalid request", request)
			return "invalid"
		}
		command := func(request string, slave uint) string {
			id := slaveId(slave)
			switch request {
			case "read-value":
				return fmt.Sprintf("*%vB1", id)
			case "read-peak":
				return fmt.Sprintf("*%vB2", id)
			case "read-valley":
				return fmt.Sprintf("*%vB3", id)
			case "reset-peak":
				return fmt.Sprintf("*%vC3", id)
			case "reset-valley":
				return fmt.Sprintf("*%vC9", id)
			case "apply-tara":
				return fmt.Sprintf("*%vCA", id)
			case "reset-tara":
				return fmt.Sprintf("*%vCB", id)
			case "reset-cold":
				return fmt.Sprintf("*%vC0", id)
			}
			panicLN("invalid request", request)
			return "invalid"
		}
		dispose = func() {
			close(exit)
		}
		queries := list.New()
		slaves := make(map[uint]*list.Element)
		push := func(sid string, slave uint, request string) {
			query := &busQueryDso{}
			query.sid = sid
			query.request = request
			query.slave = slave
			slaves[slave] = queries.PushBack(query)
		}
		feed := func() {
			if busy {
				return
			}
			ls := len(slaves)
			lq := queries.Len()
			element := queries.Front()
			assertTrue(ls == lq, "slaves != queries", ls, lq, element)
			assertTrue(element != nil || ls == 0, "ls > 0 and nil element", ls, element)
			if element != nil {
				query := element.Value.(*busQueryDso)
				queries.Remove(element)
				request := next(query.request)
				push(query.sid, query.slave, request)
				busy = true
				queue <- query
			}
		}
		dispatchs["slave"] = func(mut *Mutation) {
			args := mut.Args.(*SlaveArgs)
			element, ok := slaves[args.Slave]
			if args.Count == 0 {
				assertTrue(ok, "slave not found", args.Slave)
				delete(slaves, args.Slave)
				queries.Remove(element)
			} else {
				//0->1 and 2->1 are valid transitions
				if !ok {
					assertTrue(args.Count == 1, "slave should be 1", args.Slave, args.Count)
					push(mut.Sid, args.Slave, "read-value")
				}
			}
			feed()
		}
		dispatchs["query"] = func(mut *Mutation) {
			args := mut.Args.(*QueryArgs)
			element, ok := slaves[args.Index]
			assertTrue(ok, "slave not found", args.Index)
			queries.Remove(element)
			push(mut.Sid, args.Index, args.Request)
			feed()
		}
		dispatchs["status"] = func(mut *Mutation) {
			assertTrue(busy, "not busy")
			busy = false
			rt.Post("hub", mut)
			feed()
		}
		read := func(conn net.Conn) bool {
			defer conn.Close()
			cr := byte(13)
			reader := bufio.NewReader(conn)
			for {
				select {
				case <-exit:
					return true
				case query := <-queue:
					cmd := command(query.request, query.slave)
					buf := []byte(cmd + "\r")
					_, err := reader.Discard(reader.Buffered())
					traceIfError(log.Trace, err)
					if err != nil {
						status(query, "error")
						return false
					}
					err = conn.SetWriteDeadline(future(writetoms))
					traceIfError(log.Trace, err)
					if err != nil {
						status(query, "error")
						return false
					}
					n, err := conn.Write(buf)
					m := len(buf)
					if err == nil && n != m {
						err = fmt.Errorf("wrote %v of %v", n, m)
					}
					traceIfError(log.Trace, err)
					if err != nil {
						status(query, "error")
						return false
					}
					res := "ok"
					if strings.HasPrefix(query.request, "read-") {
						err = conn.SetReadDeadline(future(readtoms))
						traceIfError(log.Trace, err)
						if err != nil {
							status(query, "error")
							return false
						}
						res, err = reader.ReadString(cr)
						traceIfError(log.Trace, err)
						if err != nil {
							status(query, "error")
							//do not close, may timeout after cold reset
							if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
								continue
							} else {
								return false
							}
						}
					}
					status(query, strings.TrimSpace(res))
				}
			}
		}
		loop := func() {
			defer traceRecover(log.Warn)
			managed := rt.Managed(address)
			defer managed()
			for {
				conn, err := net.DialTimeout("tcp", address, millis(dialtoms))
				traceIfError(log.Trace, err)
				if err != nil {
					to := future(retryms)
					for time.Now().Before(to) {
						select {
						case <-exit:
							return
						default:
							time.Sleep(millis(sleepms))
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
	return mapDispatch(log, dispatchs)
}
