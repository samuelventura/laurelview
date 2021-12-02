package main

import (
	"fmt"
	"sync"
	"time"
)

type idDso struct {
	count  uint
	prefix string
	mutex  *sync.Mutex
}

func newId(prefix string) *idDso {
	dso := &idDso{}
	dso.mutex = new(sync.Mutex)
	dso.prefix = prefix
	return dso
}

func (dso *idDso) next() string {
	count := dso.inc()
	//1678 - January 1, 1970 UTC - 2262
	//last 3 digits are always zero
	now := time.Now()
	when := now.Format("20060102T150405.000")
	return fmt.Sprintf("%s-%s-%d", dso.prefix, when, count)
}

func (dso *idDso) inc() uint {
	defer dso.mutex.Unlock()
	dso.mutex.Lock()
	dso.count++
	count := dso.count
	return count
}
