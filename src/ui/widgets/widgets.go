package ui

import (
	"fmt"
)

func newWidget(
	name string,
	title string,
	autoscroll bool,
	editable bool,
	wrap bool,
	x0, y0 PointCalculator,
	x1, y1 PointCalculator,
	data []string,
) *Widget {
	return &Widget{
		Name:          name,
		Title:         title,
		Autoscroll:    autoscroll,
		Editable:      editable,
		Wrap:          wrap,
		x0:            x0,
		y0:            y0,
		x1:            x1,
		y1:            y1,
		data:          data,
		OnValueChange: make(chan string),
	}
}

func NewInputWidget(name string) *Widget {
	return newWidget(
		"input",
		fmt.Sprintf("%v, What's On Your Mind?", name),
		true,
		true,
		true,
		func(maxX int) int { return 0 },
		func(maxY int) int { return 3*maxY/4 - 3 },
		func(maxX int) int { return maxX - 1 },
		func(maxY int) int { return 3*maxY/4 - 1 },
		nil,
	)
}

func NewLogWidget() *Widget {
	return newWidget(
		"log",
		"Log",
		true,
		false,
		true,
		func(maxX int) int { return 0 },
		func(maxY int) int { return 3 * maxY / 4 },
		func(maxX int) int { return maxX - 1 },
		func(maxY int) int { return maxY - 1 },
		nil,
	)
}

func NewChatWidget() *Widget {
	return newWidget(
		"chat",
		"Conversation",
		true,
		false,
		true,
		func(maxX int) int { return 0 },
		func(maxY int) int { return 0 },
		func(maxX int) int { return maxX / 3 * 2 },
		func(maxY int) int { return 3*maxY/4 - 4 },
		nil,
	)
}

func NewUsersWidget() *Widget {
	return newWidget(
		"users",
		"Users",
		true,
		false,
		true,
		func(maxX int) int { return maxX/3*2 + 1 },
		func(maxY int) int { return 0 },
		func(maxX int) int { return maxX - 1 },
		func(maxY int) int { return maxY / 2 },
		nil,
	)
}

func NewHelpWidget() *Widget {
	return newWidget(
		"help",
		"Help",
		true,
		false,
		true,
		func(maxX int) int { return maxX/3*2 + 1 },
		func(maxY int) int { return maxY/2 + 1 },
		func(maxX int) int { return maxX - 1 },
		func(maxY int) int { return 3*maxY/4 - 4 },
		[]string{
			fmt.Sprintf("%-9v <ip> <port>", "/connect"),
			fmt.Sprintf("%-9v <name> <message...>", "/pm"),
			fmt.Sprintf("%-9v <name> <command...>", "/shell"),
			"/exit",
		},
	)
}
