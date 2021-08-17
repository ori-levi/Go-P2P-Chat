package client

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"time"

	"levi.ori/p2p-chat/common"
)

var logger = common.NewLogger()

type Client struct {
	common.Client
	Connection []net.Conn
}

func NewClient(name string) Client {
	return Client{
		Client:     common.NewClient(name, nil),
		Connection: make([]net.Conn, 10),
	}
}

func (c *Client) Run(remoteUrl string, ic chan interface{}) {
	c.handleConnection(remoteUrl, false)
	time.Sleep(time.Hour)
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

	reader := bufio.NewReader(os.Stdin)
	for {
		if len(c.Name) == 0 {
			fmt.Print("Please enter your name: ")
			c.Name, _ = reader.ReadString('\n')
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
		c.Name, _ = reader.ReadString('\n')
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
