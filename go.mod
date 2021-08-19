module levi.ori/p2p-chat

go 1.16

replace levi.ori/p2p-chat/server => ./server

replace levi.ori/p2p-chat/client => ./client

replace levi.ori/p2p-chat/common => ./common

require (
	github.com/jroimartin/gocui v0.5.0
	levi.ori/p2p-chat/client v0.0.0-00010101000000-000000000000
	levi.ori/p2p-chat/server v0.0.0-00010101000000-000000000000
)
