package network

import (
	"bufio"
	"fmt"
	"io"
	"levi.ori/p2p-chat/utils"
	"net"
	"strings"
)

type Conn struct {
	RawConnection net.Conn
	Reader        *bufio.Reader
	Closed        bool
	Channel       chan interface{}
	Name          string
	logChannel    chan string
	utils.Color
}

func NewConn(name string, connection net.Conn, logChannel chan string, color utils.Color) Conn {
	c := Conn{
		Closed:     false,
		Channel:    make(chan interface{}),
		Name:       name,
		logChannel: logChannel,
		Color:      color,
	}
	c.SetRawConnection(connection)
	return c
}

func (c *Conn) SetRawConnection(conn net.Conn) {
	if conn == nil {
		return
	}

	c.RawConnection = conn
	c.Reader = bufio.NewReader(conn)
}

func (c *Conn) Close() {
	close(c.Channel)
	err := c.RawConnection.Close()
	if err != nil {
		utils.Error(c.logChannel, "Failed to close connection; %v", err)
	}
}

func (c *Conn) ReadAllAsString() (int, string, error) {
	data, err := c.Reader.ReadString('\n')
	if err != nil {
		c.Closed = err == io.EOF
		return 0, "", err
	}

	data = strings.Trim(data, "\r\n")
	parts := strings.SplitN(data, " ", 2)

	code, err := utils.AsInt(parts[0])
	if err != nil {
		utils.Error(c.logChannel, "Failed to parse code as int| %v", err)
		return 0, "", err
	}

	return code, strings.Join(parts[1:], " "), nil
}

func (c *Conn) SendString(code int, format string, args ...interface{}) (int, error) {
	fullFormat := fmt.Sprintf("%v %v", code, format)
	data := fmt.Sprintf(fullFormat, args...)
	if !strings.HasSuffix(data, "\n") {
		data += "\n"
	}
	return c.RawConnection.Write([]byte(data))
}
