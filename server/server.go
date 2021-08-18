package server

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"levi.ori/p2p-chat/common"
)

type Server struct {
	Listener       net.Listener
	InternalClient *common.Client
	Clients        map[string]*common.Client
	locker         sync.RWMutex
}

func NewServer(port int, localInterfaceOnly bool) Server {
	address := AllInterfaces
	if localInterfaceOnly {
		address = InternalInterface
	}

	address = fmt.Sprintf("%v:%v", address, port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("Failed to start server on %v", address)
	}

	return Server{
		Listener: ln,
		Clients:  make(map[string]*common.Client),
	}
}

func (s *Server) Close() {
	logger.Debug("Server:Close Called")

	// 1. Close all clients
	{
		s.locker.Lock()
		defer s.locker.Unlock()

		for _, c := range s.Clients {
			c.Close()
		}
	}

	// 2. Close the listener object
	err := s.Listener.Close()
	if err != nil {
		logger.Errorf("Error occurred: %v", err)
	}
}

func (s *Server) RunServer(serverStartChannel chan bool) {
	defer s.Close()
	serverStartChannel <- true

	logger.Infof("Start listening to connection %v", s.Listener.Addr().String())

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			logger.Errorf("Failed to accept connection, error: %v", err)
		}

		go s.makeClientConnection(conn)
	}
}

func (s *Server) makeClientConnection(conn net.Conn) {
	client := common.NewClient("", conn)

	for {
		command, data := ReadCommand(&client)
		if command == common.InternalRegister {
			logger.Debugf("Internal client connected %v|%v", command, data)
			s.InternalClient = &client
			client.SendString(common.Ok, "OK\n")
			break
		}

		if command == common.Register && s.registerClient(data, &client) {
			logger.Infof("New connection from %v (%v)", client.Name, conn.RemoteAddr().String())
			if s.InternalClient != nil {
				s.InternalClient.Channel <- common.InnerCommand{
					Command: common.ClientConnect,
					Data:    &client,
				}
			}
			break
		}
	}

	go s.handleConnection(&client)
}

func (s *Server) registerClient(name string, client *common.Client) bool {
	s.locker.Lock()
	defer s.locker.Unlock()

	if _, ok := s.Clients[name]; !ok {

		client.Name = name
		s.Clients[client.Name] = client
		client.SendString(common.Ok, "Welcome %v!\n", name)

		return true
	}
	client.SendString(common.UserExists, "%v is already exists, please choose another name: ", name)
	return false
}

func ReadCommand(c *common.Client) (string, string) {
	rawData, _ := c.ReadAllAsString()
	parts := strings.SplitN(rawData, "|", 2)

	return parts[0], parts[1]
}

func (s *Server) handleConnection(client *common.Client) {
	_, err := client.ReadAllAsString()
	if client.Closed {
		s.removeConnection(client)
		return
	}

	if err != nil {
		logger.Errorf("failed to read from client, error: %v", err)
	}

	go s.handleConnection(client)
}

func (s *Server) removeConnection(client *common.Client) {
	s.locker.Lock()
	defer s.locker.Unlock()

	delete(s.Clients, client.Name)

	s.InternalClient.Channel <- common.InnerCommand{
		Command: common.ClientDisconnect,
		Data:    client,
	}
}
