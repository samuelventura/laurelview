package lvsdk

import (
	"container/list"
	"io"
)

type Cleaner interface {
	AddAction(id string, action Action)
	AddCloser(id string, closer io.Closer)
	AddChannel(id string, channel Channel)
	Remove(id string)
	Close()
}

type cleanerDso struct {
	log    Logger
	queue  Queue
	closed bool
	order  *list.List
	items  map[string]*list.Element
}

func NewCleaner(log Logger) Cleaner {
	c := &cleanerDso{}
	c.log = log
	c.queue = make(Queue)
	c.order = list.New()
	c.items = make(map[string]*list.Element)
	go c.loop()
	return c
}

func (c *cleanerDso) AddAction(id string, action Action) {
	c.queue <- func() {
		c.override(id)
		c.add(id, func() {
			c.log.Trace("close", id)
			delete(c.items, id)
			action()
		})
	}
}

func (c *cleanerDso) AddCloser(id string, closer io.Closer) {
	c.AddAction(id, func() {
		closer.Close()
	})
}

func (c *cleanerDso) AddChannel(id string, channel Channel) {
	c.AddAction(id, func() {
		close(channel)
	})
}

func (c *cleanerDso) Remove(id string) {
	c.queue <- func() {
		item, ok := c.items[id]
		if ok {
			c.safe(item)
		} else {
			c.log.Debug("nf404", id)
		}
	}
}

func (c *cleanerDso) Close() {
	c.queue <- func() {
		c.closed = true
		item := c.order.Front()
		for item != nil {
			next := item.Next()
			c.safe(item)
			item = next
		}
		c.log.Trace("count", len(c.items), c.order.Len())
	}
}

func (c *cleanerDso) loop() {
	for action := range c.queue {
		action()
	}
}

func (c *cleanerDso) override(id string) {
	item, ok := c.items[id]
	if ok {
		c.log.Debug("override", id)
		c.safe(item)
	}
}

func (c *cleanerDso) safe(item *list.Element) {
	TraceRecover(c.log.Debug)
	action := item.Value.(Action)
	c.order.Remove(item)
	action()
}

func (c *cleanerDso) add(id string, action Action) {
	c.log.Trace("add", id)
	item := c.order.PushBack(action)
	c.items[id] = item
	if c.closed {
		c.log.Trace("close", id)
		c.safe(item)
	}
}
