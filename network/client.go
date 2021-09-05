package network

import (
	"bufio"
	"fmt"
	"levi.ori/p2p-chat/utils"
	"net"
	"os"
	"strings"
	"sync"
)

var commands = map[string]func(*Client, []string){
	"/pm":      pmCommand,
	"/connect": connectCommand,
	"/shell":   shellCommand,
}

type Client struct {
	Conn
	Connections map[string]*Conn
	reader      *bufio.Reader
	locker      sync.RWMutex
	serverPort  int
	logChannel  chan string
}

func NewClient(name string, serverPort int, logChannel chan string) Client {
	return Client{
		Conn:        NewConn(name, nil, logChannel, utils.ResetColor),
		Connections: make(map[string]*Conn),
		reader:      bufio.NewReader(os.Stdin),
		serverPort:  serverPort,
		logChannel:  logChannel,
	}
}

func (c *Client) Run(serverNotification chan interface{}, input chan string) {
	go c.handleConnectionsFromServer(serverNotification)

	for {
		data := <-input
		if strings.HasPrefix(data, "/") {
			parts := strings.Split(data, " ")
			command, arguments := parts[0], parts[1:]

			if handler, ok := commands[command]; ok {
				handler(c, arguments)
			} else {
				utils.Info(c.logChannel, "Command %v is not recognized", command)
			}
		} else {
			sendToAll(c, data)
		}
	}
}

func (c *Client) handleConnectionsFromServer(ic chan interface{}) {
	for {
		rawMessage := <-ic
		msg, ok := rawMessage.(InnerCommand)
		if ok {
			if msg.Command == ClientDisconnect {
				clientName := msg.Data.(string)
				utils.Debug(c.logChannel, "client %v disconnected", clientName)
				c.removeConnection(clientName)
			} else if msg.Command == ClientConnect {
				parts := msg.Data.([]interface{})
				clientName := parts[0].(string)
				remoteAddr := strings.Split(parts[1].(string), ":")
				remotePort := parts[2].(int)

				addr := fmt.Sprintf("%v:%v", remoteAddr[0], remotePort)
				utils.Debug(c.logChannel, "client %v connect %v", clientName, addr)
				if err := c.makeConnection(addr); err != nil {
					utils.Error(c.logChannel, "Failed to make connection | %v", err)
				}
			}
		}
	}
}

func (c *Client) removeConnection(clientName string) {
	c.locker.Lock()
	defer c.locker.Unlock()

	client, ok := c.Connections[clientName]
	if !ok {
		utils.Debug(c.logChannel, "Client %v is not exists", clientName)
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

	client := NewConn("", conn, c.logChannel, utils.ResetColor)
	ok := c.register(&client, c.Name, c.serverPort)
	if ok {
		utils.Debug(c.logChannel, "Successfully connect to server")
		c.Connections[client.Name] = &client
	}
	return nil
}

func (c *Client) register(client *Conn, name string, serverPort int) bool {
	_, err := client.SendString(Ok, "%v %v %v", Register, name, serverPort)
	if err != nil {
		utils.Error(c.logChannel, "Failed to send command REGISTER")
		return false
	}

	code, data, err := client.ReadAllAsString()
	if err != nil {
		utils.Error(c.logChannel, "Failed to establish connection: %v", err)
		return false
	}

	if code != MyName {
		utils.Error(c.logChannel, "Failed to establish connection - get remote name: %v", err)
		return false
	}

	client.Name = data
	return true
}

func (c *Client) Close() {
	c.Conn.Close()
	for _, client := range c.Connections {
		client.Close()
	}
}

func pmCommand(c *Client, arguments []string) {
	c.sendPrivate(arguments, PM, "PM")
}

func connectCommand(c *Client, arguments []string) {
	utils.Debug(c.logChannel, "Connect command %v", arguments)

	addr := strings.Join(arguments[:2], ":")
	if err := c.makeConnection(addr); err != nil {
		utils.Info(c.logChannel, "Failed to connect to %v; error %v", addr, err)
	}
}

func shellCommand(c *Client, arguments []string) {
	c.sendPrivate(arguments, Shell, "Shell")
}

func (c *Client) sendPrivate(arguments []string, code int, command string) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	utils.Debug(c.logChannel, "%v command %v", command, arguments)

	name := arguments[0]
	conn, ok := c.Connections[name]
	if !ok {
		utils.Info(c.logChannel, "Failed to send %v to %v, are you sure he is connected?", command, name)
		return
	}

	msg := strings.Join(arguments[1:], " ")
	if _, err := conn.SendString(code, msg); err != nil {
		utils.Error(c.logChannel, "Failed to send %v to %v | %v", command, name, err)
	}
}

func sendToAll(c *Client, msg string) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	for name, conn := range c.Connections {
		if _, err := conn.SendString(Ok, msg); err != nil {
			utils.Error(c.logChannel, "Failed to send message to %v | %v", name, err)
		}
	}
}
