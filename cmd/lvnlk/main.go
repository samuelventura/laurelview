package main

import (
	"log"
	"os"

	"github.com/samuelventura/go-tools"
	"github.com/samuelventura/go-tree"
)

func main() {
	tools.SetupLog()

	ctrlc := tools.SetupCtrlc()
	stdin := tools.SetupStdinAll()

	log.Println("start", os.Getpid())
	defer log.Println("exit")

	rnode := tree.NewRoot("root", log.Println)
	defer rnode.WaitDisposed()
	//recover closes as well
	defer rnode.Recover()
	rnode.SetValue("name", tools.GetEnviron("LV_NLK_NAME", tools.GetHostname()))
	rnode.SetValue("record", tools.GetEnviron("LV_NLK_DOCK_RECORD", "dock.laurelview.io"))
	rnode.SetValue("pool", tools.GetEnviron("LV_NLK_DOCK_POOL", ""))
	rnode.SetValue("count", newCount())
	rnode.SetValue("id", newId("ch"))

	enode := rnode.AddChild("ship")
	defer enode.WaitDisposed()
	defer enode.Close()
	run(enode)

	select {
	case <-rnode.Closed():
	case <-enode.Closed():
	case <-ctrlc:
	case <-stdin:
	}
}
