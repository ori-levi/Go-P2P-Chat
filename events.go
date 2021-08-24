package main

import (
	"bufio"
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	ui "levi.ori/p2p-chat/src/ui/widgets"
	"strings"
)

func onInputChange(onValueChange chan string) ui.Handler {
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
			vlog, err := g.View("chat")
			if err != nil {
				return err
			}

			//color := common.ResetColor
			//if strings.HasPrefix(data, "/") {
			//	color = common.LightGreen
			//}

			//if _, err := common.ColorFprintln(vlog, color, "ME:", data); err != nil {
			if _, err := fmt.Fprintln(vlog, "ME:", data); err != nil {
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
