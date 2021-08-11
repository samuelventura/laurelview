package lvnrt

type countState struct {
	count uint
}

type Count interface {
	Count() uint
	Inc() uint
	Dec() uint
}

func NewCount() Count {
	s := &countState{}
	return s
}

func (c *countState) Count() uint {
	return c.count
}

func (c *countState) Inc() uint {
	c.count++
	return c.count
}

func (c *countState) Dec() uint {
	c.count--
	return c.count
}
