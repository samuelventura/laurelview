package lvnrt

//sync to panic in entry client
func NewCheckin(rt Runtime) Dispatch {
	state := rt.GetDispatch("state")
	log := rt.PrefixLog("checkin")
	return func(mut *Mutation) {
		switch mut.Name {
		case ":add":
			state(mut)
		case ":remove":
			state(mut)
		case "setup":
			args := mut.Args.([]Any)
			items := make([]*ItemArgs, 0, len(args))
			for _, mii := range args {
				mi := mii.(Map)
				item := &ItemArgs{}
				host, err := ParseString(mi, "host")
				PanicIfError(err)
				item.Host = host
				port, err := ParseUint(mi, "port")
				PanicIfError(err)
				item.Port = port
				slave, err := ParseUint(mi, "slave")
				PanicIfError(err)
				item.Slave = slave
				items = append(items, item)
			}
			nmut := *mut
			nmut.Args = items
			state(&nmut)
		case "query":
			args := mut.Args.(Map)
			query := &QueryArgs{}
			index, err := ParseUint(args, "index")
			PanicIfError(err)
			query.Index = index
			request, err := ParseString(args, "request")
			PanicIfError(err)
			query.Request = request
			response, err := MaybeString(args, "response", "")
			PanicIfError(err)
			query.Response = response
			errorm, err := MaybeString(args, "error", "")
			PanicIfError(err)
			query.Error = errorm
			count, err := MaybeUint(args, "count", 0)
			PanicIfError(err)
			query.Count = count
			total, err := MaybeUint(args, "total", 0)
			PanicIfError(err)
			query.Total = total
			nmut := *mut
			nmut.Args = query
			state(&nmut)
		default:
			log.Debug(mut)
		}
	}
}
