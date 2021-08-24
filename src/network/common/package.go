package netwrok

import (
	"bufio"
	"fmt"
	"io"
	"levi.ori/p2p-chat/src/utils"
	"net"
	"strings"
)

type StatusCode int

const (
	statusCodeIndex = 0
	lengthIndex     = 1
)

const (
	Nil        = StatusCode(000)
	Ok         = StatusCode(200)
	MyName     = StatusCode(207)
	PM         = StatusCode(208)
	Shell      = StatusCode(209)
	UserExists = StatusCode(409)
)

func Send(c *net.Conn, statusCode StatusCode, data string) (int, error) {
	packed := fmt.Sprintf("%v %v\n%v", statusCode, len(data), data)
	return (*c).Write([]byte(packed))
}

func Read(c *net.Conn) (StatusCode, string, error) {
	reader := bufio.NewReader(*c)
	metadata, err := reader.ReadString('\n')
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
	if _, err := io.ReadFull(*c, data); err != nil {
		return statusCode, "", err
	}

	return statusCode, string(data), err
}
