package lvsdk

type countDso struct {
	count uint
}

type Count interface {
	Count() uint
	Inc() uint
	Dec() uint
}

func NewCount() Count {
	s := &countDso{}
	return s
}

func (c *countDso) Count() uint {
	return c.count
}

func (c *countDso) Inc() uint {
	c.count++
	return c.count
}

func (c *countDso) Dec() uint {
	c.count--
	return c.count
}
