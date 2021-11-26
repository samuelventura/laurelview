package main

import (
	"sync"

	"github.com/samuelventura/go-tree"
)

type singleDso struct {
	mutex  sync.Mutex
	single tree.Node
}

func newSingle() *singleDso {
	return &singleDso{}
}

func (dso *singleDso) enter(node tree.Node) {
	dso.mutex.Lock()
	defer dso.mutex.Unlock()
	if dso.single != nil {
		dso.single.Close()
		dso.single = nil
	}
	dso.single = node
}

func (dso *singleDso) exit(node tree.Node) {
	dso.mutex.Lock()
	defer dso.mutex.Unlock()
	if dso.single == node {
		dso.single.Close()
		dso.single = nil
	}
}
