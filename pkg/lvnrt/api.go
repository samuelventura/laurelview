package lvnrt

type ItemArgs struct {
	Host  string `json:"host"`
	Port  uint   `json:"port"`
	Slave uint   `json:"slave"`
}

type QueryArgs struct {
	Index    uint   `json:"index"`
	Request  string `json:"request"`
	Response string `json:"response"`
	Count    uint   `json:"count"`
	Total    uint   `json:"total"`
	Error    string `json:"error"`
}

//slave or bus address depending on mutation
//name being status-slave or status-bus
type StatusArgs struct {
	Address  string
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
