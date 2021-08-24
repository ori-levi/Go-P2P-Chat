package main

//
//import (
//	"bufio"
//	"github.com/jroimartin/gocui"
//	"io"
//	"strings"
//)
//
//
//import (
//	"bufio"
//	"fmt"
//	"io"
//	"log"
//	"os/exec"
//	"strings"
//
//	"github.com/jroimartin/gocui"
//)
//
//const CUTSET = " \r\n" + string(common.ResetColor)
//
//var (
//	logChan    chan string
//	outputChan chan string
//
//	chatColors = map[string]common.Color{
//		"(PM)":    common.Gold,
//		"(SHELL)": common.Cyan,
//	}
//
//	logColors = map[string]common.Color{
//		"[INFO":  common.LightPurple,
//		"[DEBUG": common.LightCyan,
//		"[ERROR": common.LightRed,
//	}
//
//
//
//
//func handleViewWithChannel(
//	g *gocui.Gui,
//	channel chan string,
//	viewName string,
//	formatter func(string) string,
//	customAction func(*gocui.Gui, string) bool,
//) {
//	for {
//		msg := <-channel
//
//		g.Update(func(g *gocui.Gui) error {
//			v, err := g.View(viewName)
//			if err != nil {
//				return err
//			}
//
//			msg := strings.Trim(msg, "\r\n")
//			if customAction == nil || customAction(g, msg) {
//				if formatter != nil {
//					msg = formatter(msg)
//				}
//
//				if _, err := fmt.Fprintln(v, msg); err != nil {
//					return err
//				}
//			}
//			return nil
//		})
//	}
//}
//
//
//
//func uiMain(name string, logChannel chan string, chatChannel chan string, inputChannel chan string) {
//	logChan = logChannel
//	outputChan = inputChannel
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
//	//	v.Title = title
//	//	v.Editable = false
//	//	v.Wrap = false
//	//	v.Autoscroll = false
//	//
//	//	if _, err := fmt.Fprint(v, center(command, x1-x0-2, " ")); err != nil {
//	//		log.Panicln(err)
//	//	}
//	//}
//	//
//	//return false
//}
//
//func prefixFormatter(prefixToColor map[string]common.Color) func(string) string {
//	return func(s string) string {
//
//		finalColor := common.ResetColor
//		for prefix, color := range prefixToColor {
//			if strings.HasPrefix(s, prefix) {
//				finalColor = color
//				break
//			}
//		}
//
//		return common.ColorSprintf(finalColor, s)
//	}
//}
//
//func center(s string, n int, fill string) string {
//	filler := strings.Repeat(fill, n/2)
//
//	return fmt.Sprintf("%v%v%v", filler, s, filler)
//}
