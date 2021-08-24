package network_common

import (
	"bufio"
	"fmt"
	"io"
	"levi.ori/p2p-chat/src/utils"
	"levi.ori/p2p-chat/src/utils/colors"
	"net"
	"strings"
)

const (
	statusCodeIndex = 0
	lengthIndex     = 1
)

type Client struct {
	RawConnection net.Conn
	Reader        *bufio.Reader
	Closed        bool
	Channel       chan interface{}
	Name          string
	logChannel    chan string
	colors.Color
}

func NewClient(name string, connection net.Conn, logChannel chan string, color colors.Color) Client {
	c := Client{
		Closed:     false,
		Channel:    make(chan interface{}),
		Name:       name,
		logChannel: logChannel,
		Color:      color,
	}
	c.setRawConnection(connection)
	return c
}

func (c *Client) setRawConnection(conn net.Conn) {
	if conn == nil {
		return
	}

	c.RawConnection = conn
	c.Reader = bufio.NewReader(conn)
}

func (c *Client) Close() {
	close(c.Channel)
	err := c.RawConnection.Close()
	if err != nil {
		//todo enable again log channel
		//Error(c.logChannel, "Failed to close connection; %v", err)
	}
}

func (c *Client) Send(statusCode StatusCode, data string) (int, error) {
	packed := fmt.Sprintf("%v %v\n%v", statusCode, len(data), data)
	return c.RawConnection.Write([]byte(packed))
}

func (c *Client) Read() (StatusCode, string, error) {
	metadata, err := c.Reader.ReadString('\n')
	if err != nil {
		return Nil, "", err
	}

	parts := strings.Split(metadata, " ")

	rawStatusCode, err := utils.AsInt(parts[statusCodeIndex])
	statusCode := StatusCode(rawStatusCode)
	if err != nil {
		return statusCode, "", err
	}

	length, err := utils.AsInt(parts[lengthIndex])
	if err != nil {
		return statusCode, "", err
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(c.RawConnection, data); err != nil {
		return statusCode, "", err
	}

	return statusCode, string(data), err
}
