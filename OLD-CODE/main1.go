package main1

import (
	"flag"
	"levi.ori/p2p-chat/OLD-CODE/client"
	"levi.ori/p2p-chat/OLD-CODE/server"
	"log"
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
		log.Fatalln("Name is missing please run with -name <name>")
	}

	inputChannel := make(chan string)
	logChannel := make(chan string)

	serverApp := server.NewServer(name, port, localInterfaceOnly, logChannel)
	go serverApp.RunServer()

	clientApp := client.NewClient(name, port, logChannel)
	go clientApp.Run(serverApp.InChannel, inputChannel)

	OLD_CODE.uiMain(name, logChannel, serverApp.OutChannel, inputChannel)
}
