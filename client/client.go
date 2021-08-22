package client

import (
	"bufio"
	"fmt"
	"levi.ori/p2p-chat/common"
	"net"
	"os"
	"strings"
	"sync"
)

var logger = common.NewLogger()

var commands = map[string]func(*Client, []string){
	"/pm":      pmCommand,
	"/connect": connectCommand,
	"/shell":   shellCommand,
}

type Client struct {
	common.Client
	Connections map[string]*common.Client
	reader      *bufio.Reader
	locker      sync.RWMutex
}

func NewClient(name string) Client {
	return Client{
		Client:      common.NewClient(name, nil),
		Connections: make(map[string]*common.Client),
		reader:      bufio.NewReader(os.Stdin),
	}
}

func (c *Client) Run(in chan interface{}, out chan interface{}) {
	go c.handleConnectionsFromServer(in)

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
				handler(c, arguments)
			} else {
				logger.Infof("Command %v is not recognized", command)
			}
		} else {
			sendToAll(c, data)
		}
	}
}

func (c *Client) handleConnectionsFromServer(ic chan interface{}) {
	for {
		rawMessage := <-ic
		msg, ok := rawMessage.(common.InnerCommand)
		if ok {
			if msg.Command == common.ClientDisconnect {
				clientName := msg.Data.(string)
				logger.Debugf("client %v disconnected", clientName)
				//c.removeConnection(clientName)
			} else if msg.Command == common.ClientConnect {
				client := msg.Data.(*common.Client)
				logger.Debugf("client %v connect", client.Name)
				//c.makeConnection(client)
			}
		}
	}
}

func (c *Client) removeConnection(clientName string) {
	c.locker.Lock()
	defer c.locker.Unlock()

	client, ok := c.Connections[clientName]
	if !ok {
		logger.Debugf("Client %v is not exists", clientName)
		return
	}

	client.Close()
	delete(c.Connections, clientName)
}

func (c *Client) makeConnection(addr string) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	client := common.NewClient("", conn)
	ok := register(&client, c.Name)
	if ok {
		logger.Debugf("Successfully connect to server")
		c.Connections[client.Name] = &client
		//go handleConnection(client)
	}
	return nil
}

func register(client *common.Client, name string) bool {
	_, err := client.SendString(common.Ok, "%v %v", common.Register, name)
	if err != nil {
		logger.Error("Failed to send command REGISTER")
		return false
	}

	code, data, err := client.ReadAllAsString()
	if err != nil {
		logger.Errorf("Failed to establish connection: %v", err)
		return false
	}

	if code != common.MyName {
		logger.Errorf("Failed to establish connection - get remote name: %v", err)
		return false
	}

	client.Name = data
	return true
}

func (c *Client) Close() {
	c.Client.Close()
	for _, client := range c.Connections {
		client.Close()
	}
}

func pmCommand(c *Client, arguments []string) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	logger.Debugf("PM command %v", arguments)

	name := arguments[0]
	conn, ok := c.Connections[name]
	if !ok {
		logger.Infof("Failed to send PM to %v, are you sure he is connected?", name)
		return
	}

	msg := strings.Join(arguments[1:], " ")
	conn.SendString(common.PM, msg)
}

func connectCommand(c *Client, arguments []string) {
	logger.Debugf("Connect command %v", arguments)

	addr := strings.Join(arguments[:2], ":")
	if err := c.makeConnection(addr); err != nil {
		logger.Infof("Failed to connect to %v; error %v", addr, err)
	}
}

func shellCommand(c *Client, arguments []string) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	logger.Debugf("Shell command %v", arguments)
}

func sendToAll(c *Client, msg string) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	for _, conn := range c.Connections {
		conn.SendString(common.Ok, msg)
	}
}
