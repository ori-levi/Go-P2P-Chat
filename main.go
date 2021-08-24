package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "", "Your client name")

	var port int
	flag.IntVar(&port, "port", 9090, "local server port for listening")

	flag.Parse()

	if len(name) == 0 {
		log.Fatalln("Name is missing please run with -name <name>")
	}

	inputChannel := make(chan string)
	logChannel := make(chan string)

	fmt.Println(inputChannel)
	fmt.Println(logChannel)

	//serverApp := server.NewServer(name, port, localInterfaceOnly, logChannel)
	//go serverApp.RunServer()
	//
	//clientApp := client.NewClient(name, port, logChannel)
	//go clientApp.Run(serverApp.InChannel, inputChannel)
	//
	//old_code.uiMain(name, logChannel, serverApp.OutChannel, inputChannel)
}
