package lvsdk

import (
	"fmt"
	"testing"
)

func TestSdkCleanerBasic(t *testing.T) {
	to := NewTestOutput()
	defer to.Close()
	log := to.Logger()
	cleaner := NewCleaner(log)
	cleaner.AddAction("id0", func() {
		log.Trace("actionN")
	})
	to.MatchWait(t, 200, "trace", "add", "id0")
	cleaner.AddAction("id0", func() {
		log.Trace("action0")
	})
	to.MatchWait(t, 200, "debug", "override", "id0")
	to.MatchWait(t, 200, "trace", "close", "id0")
	to.MatchWait(t, 200, "trace", "actionN")
	to.MatchWait(t, 200, "trace", "add", "id0")
	cleaner.AddAction("id1", func() {
		log.Trace("action1")
	})
	to.MatchWait(t, 200, "trace", "add", "id1")
	cleaner.AddAction("id2", func() {
		log.Trace("action2")
	})
	to.MatchWait(t, 200, "trace", "add", "id2")
	cleaner.Remove("id0")
	to.MatchWait(t, 200, "trace", "close", "id0")
	to.MatchWait(t, 200, "trace", "action0")
	cleaner.Remove("id0")
	to.MatchWait(t, 200, "debug", "nf404", "id0")
	cleaner.Close()
	to.MatchWait(t, 200, "trace", "close", "id1")
	to.MatchWait(t, 200, "trace", "action1")
	to.MatchWait(t, 200, "trace", "close", "id2")
	to.MatchWait(t, 200, "trace", "action2")
	cleaner.Status(func(any Any) {
		c := any.(*cleanerDso)
		log.Trace("status", c.closed, len(c.items), c.order.Len())
	})
	to.MatchWait(t, 200, "trace", "status", "true", "0", "0")
	// immediate removal of anything arriving after close
	cleaner.AddAction("id3", func() {
		log.Trace("action3")
	})
	to.MatchWait(t, 200, "trace", "add", "id3")
	to.MatchWait(t, 200, "trace", "close", "id3")
	to.MatchWait(t, 200, "trace", "action3")
	cleaner.Status(func(any Any) {
		c := any.(*cleanerDso)
		log.Trace("status", c.closed, len(c.items), c.order.Len())
	})
	to.MatchWait(t, 200, "trace", "status", "true", "0", "0")
}

func TestSdkCleanerLoad(t *testing.T) {
	to := NewTestOutput()
	defer to.Close()
	cleaner := NewCleaner(NopLogger())
	dones := make(map[int]Channel)
	for i := 0; i < 100; i++ {
		ii := i
		done := make(Channel)
		dones[ii] = done
		go func() {
			defer close(done)
			for j := 0; j < 1000; j++ {
				id := fmt.Sprintf("id_%v_%v", ii, j)
				cleaner.AddAction(id, func() {})
			}
		}()
	}
	for _, done := range dones {
		goit := false
		for !goit {
			select {
			case <-done:
				goit = true
			default:
				cleaner.Status(func(any Any) {
					c := any.(*cleanerDso)
					AssertTrue(len(c.items) == c.order.Len())
				})
				continue
			}
		}
	}
	cleaner.Close()
	log := to.Logger()
	cleaner.Status(func(any Any) {
		c := any.(*cleanerDso)
		log.Trace("status", c.closed, len(c.items), c.order.Len())
	})
	to.MatchWait(t, 200, "trace", "status", "true", "0", "0")
	//after close
	dones = make(map[int]Channel)
	for i := 0; i < 100; i++ {
		ii := i
		done := make(Channel)
		dones[ii] = done
		go func() {
			defer close(done)
			for j := 0; j < 1000; j++ {
				id := fmt.Sprintf("id_%v_%v", ii, j)
				cleaner.AddAction(id, func() {})
			}
		}()
	}
	for _, done := range dones {
		<-done
	}
	cleaner.Status(func(any Any) {
		c := any.(*cleanerDso)
		log.Trace("status", c.closed, len(c.items), c.order.Len())
	})
	to.MatchWait(t, 200, "trace", "status", "true", "0", "0")
}
