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
	rnode.SetValue("driver", tools.GetEnviron("LV_CBE_DB_DRIVER", "sqlite"))
	rnode.SetValue("source", tools.GetEnviron("LV_CBE_DB_SOURCE", tools.WithExtension("db3")))
	rnode.SetValue("stom", tools.GetEnvironInt("LV_CBE_SESSION_TO", 10, 32, 15)) //minutes
	rnode.SetValue("webep", tools.GetEnviron("LV_CBE_ENDPOINT", "127.0.0.1:5003"))
	dao := newDao(rnode) //close on root
	rnode.AddAction("dao", dao.close)
	rnode.SetValue("dao", dao)
	api := newApi(rnode)
	rnode.SetValue("api", api)

	wnode := rnode.AddChild("web")
	defer wnode.WaitDisposed()
	defer wnode.Close()
	web(wnode)

	select {
	case <-rnode.Closed():
	case <-wnode.Closed():
	case <-ctrlc:
	case <-stdin:
	}
}
