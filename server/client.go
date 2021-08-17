package server

import (
	"net"
	"strings"

	"levi.ori/p2p-chat/common"
)

type client struct {
	common.Client
	ChannelStarted bool
}

func newFromConnection(conn net.Conn) client {
	return client{
		Client:         common.NewClient("", conn),
		ChannelStarted: false,
	}
}

func (c *client) handleChannel() {
	c.ChannelStarted = true

	rawMessage := <-c.Channel
	if rawMessage != nil {
		msg := rawMessage.(InnerCommand)
		if msg.command == ClientDisconnect && msg.data != nil {
			client := msg.data.(*client)
			_, err := c.SendString(common.OK, "DISCONNECTED|%v|%v\n", client.Name, client.RawConnection.RemoteAddr().String())
			if err != nil {
				logger.Errorf("Error occurred: %v", err)
			}
		} else if msg.command == NewClient && msg.data != nil {
			client := msg.data.(*client)
			_, err := c.SendString(common.OK, "NewConnection|%v|%v\n", client.Name, client.RawConnection.RemoteAddr().String())
			if err != nil {
				logger.Errorf("Error occurred: %v", err)
			}
		}
	}

	go c.handleChannel()
}

func (c *client) ReadCommand() (string, string) {
	rawData, _ := c.ReadAllAsString()
	parts := strings.SplitN(rawData, "|", 2)

	return parts[0], parts[1]
}
