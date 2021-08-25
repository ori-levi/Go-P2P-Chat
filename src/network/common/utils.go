package network_common

type Command int

type innerCommand struct {
	Command
	Data interface{}
}

const (
	ClientDisconnect Command = iota
	ClientConnect
)

var (
	NotificationChannel = make(chan innerCommand)
)
