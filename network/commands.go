package network

const (
	Register = "REGISTER"
)

const (
	ClientDisconnect = iota
	ClientConnect
)

type InnerCommand struct {
	Command int
	Data    interface{}
}
