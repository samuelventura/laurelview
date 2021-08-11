package lvnrt

type Output = func(string, ...Any)
type Map = map[string]Any
type Queue = chan Action
type Channel = chan Any
type Any = interface{}
type Action = func()

type Dispatch = func(*Mutation)
type Factory = func() Dispatch

var NopAction = func() {}
var NopDispatch = func(*Mutation) {}

type Mutation struct {
	Sid  string
	Name string
	Args Any
}

type AddArgs struct {
	Callback Dispatch
}

type RemoveArgs struct {
}

type DisposeArgs struct {
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
}

type StatusArgs struct {
	Slave    string
	Request  string
	Response string
}

type ChannelArgs struct {
	Dispatch Dispatch
	Host     string
	Port     uint
}

type SlaveArgs struct {
	Slave uint
	Count uint
}
