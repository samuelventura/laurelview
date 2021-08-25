package lvsdk

func SendChannel(channel Channel, any Any) {
	channel <- any
}

func WaitChannel(channel Channel) Any {
	any := <-channel
	return any
}

func WaitClose(close func() Channel) Any {
	any := <-close()
	return any
}
