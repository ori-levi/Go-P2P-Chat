package main

const (
	CLIENT_DISCONNECT = iota
	CLOSE_ALL
)

type InnerCommand struct {
	command int
	data    interface{}
}
