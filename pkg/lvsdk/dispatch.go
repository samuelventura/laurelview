package lvsdk

func ClearDispatch(dispatchs map[string]Dispatch) {
	for name := range dispatchs {
		delete(dispatchs, name)
	}
}

func MapDispatch(log Logger, dispmap map[string]Dispatch) Dispatch {
	return func(mut Mutation) {
		//recover by default to make AssertTrue the input check standard
		//:remove may be duplicated
		//state.query may arrive after :remove
		defer TraceRecoverMut(log.Debug, mut)
		dispatch, ok := dispmap[mut.Name]
		if ok {
			log.Trace(mut)
			dispatch(mut)
		} else {
			log.Debug(mut)
		}
	}
}

//like cleaner should never close anything and aim for idempotency
//it needs to be pair with state immutability for real efficacy
func AsyncDispatch(log Logger, dispatch Dispatch) Dispatch {
	queue := make(chan Mutation)
	catch := func(mut Mutation) {
		r := recover()
		if r != nil {
			log.Error("recover", mut, r)
		}
	}
	safe := func(disp Dispatch, mut Mutation) {
		defer catch(mut)
		disp(mut)
	}
	loop := func() {
		for mut := range queue {
			safe(dispatch, mut)
		}
	}
	go loop()
	return func(mut Mutation) {
		queue <- mut
	}
}

func DisposeArgs(arg Any) {
	action, ok := arg.(Action)
	if ok {
		action()
	}
	channel, ok := arg.(Channel)
	if ok {
		close(channel)
	}
}
