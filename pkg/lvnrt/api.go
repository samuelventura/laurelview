package lvnrt

type AddArgs struct {
	Callback Dispatch
}

type SetupArgs struct {
	Items []*ItemArgs
}

type ItemArgs struct {
	Host  string
	Port  uint
	Slave uint
}

type QueryArgs struct {
	Index    uint
	Request  string
	Response string
	Count    uint
	Total    uint
	Error    string
}

type StatusArgs struct {
	Slave    string
	Request  string
	Response string
	Error    string
}

type BusArgs struct {
	Host string
	Port uint
}

type SlaveArgs struct {
	Slave uint
	Count uint
}
