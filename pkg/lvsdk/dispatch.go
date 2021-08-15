package lvsdk

func ClearDispatch(dispatchs map[string]Dispatch) {
	for name := range dispatchs {
		delete(dispatchs, name)
	}
}

func MapDispatch(log Logger, dispmap map[string]Dispatch) Dispatch {
	return func(mut *Mutation) {
		dispatch, ok := dispmap[mut.Name]
		if ok {
			log.Trace(mut)
			dispatch(mut)
		} else {
			log.Debug(mut)
		}
	}
}

func AsyncDispatch(output Output, dispatch Dispatch) Dispatch {
	queue := make(chan *Mutation)
	loop := func() {
		defer TraceRecover(output)
		for mut := range queue {
			dispatch(mut)
		}
	}
	go loop()
	return func(mut *Mutation) {
		//do not close queue nor state dispose
		//let map dispatch report the ignore
		queue <- mut
	}
}
