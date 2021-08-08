package lvnbe

type Output = func(string, ...Any)
type Map = map[string]Any
type Queue = chan Action
type Any = interface{}
type Channel = chan Any
type Action = func()

type Callback = func(*Mutation)

var NopCallback = func(*Mutation) {}

type Mutation struct {
	Sid  string
	Name string
	Args Any
}

type AllArgs struct {
	Items []*OneArgs
}

type OneArgs struct {
	Id   uint
	Name string
	Json string
}

type CreateArgs struct {
	Id   uint
	Name string
	Json string
}

type DeleteArgs struct {
	Id uint
}

type UpdateArgs struct {
	Id   uint
	Name string
	Json string
}
