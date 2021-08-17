package common

const (
	Register         = "REGISTER"
	InternalRegister = "INTERNAL-REGISTER"
)

const (
	ClientDisconnect = iota
	ClientConnect
)

type InnerCommand struct {
	Command int
	Data    interface{}
}
