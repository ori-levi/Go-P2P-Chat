package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Listener       net.Listener
	MegaChannel    InnerCommandChannel
	ClientChannels map[Client]InnerCommandChannel
}

func NewServer(port int, listenAllInterfaces bool) Server {
	address := ""
	if listenAllInterfaces {
		address = "0.0.0.0"
	}

	address = fmt.Sprintf("%v:%v", address, port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Failed to start server on port", port)
	}

	return Server{
		Listener:       ln,
		MegaChannel:    make(InnerCommandChannel),
		ClientChannels: make(map[Client]InnerCommandChannel),
	}
}

func handleMegaChannel(s *Server, message InnerCommand) {
	if message.command == CLIENT_DISCONNECT {
		client := message.data.(Client)
		fmt.Println("bye bye", client.RawConnection.LocalAddr().String(), client.RawConnection.RemoteAddr().String())
		close(s.ClientChannels[client])
		delete(s.ClientChannels, client)
		for _, c := range s.ClientChannels {
			c <- InnerCommand{
				command: CLIENT_DISCONNECT,
				data:    client,
			}
		}
	}

}

func (s *Server) RunServer() {
	defer s.Close()

	fmt.Println("Starting listening to connection", s.Listener.Addr().String())
loop:
	for {
		select {
		case msg := <-s.MegaChannel:
			if msg.command == CLOSE_ALL {
				break loop
			}
			handleMegaChannel(s, msg)
		default:
			s.acceptNewConnection()
		}
	}
}

func (s *Server) acceptNewConnection() {
	conn, err := s.Listener.Accept()
	if err != nil {
		fmt.Println("Failed to accept connection, error:", err)
	}

	fmt.Println("Accept client from", conn.RemoteAddr().String())
	client := NewFromConnection(conn)

	channel := make(InnerCommandChannel)
	s.ClientChannels[client] = channel
	go handleConnection(&client, channel, s.MegaChannel)
}

func (s *Server) Close() {
	fmt.Println("close called")
	s.Listener.Close()
	close(s.MegaChannel)
	for _, c := range s.ClientChannels {
		close(c)
	}
}

func handleConnection(client *Client, clientChannel InnerCommandChannel, megaChannel InnerCommandChannel) {
	select {
	case msg := <-clientChannel:
		if msg.command == CLIENT_DISCONNECT {
			client.SendString("DISCONNECTED:%v", client.RawConnection.RemoteAddr().String())
		}
	default:
		data, err := client.Read()
		if client.Closed {
			megaChannel <- InnerCommand{
				command: CLIENT_DISCONNECT,
				data:    client,
			}
			return
		}

		if err != nil {
			fmt.Println("failed to read from client, error:", err)
		}

		fmt.Println("Got from client", data)
		_, err = client.SendString("Answer: %v", data)
		if err != nil {
			fmt.Println("Failed to send data")
		}
	}

	go handleConnection(client, clientChannel, megaChannel)
}
