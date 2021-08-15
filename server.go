package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Listener    net.Listener
	MegaChannel chan InnerCommand
	Clients     map[*Client]net.Addr
}

func NewServer(port int, listenAllInterfaces bool) Server {
	address := "127.0.0.1"
	if listenAllInterfaces {
		address = "0.0.0.0"
	}

	address = fmt.Sprintf("%v:%v", address, port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Failed to start server on", address)
	}

	return Server{
		Listener:    ln,
		MegaChannel: make(chan InnerCommand), // 32),
		Clients:     make(map[*Client]net.Addr),
	}
}

func (s *Server) Close() {
	fmt.Println("close called")
	s.Listener.Close()
	close(s.MegaChannel)
	for c, _ := range s.Clients {
		c.Close()
	}
}

func (s *Server) RunServer() {
	defer s.Close()

	fmt.Println("Start listening to connection", s.Listener.Addr().String())
	go s.handleMegaChannel()

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection, error:", err)
		}

		fmt.Println("Accept client from", conn.RemoteAddr().String())
		client := NewFromConnection(conn)

		s.Clients[&client] = client.RawConnection.RemoteAddr()
		go handleConnection(&client, s.MegaChannel)
	}
}

func (s *Server) handleMegaChannel() {
	message := <-s.MegaChannel

	if message.command == CLIENT_DISCONNECT {
		client := message.data.(*Client)
		fmt.Println("Client disconnect", client.RawConnection.RemoteAddr().String())
		client.Close()
		delete(s.Clients, client)
		for c := range s.Clients {
			c.Channel <- InnerCommand{
				command: CLIENT_DISCONNECT,
				data:    client,
			}
		}
	}

	go s.handleMegaChannel()
}

func handleConnection(client *Client, megaChannel chan InnerCommand) {
	if !client.ChannelStarted {
		go client.handleChanndel()
	}

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

	go handleConnection(client, megaChannel)
}
