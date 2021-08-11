package lvnrt

type channelDso struct {
}

func NewChannel() Dispatch {
	channel := &channelDso{}
	return channel.apply
}

func (channel *channelDso) apply(mut *Mutation) {

}
