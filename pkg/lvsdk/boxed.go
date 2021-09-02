package lvsdk

type countDso struct {
	value uint
}

type Count interface {
	Get() uint
	Set(uint)
	Inc() uint
	Dec() uint
}

func NewCount() Count {
	return &countDso{}
}

func (c *countDso) Get() uint {
	return c.value
}

func (c *countDso) Set(value uint) {
	c.value = value
}

func (c *countDso) Inc() uint {
	c.value++
	return c.value
}

func (c *countDso) Dec() uint {
	c.value--
	return c.value
}

type flagDso struct {
	value bool
}

type Flag interface {
	Get() bool
	Set(bool)
}

func NewFlag() Flag {
	return &flagDso{}
}

func (f *flagDso) Get() bool {
	return f.value
}

func (f *flagDso) Set(value bool) {
	f.value = value
}
