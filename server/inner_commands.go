package server

const (
	ClientDisconnect = iota
	NewClient
)

type InnerCommand struct {
	command int
	data    interface{}
}
