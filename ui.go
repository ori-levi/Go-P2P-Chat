package main

import (
	"bufio"
	"fmt"
	"io"
	"levi.ori/p2p-chat/common"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

type Widget struct {
	Name       string
	Title      string
	Editable   bool
	Autoscroll bool
	Wrap       bool
	x0, y0     func(int) int
	x1, y1     func(int) int
	data       []string
}

var (
	chatColors = map[string]common.Color{
		"(PM)":    common.Gold,
		"(SHELL)": common.Cyan,
	}

	logColors = map[string]common.Color{
		"[INFO":  common.LightPurple,
		"[DEBUG": common.LightCyan,
		"[ERROR": common.LightRed,
	}

	widgets = []Widget{
		{
			Name:       "input",
			Title:      "What's On Your Mind?",
			Autoscroll: true,
			Editable:   true,
			Wrap:       true,
			x0:         func(maxX int) int { return 0 },
			y0:         func(maxY int) int { return 3*maxY/4 - 3 },
			x1:         func(maxX int) int { return maxX - 1 },
			y1:         func(maxY int) int { return 3*maxY/4 - 1 },
		},
		{
			Name:       "log",
			Title:      "Log",
			Autoscroll: true,
			Editable:   false,
			Wrap:       true,
			x0:         func(maxX int) int { return 0 },
			y0:         func(maxY int) int { return 3 * maxY / 4 },
			x1:         func(maxX int) int { return maxX - 1 },
			y1:         func(maxY int) int { return maxY - 1 },
		},
		{
			Name:       "chat",
			Title:      "Conversation",
			Autoscroll: true,
			Editable:   false,
			Wrap:       true,
			x0:         func(maxX int) int { return 0 },
			y0:         func(maxY int) int { return 0 },
			x1:         func(maxX int) int { return maxX / 3 * 2 },
			y1:         func(maxY int) int { return 3*maxY/4 - 4 },
		},
		{
			Name:       "users",
			Title:      "Users",
			Autoscroll: true,
			Editable:   false,
			Wrap:       true,
			x0:         func(maxX int) int { return maxX/3*2 + 1 },
			y0:         func(maxY int) int { return 0 },
			x1:         func(maxX int) int { return maxX - 1 },
			y1:         func(maxY int) int { return maxY / 2 },
		},
		{
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
		},
	}
)

func layout(g *gocui.Gui, name string) error {
	maxX, maxY := g.Size()
	for _, w := range widgets {
		if v, err := g.SetView(w.Name, w.x0(maxX), w.y0(maxY), w.x1(maxX), w.y1(maxY)); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			title := w.Title
			if w.Name == "input" {
				title = fmt.Sprint(name, ", ", title)
			}

			v.Title = title
			v.Editable = w.Editable
			v.Wrap = w.Wrap
			v.Autoscroll = w.Autoscroll

			if len(w.data) > 0 {
				v.Clear()
				for _, d := range w.data {
					if _, err := fmt.Fprintln(v, d); err != nil {
						return err
					}
				}
			}
		}
	}

	_, err := g.SetCurrentView("input")
	return err
}

func quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}

func handleViewWithChannel(g *gocui.Gui, channel chan string, viewName string, formatter func(string) string, customAction func(string) bool) {
	for {
		msg := <-channel

		g.Update(func(g *gocui.Gui) error {
			v, err := g.View(viewName)
			if err != nil {
				return err
			}

			msg := strings.Trim(msg, "\r\n")
			if customAction == nil || customAction(msg) {
				if formatter != nil {
					msg = formatter(msg)
				}

				if _, err := fmt.Fprintln(v, msg); err != nil {
					return err
				}
			}
			return nil
		})
	}
}

func sendData(input chan string) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		reader := bufio.NewReader(v)
		data, err := reader.ReadString('\n')
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		data = strings.Trim(data, "\r\n")
		if data == "/exit" {
			return gocui.ErrQuit
		}

		input <- data
		g.Update(func(gui *gocui.Gui) error {
			vlog, err := g.View("chat")
			if err != nil {
				return err
			}

			color := common.ResetColor
			if strings.HasPrefix(data, "/") {
				color = common.LightGreen
			}

			if _, err := common.ColorFprintln(vlog, color, "ME:", data); err != nil {
				return err
			}

			return nil
		})

		v.Clear()
		if err := v.SetCursor(0, 0); err != nil {
			return nil
		}
		if err := v.SetOrigin(0, 0); err != nil {
			return nil
		}
		return nil
	}
}

func uiMain(name string, logChannel chan string, chatChannel chan string, inputChannel chan string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.InputEsc = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(func(g *gocui.Gui) error {
		return layout(g, name)
	})

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, sendData(inputChannel)); err != nil {
		log.Panicln(err)
	}

	go handleViewWithChannel(g, logChannel, "log", prefixFormatter(logColors), nil)
	go handleViewWithChannel(g, chatChannel, "chat", prefixFormatter(chatColors), shellPrompt)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func shellPrompt(s string) bool {
	if !strings.HasPrefix(s, "(SHELL)") {
		return true
	}

	return false
}

func prefixFormatter(prefixToColor map[string]common.Color) func(string) string {
	return func(s string) string {

		finalColor := common.ResetColor
		for prefix, color := range prefixToColor {
			if strings.HasPrefix(s, prefix) {
				finalColor = color
				break
			}
		}

		return common.ColorSprintf(finalColor, s)
	}
}
