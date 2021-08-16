package lvsdk

import (
	"testing"
)

func TestSdkCleaner(t *testing.T) {
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
	to.MatchWait(t, 200, "trace", "count", "0", "0")
	// immediate removal of anything arriving after close
	cleaner.AddAction("id3", func() {
		log.Trace("action3")
	})
	to.MatchWait(t, 200, "trace", "add", "id3")
	to.MatchWait(t, 200, "trace", "close", "id3")
	to.MatchWait(t, 200, "trace", "action3")
}
