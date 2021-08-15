package lvsdk

import (
	"fmt"
	"sync"
	"time"
)

type idDso struct {
	count  uint64
	prefix string
	mutex  *sync.Mutex
}

type Id interface {
	Next() string
}

func NewId(prefix string) Id {
	s := &idDso{}
	s.mutex = new(sync.Mutex)
	s.prefix = prefix
	return s
}

func (id *idDso) Next() string {
	count := id.next()
	//1678 - January 1, 1970 UTC - 2262
	//last 3 digits are always zero
	now := time.Now().UnixNano() / 1000
	return fmt.Sprintf("%s-%d-%d", id.prefix, count, now)
}

func (id *idDso) next() uint64 {
	defer id.mutex.Unlock()
	id.mutex.Lock()
	id.count++
	count := id.count
	return count
}
