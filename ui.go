package main

import (
	"fmt"
	"levi.ori/p2p-chat/src/ui"
	"levi.ori/p2p-chat/src/utils/colors"
	"strings"
)

//const CUTSET = " \r\n" + string(common.ResetColor)

var (
	chatColors = map[string]colors.Color{
		"(PM)":    colors.Gold,
		"(SHELL)": colors.Cyan,
	}

	logColors = map[string]colors.Color{
		"[INFO":  colors.LightPurple,
		"[DEBUG": colors.LightCyan,
		"[ERROR": colors.LightRed,
	}
)

func NewInputWidget(name string) *ui.Widget {
	return ui.NewWidget(
		"input",
		fmt.Sprintf("%v, What's On Your Mind?", name),
		true,
		true,
		true,
		func(maxX int) int { return 0 },
		func(maxY int) int { return 3*maxY/4 - 3 },
		func(maxX int) int { return maxX - 1 },
		func(maxY int) int { return 3*maxY/4 - 1 },
		true,
		nil,
	)
}

func NewLogWidget() *ui.Widget {
	return ui.NewWidget(
		"log",
		"Log",
		true,
		false,
		true,
		func(maxX int) int { return 0 },
		func(maxY int) int { return 3 * maxY / 4 },
		func(maxX int) int { return maxX - 1 },
		func(maxY int) int { return maxY - 1 },
		false,
		nil,
	)
}

func NewChatWidget() *ui.Widget {
	return ui.NewWidget(
		"chat",
		"Conversation",
		true,
		false,
		true,
		func(maxX int) int { return 0 },
		func(maxY int) int { return 0 },
		func(maxX int) int { return maxX / 3 * 2 },
		func(maxY int) int { return 3*maxY/4 - 4 },
		false,
		nil,
	)
}

func NewUsersWidget() *ui.Widget {
	return ui.NewWidget(
		"users",
		"Users",
		true,
		false,
		true,
		func(maxX int) int { return maxX/3*2 + 1 },
		func(maxY int) int { return 0 },
		func(maxX int) int { return maxX - 1 },
		func(maxY int) int { return maxY / 2 },
		false,
		nil,
	)
}

func NewHelpWidget() *ui.Widget {
	return ui.NewWidget(
		"help",
		"Help",
		true,
		false,
		true,
		func(maxX int) int { return maxX/3*2 + 1 },
		func(maxY int) int { return maxY/2 + 1 },
		func(maxX int) int { return maxX - 1 },
		func(maxY int) int { return 3*maxY/4 - 4 },
		false,
		[]string{
			fmt.Sprintf("%-9v <ip> <port>", "/connect"),
			fmt.Sprintf("%-9v <name> <message...>", "/pm"),
			fmt.Sprintf("%-9v <name> <command...>", "/shell"),
			"/exit",
		},
	)
}

//func sendData(input chan string) func(*gocui.Gui, *gocui.View) error {
//	return func(g *gocui.Gui, v *gocui.View) error {
//		reader := bufio.NewReader(v)
//		data, err := reader.ReadString('\n')
//		if err == io.EOF {
//			return nil
//		}
//
//		if err != nil {
//			return err
//		}
//
//		data = strings.Trim(data, "\r\n")
//		if data == "/exit" {
//			return gocui.ErrQuit
//		}
//
//		input <- data
//		g.Update(func(gui *gocui.Gui) error {
//			vlog, err := g.View("chat")
//			if err != nil {
//				return err
//			}
//
//			color := common.ResetColor
//			if strings.HasPrefix(data, "/") {
//				color = common.LightGreen
//			}
//
//			if _, err := common.ColorFprintln(vlog, color, "ME:", data); err != nil {
//				return err
//			}
//
//			return nil
//		})
//
//		v.Clear()
//		if err := v.SetCursor(0, 0); err != nil {
//			return nil
//		}
//		if err := v.SetOrigin(0, 0); err != nil {
//			return nil
//		}
//		return nil
//	}
//}
//
//func uiMain(name string, logChannel chan string, chatChannel chan string, inputChannel chan string) {
//	logChan = logChannel
//	outputChan = inputChannel
//
//	g, err := gocui.NewGui(gocui.OutputNormal)
//	if err != nil {
//		log.Panicln(err)
//	}
//	defer g.Close()
//
//	g.Highlight = true
//	g.Cursor = true
//	g.InputEsc = true
//	g.SelFgColor = gocui.ColorGreen
//
//	g.SetManagerFunc(func(g *gocui.Gui) error {
//		return layout(g, name)
//	})
//
//	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
//		log.Panicln(err)
//	}
//	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
//		log.Panicln(err)
//	}
//	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, sendData(inputChannel)); err != nil {
//		log.Panicln(err)
//	}
//
//	go handleViewWithChannel(g, logChannel, "log", prefixFormatter(logColors), nil)
//	go handleViewWithChannel(g, chatChannel, "chat", prefixFormatter(chatColors), shellPrompt)
//
//	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
//		log.Panicln(err)
//	}
//}
//
//func shellPrompt(g *gocui.Gui, s string) bool {
//	if !strings.Contains(s, "(SHELL)") {
//		return true
//	}
//
//	parts := strings.SplitN(s, ":", 2)
//	command := strings.Trim(parts[1], CUTSET)
//
//	parts = strings.SplitN(parts[0], ")", 2)
//	name := strings.Trim(parts[1], CUTSET)
//
//	commandParts := strings.Split(command, " ")
//	app := commandParts[0]
//	arguments := commandParts[1:]
//
//	cmd := exec.Command(app, arguments...)
//	stdout, err := cmd.Output()
//	if err != nil {
//		common.Error(logChan, "Failed to execute command %v from %v", command, name)
//		return true
//	}
//
//	data := strings.Split(string(stdout), "\n")
//	for i := 0; i < len(data); i++ {
//		data[i] = fmt.Sprintf("%v %v", common.Ok, data[i])
//	}
//
//	outputChan <- strings.Join(data, "\n")
//	return true
//
//	//title := fmt.Sprintf("%v want's to run the shell command", name)
//	//
//	//commandLength := len(command)
//	//minimumWidth := len(title)
//	//
//	//width := commandLength
//	//height := 5
//	//if commandLength < minimumWidth {
//	//	width = minimumWidth
//	//}
//	//
//	//maxX, maxY := g.Size()
//	//x0 := maxX/2 - width/2
//	//y0 := maxY/2 - height/2
//	//x1 := maxX/2 + width/2 + 3
//	//y1 := maxY/2 + height/2
//	//
//	//if v, err := g.SetView("shell-prompt", x0, y0, x1, y1); err != nil {
//	//	if err != gocui.ErrUnknownView {
//	//		log.Panicln(err)
//	//	}
//	//
//	//	v.title = title
//	//	v.editable = false
//	//	v.wrap = false
//	//	v.autoscroll = false
//	//
//	//	if _, err := fmt.Fprint(v, center(command, x1-x0-2, " ")); err != nil {
//	//		log.Panicln(err)
//	//	}
//	//}
//	//
//	//return false
//}
//

func prefixFormatter(prefixToColor map[string]colors.Color) func(string) string {
	return func(s string) string {
		finalColor := colors.ResetColor
		for prefix, color := range prefixToColor {
			if strings.HasPrefix(s, prefix) {
				finalColor = color
				break
			}
		}

		return colors.ColorSprintf(finalColor, s)
	}
}

//func center(s string, n int, fill string) string {
//	filler := strings.Repeat(fill, n/2)
//
//	return fmt.Sprintf("%v%v%v", filler, s, filler)
//}
