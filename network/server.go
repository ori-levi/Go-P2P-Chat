package network

import (
	"fmt"
	"levi.ori/p2p-chat/utils"
	"log"
	"net"
	"strings"
	"sync"
)

const (
	InternalInterface = "127.0.0.1"
	AllInterfaces     = "0.0.0.0"
	DefaultPort       = 9090
)

type Server struct {
	name       string
	Listener   net.Listener
	Clients    map[string]*Conn
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
		Clients:    make(map[string]*Conn),
		InChannel:  make(chan interface{}),
		OutChannel: make(chan string),
		logChannel: logChannel,
		port:       port,
	}
}

func (s *Server) Close() {
	utils.Debug(s.logChannel, "Server:Close Called")

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
		utils.Error(s.logChannel, "Error occurred: %v", err)
	}
}

func (s *Server) RunServer() {
	defer s.Close()

	utils.Info(s.logChannel, "Start listening to connection %v", s.Listener.Addr().String())
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			utils.Error(s.logChannel, "Failed to accept connection, error: %v", err)
		}

		go s.makeClientConnection(conn)
	}
}

func (s *Server) makeClientConnection(conn net.Conn) {
	client := NewConn("", conn, s.logChannel, utils.RandomColor())

	command, data := ReadCommand(&client)
	if command == Register {
		if serverPort, ok := s.registerClient(data, &client); ok && serverPort != s.port {
			utils.Info(s.logChannel, "New connection from %v (%v)", client.Name, conn.RemoteAddr().String())
			s.InChannel <- InnerCommand{
				Command: ClientConnect,
				Data:    []interface{}{client.Name, client.RawConnection.RemoteAddr().String(), serverPort},
			}
			go s.handleConnection(&client)
		}
	} else {
		client.Close()
	}
}

func (s *Server) registerClient(data string, client *Conn) (int, bool) {
	s.locker.Lock()
	defer s.locker.Unlock()

	parts := strings.Split(data, " ")
	lastIndex := len(parts) - 1
	name := strings.Join(parts[:lastIndex], " ")
	remotePort, err := utils.AsInt(parts[lastIndex])
	if err != nil {
		utils.Error(s.logChannel, "Failed to parse remote port| %v", err)
		return 0, false
	}

	if _, ok := s.Clients[name]; !ok && name != s.name {
		client.Name = name
		s.Clients[client.Name] = client
		if _, err := client.SendString(MyName, "%v\n", s.name); err != nil {
			utils.Error(s.logChannel, "Failed to notify new user my name | %v", err)
		}
		return remotePort, true
	}
	if _, err := client.SendString(UserExists, "%v is already exists, please choose another name: ", name); err != nil {
		utils.Error(s.logChannel, "Failed to notify user already exists | %v", err)
	}
	return 0, false
}

func ReadCommand(c *Conn) (string, string) {
	_, rawData, _ := c.ReadAllAsString()
	parts := strings.SplitN(rawData, " ", 2)

	return parts[0], parts[1]
}

func (s *Server) handleConnection(client *Conn) {
	code, data, err := client.ReadAllAsString()
	if client.Closed {
		s.removeConnection(client)
		return
	}

	if err != nil {
		utils.Error(s.logChannel, "failed to read from client, error: %v", err)
	}

	name := client.Name
	if code == PM {
		name = fmt.Sprintf("(PM) %v", name)
	} else if code == Shell {
		name = fmt.Sprintf("(SHELL) %v", name)
	}

	name = utils.ColorSprintf(client.Color, "%v:", name)
	s.OutChannel <- fmt.Sprintf("%v %v", name, data)

	go s.handleConnection(client)
}

func (s *Server) removeConnection(client *Conn) {
	s.locker.Lock()
	defer s.locker.Unlock()

	client.Close()
	delete(s.Clients, client.Name)

	s.InChannel <- InnerCommand{
		Command: ClientDisconnect,
		Data:    client.Name,
	}
}
