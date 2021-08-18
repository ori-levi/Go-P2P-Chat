package main

import (
	"flag"

	"levi.ori/p2p-chat/client"
	"levi.ori/p2p-chat/common"
	"levi.ori/p2p-chat/server"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "", "Your client name")

	var localServerPort int
	flag.IntVar(&localServerPort, "server-port", server.DefaultPort, "local server port for listening")

	var remoteUrl string
	flag.StringVar(&remoteUrl, "remote", "", "remote server url (format: <ip>:<port>)")

	var localInterfaceOnly bool
	flag.BoolVar(&localInterfaceOnly, "local-iface", false, "listening only for local interface ("+server.InternalInterface+")")

	flag.Parse()

	if len(remoteUrl) == 0 {
		common.Logger.Fatalf("Remote url is missing please run with -remote <ip>:<port>")
	}

	serverStartChannel := make(chan bool)
	serverApp := startServer(localServerPort, localInterfaceOnly, serverStartChannel)

	// waiting for server to up and running
	<-serverStartChannel

	clientApp := client.NewClient(name)
	clientApp.MakeInternalConnection(localServerPort)
	clientApp.Run(remoteUrl, serverApp.InternalClient.Channel)
}

func startServer(port int, internalIface bool, startChannel chan bool) *server.Server {
	app := server.NewServer(port, internalIface)
	go app.RunServer(startChannel)
	return &app
}
