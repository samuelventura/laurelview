package lvnbe

import (
	"fmt"
	"sync"
	"time"
)

type idState struct {
	count  uint
	prefix string
	mutex  sync.Mutex
}

type Id interface {
	Next() string
}

func NewId(prefix string) Id {
	s := &idState{}
	s.prefix = prefix
	return s
}

func (id *idState) Next() string {
	count := id.next()
	//1678 - January 1, 1970 UTC - 2262
	//last 3 digits are always zero
	now := time.Now().UnixNano() / 1000
	return fmt.Sprintf("%s-%d-%d", id.prefix, count, now)
}

func (id *idState) next() uint {
	defer id.mutex.Unlock()
	id.mutex.Lock()
	id.count++
	count := id.count
	return count
}
