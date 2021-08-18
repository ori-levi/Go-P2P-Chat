package client

import (
	"bufio"
	"fmt"
	"levi.ori/p2p-chat/common"
	"net"
	"net/textproto"
	"os"
	"strings"
	"sync"
)

var logger = common.NewLogger()

var commands = map[string]func([]string){
	"/pm":      pmCommand,
	"/connect": connectCommand,
	"/shell":   shellCommand,
}

type Client struct {
	common.Client
	Connection []net.Conn
	reader     *bufio.Reader
	locker     sync.RWMutex
}

func NewClient(name string) Client {
	return Client{
		Client:     common.NewClient(name, nil),
		Connection: make([]net.Conn, 10),
		reader:     bufio.NewReader(os.Stdin),
	}
}

func (c *Client) Run(ic chan interface{}) {
	//go c.handleInternalChannel(ic)

	for {
		fmt.Print(">>> ")
		data, err := c.reader.ReadString('\n')
		if err != nil {
			logger.Fatalf("Failed to read from stdin, error: %v", err)
		}

		data = strings.Trim(data, " \n\r")
		if strings.HasPrefix(data, "/") {
			parts := strings.Split(data, " ")
			command, arguments := parts[0], parts[1:]

			if command == "/exit" {
				break
			}

			if handler, ok := commands[command]; ok {
				handler(arguments)
			} else {
				logger.Infof("Command %v is not recognized", command)
			}
		} else {
			c.sendToAll(data)
		}
	}
}

func (c *Client) handleInternalChannel(ic chan interface{}) {
	for {
		rawMessage := <-ic
		msg, ok := rawMessage.(common.InnerCommand)
		if ok {
			if msg.Command == common.ClientDisconnect {
				client := msg.Data.(*common.Client)
				logger.Debugf("client %v disconnected", client.RawConnection.RemoteAddr().String())
			} else if msg.Command == common.ClientConnect {
				client := msg.Data.(*common.Client)
				logger.Debugf("client %v connect", client.RawConnection.RemoteAddr().String())
			}
		}
	}
}

func (c *Client) MakeInternalConnection(serverPort int) {
	addr := fmt.Sprintf("127.0.0.1:%v", serverPort)

	err := c.handleConnection(addr, true)
	if err != nil {
		logger.Fatalf("Failed to establish internal connection; %v", err)
	}
}

func (c *Client) handleConnection(addr string, internal bool) error {
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return err
	}

	for {
		if len(c.Name) == 0 {
			fmt.Print("Please enter your name: ")
			c.Name, _ = c.reader.ReadString('\n')
		}

		id := register(conn, c.Name, internal)
		conn.StartResponse(id)
		defer conn.EndResponse(id)

		code, message, err := conn.ReadCodeLine(200)
		if err == nil {
			logger.Debugf("Successfully connect to server")
			break
		}

		fmt.Printf("%v %v", code, message)
		c.Name, _ = c.reader.ReadString('\n')
	}

	return nil
}

func register(conn *textproto.Conn, name string, internal bool) uint {
	cmd := common.Register
	if internal {
		cmd = common.InternalRegister
	}

	id, err := conn.Cmd("%v|%v", cmd, name)
	if err != nil {
		logger.Fatalf("Failed to send command REGISTER")
	}
	return id
}

func (c *Client) Close() {
	c.Client.Close()
	for _, client := range c.Connection {
		err := client.Close()
		if err != nil {
			logger.Errorf("Failed to close connection; %v", err)
		}
	}
}

func pmCommand(arguments []string) {
	logger.Debugf("PM command %v", arguments)
}

func connectCommand(arguments []string) {
	logger.Debugf("Connect command %v", arguments)
}

func shellCommand(arguments []string) {
	logger.Debugf("Shell command %v", arguments)
}

func (c *Client) sendToAll(msg string) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	for _, conn := range c.Connection {
		conn.Write([]byte(msg))
	}
}
