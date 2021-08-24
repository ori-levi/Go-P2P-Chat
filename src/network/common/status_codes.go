package network_common

type StatusCode int

const (
	Nil        = StatusCode(000)
	Ok         = StatusCode(200)
	MyName     = StatusCode(207)
	PM         = StatusCode(208)
	Shell      = StatusCode(209)
	UserExists = StatusCode(409)
)
