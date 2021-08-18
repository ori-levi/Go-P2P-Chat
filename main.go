package main

import (
	"flag"
	"levi.ori/p2p-chat/common"

	"levi.ori/p2p-chat/server"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "", "Your client name")

	var port int
	flag.IntVar(&port, "port", server.DefaultPort, "local server port for listening")

	var localInterfaceOnly bool
	flag.BoolVar(&localInterfaceOnly, "local-iface", false, "listening only for local interface ("+server.InternalInterface+")")

	flag.Parse()

	if len(name) == 0 {
		common.Logger.Fatalf("Name is missing please run with -name <name>")
	}

	serverApp := server.NewServer(port, localInterfaceOnly)
	serverApp.RunServer()
	//go serverApp.RunServer()

	//clientApp := client.NewClient(name)
	//clientApp.MakeInternalConnection(port)
	//clientApp.Run(serverApp.InternalClient.Channel)
}
