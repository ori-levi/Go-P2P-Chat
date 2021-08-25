package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	app "levi.ori/p2p-chat/src/ui"
	"levi.ori/p2p-chat/src/utils/colors"
	"strings"
)

func onInputChange(onValueChange chan string, displayViewName string) app.KeyHandler {
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

		onValueChange <- data
		g.Update(func(gui *gocui.Gui) error {
			vlog, err := g.View(displayViewName)
			if err != nil {
				return err
			}

			color := colors.ResetColor
			if strings.HasPrefix(data, "/") {
				color = colors.LightGreen
			}

			if _, err := colors.ColorFprintln(vlog, color, "ME:", data); err != nil {
				return err
			}

			return nil
		})

		v.Clear()
		if err := v.SetCursor(0, 0); err != nil {
			return err
		}
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
		return nil
	}
}

func onChannelChanged(
	viewName string,
	formatter func(string) string,
	//customAction func(*gocui.Gui, string) bool,
) app.LogConsumer {
	return func(g *gocui.Gui, rawMsg interface{}) {
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View(viewName)
			if err != nil {
				return err
			}

			msg, ok := rawMsg.(string)
			if !ok {
				return errors.New(fmt.Sprintf("Log Consumer got msg with type %T expected string", rawMsg))
			}

			msg = strings.Trim(msg, "\r\n")
			//if customAction == nil || customAction(g, msg) {
			if formatter != nil {
				msg = formatter(msg)
			}

			if _, err := fmt.Fprintln(v, msg); err != nil {
				return err
			}
			//}
			return nil
		})
	}
}
