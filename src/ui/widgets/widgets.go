package ui

import "fmt"

func NewInputWidget(name string) *Widget {
	return &Widget{
		Name:       "input",
		Title:      fmt.Sprintf("%v, What's On Your Mind?", name),
		Autoscroll: true,
		Editable:   true,
		Wrap:       true,
		x0:         func(maxX int) int { return 0 },
		y0:         func(maxY int) int { return 3*maxY/4 - 3 },
		x1:         func(maxX int) int { return maxX - 1 },
		y1:         func(maxY int) int { return 3*maxY/4 - 1 },
	}
}

func NewLogWidget() *Widget {
	return &Widget{
		Name:       "log",
		Title:      "Log",
		Autoscroll: true,
		Editable:   false,
		Wrap:       true,
		x0:         func(maxX int) int { return 0 },
		y0:         func(maxY int) int { return 3 * maxY / 4 },
		x1:         func(maxX int) int { return maxX - 1 },
		y1:         func(maxY int) int { return maxY - 1 },
	}
}

func NewChatWidget() *Widget {
	return &Widget{
		Name:       "chat",
		Title:      "Conversation",
		Autoscroll: true,
		Editable:   false,
		Wrap:       true,
		x0:         func(maxX int) int { return 0 },
		y0:         func(maxY int) int { return 0 },
		x1:         func(maxX int) int { return maxX / 3 * 2 },
		y1:         func(maxY int) int { return 3*maxY/4 - 4 },
	}
}

func NewUsersWidget() *Widget {
	return &Widget{
		Name:       "users",
		Title:      "Users",
		Autoscroll: true,
		Editable:   false,
		Wrap:       true,
		x0:         func(maxX int) int { return maxX/3*2 + 1 },
		y0:         func(maxY int) int { return 0 },
		x1:         func(maxX int) int { return maxX - 1 },
		y1:         func(maxY int) int { return maxY / 2 },
	}
}

func NewHelpWidget() *Widget {
	return &Widget{
		Name:       "help",
		Title:      "Help",
		Autoscroll: true,
		Editable:   false,
		Wrap:       true,
		x0:         func(maxX int) int { return maxX/3*2 + 1 },
		y0:         func(maxY int) int { return maxY/2 + 1 },
		x1:         func(maxX int) int { return maxX - 1 },
		y1:         func(maxY int) int { return 3*maxY/4 - 4 },
		data: []string{
			fmt.Sprintf("%-9v <ip> <port>", "/connect"),
			fmt.Sprintf("%-9v <name> <message...>", "/pm"),
			fmt.Sprintf("%-9v <name> <command...>", "/shell"),
			"/exit",
		},
	}
}
