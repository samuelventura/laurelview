package lvsdk

func SendChannel(channel Channel, any Any) {
	channel <- any
}

//FIXME pattern not mature
func CloseChannel(channel Channel) {
	select {
	case <-channel:
	default:
		close(channel)
	}
}

//FIXME pattern not mature
func WaitChannel(channel Channel, output Output) {
	output("waiting channel...")
	<-channel
	output("waiting channel done")
}
