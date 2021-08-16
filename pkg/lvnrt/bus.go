package lvnrt

import (
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
	PanicLN("invalid slave", slave)
	return "invalid"
}

func NewBus(rt Runtime) Dispatch {
	dispose := NopAction
	log := PrefixLogger(rt.Log, "bus")
	cleaner := rt.Getc("bus")
	dispatchs := make(map[string]Dispatch)
	dispatchs["dispose"] = func(mut *Mutation) {
		defer DisposeArgs(mut.Args)
		defer dispose()
		ClearDispatch(dispatchs)
	}
	dispatchs["setup"] = func(mut *Mutation) {
		delete(dispatchs, "bus")
		dialtoms := rt.Getv("bus.dialtoms").(int64)
		writetoms := rt.Getv("bus.writetoms").(int64)
		readtoms := rt.Getv("bus.readtoms").(int64)
		sleepms := rt.Getv("bus.sleepms").(int64)
		retryms := rt.Getv("bus.retryms").(int64)
		resetms := rt.Getv("bus.resetms").(int64)
		discardms := rt.Getv("bus.discardms").(int64)
		bus := mut.Args.(*BusArgs)
		address := fmt.Sprintf("%v:%v", bus.Host, bus.Port)
		log := PrefixLogger(rt.Log, "bus", address)
		//size = 1 may be in reconnecting loop
		queue := make(chan *busQueryDso, 1)
		exit := make(Channel)
		busy := false
		status := func(query *busQueryDso, response string, err error) {
			mut := &Mutation{}
			mut.Sid = query.sid
			mut.Name = "status"
			mut.Args = &StatusArgs{
				Slave:    fmt.Sprintf("%v:%v:%v", bus.Host, bus.Port, query.slave),
				Request:  query.request,
				Response: response, //+ fmt.Sprint(err),
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
			PanicLN("invalid request", request)
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
			PanicLN("invalid request", request)
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
			AssertTrue(ls == lq, "slaves != queries", ls, lq, element)
			AssertTrue(element != nil || ls == 0, "ls > 0 and nil element", ls, element)
			if element != nil {
				//FIXME 20210814T013131.279 warn bus 127.0.0.1:54496 recover interface conversion: interface {} is nil, not *lvnrt.busQueryDso goroutine 39 [running]:
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
			//all transitions are valid
			if args.Count == 0 {
				delete(slaves, args.Slave)
				queries.Remove(element)
			} else {
				if !ok {
					push(mut.Sid, args.Slave, "read-value")
				}
			}
			feed()
		}
		dispatchs["query"] = func(mut *Mutation) {
			args := mut.Args.(*QueryArgs)
			element, ok := slaves[args.Index]
			AssertTrue(ok, "slave not found", args.Index)
			queries.Remove(element)
			push(mut.Sid, args.Index, args.Request)
			feed()
		}
		dispatchs["status"] = func(mut *Mutation) {
			AssertTrue(busy, "not busy")
			busy = false
			rt.Post("hub", mut)
			feed()
		}
		read := func(conn net.Conn) bool {
			cleaner.AddCloser(address, conn)
			defer cleaner.Remove(address)
			socket := NewSocketConn(conn)
			defer socket.Close()
			for {
				select {
				case <-exit:
					return true
				case query := <-queue:
					cmd := command(query.request, query.slave)
					err := socket.Discard(discardms)
					TraceIfError(log.Trace, err)
					if err != nil {
						status(query, "error", err)
						return false
					}
					//log.Info("REQUEST >", cmd)
					err = socket.WriteLine(cmd, writetoms)
					TraceIfError(log.Trace, err)
					if err != nil {
						status(query, "error", err)
						return false
					}
					res := "ok"
					if strings.HasPrefix(query.request, "read-") {
						res, err = socket.ReadLine(readtoms)
						TraceIfError(log.Trace, err)
						if err != nil {
							status(query, "error", err)
							//do not close, may timeout after cold reset
							nerr, ok := err.(net.Error)
							if ok && nerr.Timeout() {
								//log.Info("TIMEOUT <", cmd)
								continue
							} else {
								return false
							}
						}
					} else {
						//bus get unresponsive after resets 400ms works
						time.Sleep(Millis(resetms))
					}
					//log.Info("RESPONSE <", strings.TrimSpace(res))
					status(query, strings.TrimSpace(res), nil)
				}
			}
		}
		loop := func() {
			defer TraceRecover(log.Debug)
			for {
				conn, err := net.DialTimeout("tcp", address, Millis(dialtoms))
				TraceIfError(log.Trace, err)
				if err != nil {
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
