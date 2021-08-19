package server

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"levi.ori/p2p-chat/common"
)

var logger = common.Logger

const (
	InternalInterface = "127.0.0.1"
	AllInterfaces     = "0.0.0.0"
	DefaultPort       = 9090
)

type Server struct {
	name       string
	Listener   net.Listener
	Clients    map[string]*common.Client
	InChannel  chan interface{}
	OutChannel chan interface{}
	locker     sync.RWMutex
}

func NewServer(name string, port int, localInterfaceOnly bool) Server {
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
		name:       name,
		Listener:   ln,
		Clients:    make(map[string]*common.Client),
		InChannel:  make(chan interface{}),
		OutChannel: make(chan interface{}),
	}
}

func (s *Server) Close() {
	logger.Debug("Server:Close Called")

	// 1. Close all clients
	for _, c := range s.Clients {
		c.Close()
	}

	// 2. close channels
	close(s.InChannel)
	close(s.OutChannel)

	// 3. Close the listener object
	err := s.Listener.Close()
	if err != nil {
		logger.Errorf("Error occurred: %v", err)
	}
}

func (s *Server) RunServer() {
	defer s.Close()

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
		if command == common.Register && s.registerClient(data, &client) {
			logger.Infof("New connection from %v (%v)", client.Name, conn.RemoteAddr().String())
			s.InChannel <- common.InnerCommand{
				Command: common.ClientConnect,
				Data:    &client,
			}
			break
		}
	}

	// todo: do i need this?
	go s.handleConnection(&client)
}

func (s *Server) registerClient(name string, client *common.Client) bool {
	s.locker.Lock()
	defer s.locker.Unlock()

	if _, ok := s.Clients[name]; !ok && name != s.name {
		client.Name = name
		s.Clients[client.Name] = client
		client.SendString(common.MyName, "%v\n", s.name)
		return true
	}
	client.SendString(common.UserExists, "%v is already exists, please choose another name: ", name)
	return false
}

func ReadCommand(c *common.Client) (string, string) {
	rawData, _ := c.ReadAllAsString()
	parts := strings.SplitN(rawData, " ", 3)

	return parts[1], parts[2]
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

	s.InChannel <- common.InnerCommand{
		Command: common.ClientDisconnect,
		Data:    client.Name,
	}
}
