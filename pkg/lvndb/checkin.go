package lvndb

//sync to panic in entry client
func NewCheckin(rt Runtime) Dispatch {
	stateDispatch := rt.GetDispatch("state")
	log := rt.PrefixLog("checkin")
	return func(mut Mutation) {
		switch mut.Name {
		case "create":
			args := mut.Args.(Map)
			one := OneArgs{}
			name, err := ParseString(args, "name")
			PanicIfError(err)
			one.Name = name
			json, err := ParseString(args, "json")
			PanicIfError(err)
			one.Json = json
			nmut := mut
			nmut.Args = one
			stateDispatch(nmut)
		case "update":
			args := mut.Args.(Map)
			one := OneArgs{}
			id, err := ParseUint(args, "id")
			PanicIfError(err)
			one.Id = id
			name, err := ParseString(args, "name")
			PanicIfError(err)
			one.Name = name
			json, err := ParseString(args, "json")
			PanicIfError(err)
			one.Json = json
			nmut := mut
			nmut.Args = one
			stateDispatch(nmut)
		case "delete":
			id, err := CastUint(mut.Args, "id")
			PanicIfError(err)
			nmut := mut
			nmut.Args = id
			stateDispatch(nmut)
		case ":add":
			stateDispatch(mut)
		case ":remove":
			stateDispatch(mut)
		case ":dispose":
			stateDispatch(mut)
		default:
			log.Debug(mut)
		}
	}
}
