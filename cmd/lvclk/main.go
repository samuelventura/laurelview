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
	rnode.SetValue("endpoint", tools.GetEnviron("LV_CLK_ENDPOINT_SSH", "0.0.0.0:31622"))
	rnode.SetValue("proxy", tools.GetEnviron("LV_CLK_ENDPOINT_PROXY", "127.0.0.1:31680"))
	rnode.SetValue("single", newSingle())

	enode := rnode.AddChild("ssh")
	defer enode.WaitDisposed()
	defer enode.Close()
	sshd(enode)

	select {
	case <-rnode.Closed():
	case <-enode.Closed():
	case <-ctrlc:
	case <-stdin:
	}
}
