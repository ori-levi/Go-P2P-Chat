package main

import "flag"

const DEFAULT_PORT = 9090

func main() {

	var port int
	flag.IntVar(&port, "port", DEFAULT_PORT, "port for listening")

	var listingAllInterfaces bool
	flag.BoolVar(&listingAllInterfaces, "all-interfaces", false, "Should listining on all interfaces")

	flag.Parse()

	server := NewServer(port, listingAllInterfaces)
	server.RunServer()
}
