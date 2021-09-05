package main

import (
	"encoding/json"
	"time"
)

func cycle(log Logger, ep string) {
	defer TraceRecover(log.Debug)
	log.Trace("connecting to /ws/db...", ep)
	conn := connect(ep, "/ws/db")
	log.Trace("connected to /ws/db", ep)
	defer conn.Close()
	cleandb := NewCleaner(log)
	defer cleandb.Close()
	onem := make(map[uint]OneArgs)
	for {
		mut := ReadMutation(conn)
		log.Trace("db.mut", mut)
		switch mut.Name {
		case "all":
			its := mut.Args.([]Any)
			for _, it := range its {
				itm := it.(Map)
				one := parse(itm)
				onem[one.Id] = one
			}
		case "create":
			itm := mut.Args.(Map)
			one := parse(itm)
			onem[one.Id] = one
		case "update":
			itm := mut.Args.(Map)
			one := parse(itm)
			onem[one.Id] = one
		case "delete":
			id := uint(mut.Args.(float64))
			delete(onem, id)
		case "ping": //keepalive
		default: //zero on error
			return
		}
		items := make([]ItemArgs, 0, len(onem))
		ones := make([]OneArgs, 0, len(onem))
		for _, it := range onem {
			var any Any
			err := json.Unmarshal([]byte(it.Json), &any)
			PanicIfError(err)
			jm := any.(Map)
			ia := ItemArgs{}
			host, err := ParseString(jm, "host")
			PanicIfError(err)
			ia.Host = host
			port, err := ParseUint(jm, "port")
			PanicIfError(err)
			ia.Port = port
			slave, err := ParseUint(jm, "slave")
			PanicIfError(err)
			ia.Slave = slave
			items = append(items, ia)
			ones = append(ones, it)
		}
		cleanrt := NewCleaner(log)
		cleandb.AddCleaner("clean", cleanrt)
		exit := make(Channel)
		cleanrt.AddChannel("exit", exit)
		cycle := func() {
			defer TraceRecover(log.Debug)
			log.Trace("connecting to /ws/rt...", ep)
			conn := connect(ep, "/ws/rt")
			log.Trace("connected to /ws/rt", ep)
			defer conn.Close()
			cleanrt.AddCloser("conn.rt", conn)
			setup := Mna("setup", items)
			log.Trace("rt.setup", setup)
			WriteMutation(conn, setup)
			for {
				mut := ReadMutation(conn)
				log.Trace("rt.mut", mut)
				switch mut.Name {
				case "query":
				default: //zero on error
					return
				}
			}
		}
		loop := func() {
			for {
				select {
				case <-exit:
					return
				default:
					cycle()
				}
			}
		}
		go loop()
	}
}

func parse(itm Map) OneArgs {
	one := OneArgs{}
	id, err := ParseUint(itm, "id")
	PanicIfError(err)
	one.Id = id
	name, err := ParseString(itm, "name")
	PanicIfError(err)
	one.Name = name
	json, err := ParseString(itm, "json")
	PanicIfError(err)
	one.Json = json
	return one
}

func run(log Logger, ep string) {
	for {
		cycle(log, ep)
		time.Sleep(Millis(1000))
	}
}
