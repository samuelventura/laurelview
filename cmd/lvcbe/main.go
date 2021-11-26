package main

import (
	"log"
	"os"

	"github.com/samuelventura/go-state"
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
	rnode.SetValue("endpoint", tools.GetEnviron("LV_CBE_ENDPOINT", "127.0.0.1:31603"))
	rnode.SetValue("state", tools.GetEnviron("LV_CBE_STATE", tools.WithExtension("state")))

	snode := state.Serve(rnode, rnode.GetValue("state").(string))
	defer snode.WaitDisposed()
	defer snode.Close()

	anode := rnode.AddChild("api")
	defer anode.WaitDisposed()
	defer anode.Close()
	api(anode)

	select {
	case <-rnode.Closed():
	case <-snode.Closed():
	case <-anode.Closed():
	case <-ctrlc:
	case <-stdin:
	}
}
