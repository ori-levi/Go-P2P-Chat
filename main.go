package main

import (
	"flag"
	"fmt"
)

const DEFAULT_PORT = 9090

func main() {

	var port int
	flag.IntVar(&port, "port", DEFAULT_PORT, fmt.Sprintf("port for listening [Default: %v]", DEFAULT_PORT))

	flag.Parse()

	server := NewServer(port, false)
	server.RunServer()
}
