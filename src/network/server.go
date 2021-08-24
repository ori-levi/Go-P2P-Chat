package network

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"levi.ori/p2p-chat/common"
)

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
	OutChannel chan string
	locker     sync.RWMutex
	port       int
	logChannel chan string
}

func NewServer(name string, port int, localInterfaceOnly bool, logChannel chan string) Server {
	address := AllInterfaces
	if localInterfaceOnly {
		address = InternalInterface
	}

	address = fmt.Sprintf("%v:%v", address, port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start server on %v", address)
	}

	return Server{
		name:       name,
		Listener:   ln,
		Clients:    make(map[string]*common.Client),
		InChannel:  make(chan interface{}),
		OutChannel: make(chan string),
		logChannel: logChannel,
		port:       port,
	}
}

func (s *Server) Close() {
	common.Debug(s.logChannel, "Server:Close Called")

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
		common.Error(s.logChannel, "Error occurred: %v", err)
	}
}

func (s *Server) RunServer() {
	defer s.Close()

	common.Info(s.logChannel, "Start listening to connection %v", s.Listener.Addr().String())
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			common.Error(s.logChannel, "Failed to accept connection, error: %v", err)
		}

		go s.makeClientConnection(conn)
	}
}

func (s *Server) makeClientConnection(conn net.Conn) {
	client := common.NewClient("", conn, s.logChannel, common.RandomColor())

	command, data := ReadCommand(&client)
	if command == common.Register {
		if serverPort, ok := s.registerClient(data, &client); ok && serverPort != s.port {
			common.Info(s.logChannel, "New connection from %v (%v)", client.Name, conn.RemoteAddr().String())
			s.InChannel <- common.InnerCommand{
				Command: common.ClientConnect,
				Data:    []interface{}{client.Name, client.RawConnection.RemoteAddr().String(), serverPort},
			}
			go s.handleConnection(&client)
		}
	} else {
		client.Close()
	}
}

func (s *Server) registerClient(data string, client *common.Client) (int, bool) {
	s.locker.Lock()
	defer s.locker.Unlock()

	parts := strings.Split(data, " ")
	lastIndex := len(parts) - 1
	name := strings.Join(parts[:lastIndex], " ")
	remotePort, err := common.AsInt(parts[lastIndex])
	if err != nil {
		common.Error(s.logChannel, "Failed to parse remote port| %v", err)
		return 0, false
	}

	if _, ok := s.Clients[name]; !ok && name != s.name {
		client.Name = name
		s.Clients[client.Name] = client
		if _, err := client.SendString(common.MyName, "%v\n", s.name); err != nil {
			common.Error(s.logChannel, "Failed to notify new user my name | %v", err)
		}
		return remotePort, true
	}
	if _, err := client.SendString(common.UserExists, "%v is already exists, please choose another name: ", name); err != nil {
		common.Error(s.logChannel, "Failed to notify user already exists | %v", err)
	}
	return 0, false
}

func ReadCommand(c *common.Client) (string, string) {
	_, rawData, _ := c.ReadAllAsString()
	parts := strings.SplitN(rawData, " ", 2)

	return parts[0], parts[1]
}

func (s *Server) handleConnection(client *common.Client) {
	code, data, err := client.ReadAllAsString()
	if client.Closed {
		s.removeConnection(client)
		return
	}

	if err != nil {
		common.Error(s.logChannel, "failed to read from client, error: %v", err)
	}

	name := client.Name
	if code == common.PM {
		name = fmt.Sprintf("(PM) %v", name)
	} else if code == common.Shell {
		name = fmt.Sprintf("(SHELL) %v", name)
	}

	name = common.ColorSprintf(client.Color, "%v:", name)
	s.OutChannel <- fmt.Sprintf("%v %v", name, data)

	go s.handleConnection(client)
}

func (s *Server) removeConnection(client *common.Client) {
	s.locker.Lock()
	defer s.locker.Unlock()

	client.Close()
	delete(s.Clients, client.Name)

	s.InChannel <- common.InnerCommand{
		Command: common.ClientDisconnect,
		Data:    client.Name,
	}
}
