package common

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

var logger = NewLogger()

type Client struct {
	RawConnection net.Conn
	Reader        *bufio.Reader
	Closed        bool
	Channel       chan interface{}
	Name          string
}

func NewClient(name string, connection net.Conn) Client {
	c := Client{
		Closed:  false,
		Channel: make(chan interface{}),
		Name:    name,
	}
	c.SetRawConnection(connection)
	return c
}

func (c *Client) SetRawConnection(conn net.Conn) {
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
		logger.Errorf("Failed to close connection; %v", err)
	}
}

func (c *Client) ReadAllAsString() (int, string, error) {
	data, err := c.Reader.ReadString('\n')
	if err != nil {
		c.Closed = err == io.EOF
		return 0, "", err
	}

	data = strings.Trim(data, "\r\n")
	parts := strings.SplitN(data, " ", 2)

	code, err := AsInt(parts[0])
	if err != nil {
		logger.Errorf("Failed to establish connection - get remote name: %v", err)
		return 0, "", err
	}

	return code, strings.Join(parts[1:], " "), nil
}

func (c *Client) SendString(code int, format string, args ...interface{}) (int, error) {
	fullFormat := fmt.Sprintf("%v %v", code, format)
	data := fmt.Sprintf(fullFormat, args...)
	if !strings.HasSuffix(data, "\n") {
		data += "\n"
	}
	return c.RawConnection.Write([]byte(data))
}
