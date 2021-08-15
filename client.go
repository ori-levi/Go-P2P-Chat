package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

const (
	ALL      = "ALL"
	GET_FIEL = "get-file"
)

type Protocol struct {
	action string
	data   string
}

type Client struct {
	RawConnection net.Conn
	Reader        *bufio.Reader
	Closed        bool
}

func NewFromConnection(conn net.Conn) Client {
	return Client{
		RawConnection: conn,
		Reader:        bufio.NewReader(conn),
		Closed:        false,
	}
}

func (c *Client) ReadAllAsString() (string, error) {
	data, err := c.Reader.ReadString('\n')
	if err != nil {
		c.Closed = err == io.EOF
		return "", err
	}

	return strings.Trim(data, "\r\n"), nil
}

func (c *Client) Read() (Protocol, error) {
	data, err := c.ReadAllAsString()
	if err != nil {
		return Protocol{}, err
	}

	// parts := strings.SplitN(data, ":", 2)
	return Protocol{
		action: data,
		data:   data,
	}, nil

}

func (c *Client) SendString(format string, args ...interface{}) (int, error) {
	data := fmt.Sprintf(format, args...)
	return c.RawConnection.Write([]byte(data))
}
